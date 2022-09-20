package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ozan1338/api"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/gapi"
	"github.com/ozan1338/pb"
	"github.com/ozan1338/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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