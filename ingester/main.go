package main

import (
	"fmt"
	"log"
	"net"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/ingester/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var encryptionKey string

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	url := fmt.Sprintf(":%d", 7777)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := handler.Server{}

	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile("../cert/server-cert.pem", "../cert/server-key.pem")
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
	)

	// attach the Ping service to the server
	apiproto.RegisterEventServiceServer(grpcServer, &s)

	// start the server
	fmt.Printf("Starting server at %s\n", url)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
