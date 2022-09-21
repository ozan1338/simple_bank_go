package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader 	= "grpcgateway-user-agent"
	userAgentHeader 			= "user-agent"
	xForwardedForHeader			= "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP string
}

func (server *Server) extractMetaData(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		// log.Printf("md: %v\n",md)
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			// log.Printf("user agent: %v\n",userAgents[0])
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		clientAPIs := md.Get(xForwardedForHeader)

		if len(clientAPIs) > 0 {
			log.Printf("ip : %v\n", clientAPIs[0])
			mtdt.ClientIP = clientAPIs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}