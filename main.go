package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ozan1338/api"
	db "github.com/ozan1338/db/sqlc"
	gapi "github.com/ozan1338/gapi"
	pb "github.com/ozan1338/pb"
	"github.com/ozan1338/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config,err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config: ",err)
	}

	conn, err := sql.Open(config.DBDrvier, config.DBSource)

	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	
	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server,err := gapi.NewServer(config,store)
	
	if err != nil {
		log.Fatal("can't start server: ", err)
	}
	
	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("can't create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Can't start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server,err := gapi.NewServer(config,store)
	
	if err != nil {
		log.Fatal("can't start server: ", err)
	}
	
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("can't register handler server")
	}

	

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("can't create listener")
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Can't start HTTP gateway server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server,err := api.NewServer(config,store)

	if err != nil {
		log.Fatal("can't start server: ", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ",err)
	}
}