package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config holds the server specific config
type Config struct {
	Decrypt       bool                       // Decrypt messages received from fluentbit-agent
	AccessTokenDB map[string]apiproto.Config // List of access tokens and their config
}

// Server represents the gRPC server
type Server struct {
	apiproto.UnimplementedEventServiceServer
	config *Config
}

// main start a gRPC server and waits for connection
func NewServer(c *Config) *Server {
	return &Server{
		config: c,
	}
}

func (s *Server) Start() {
	// create a listener on TCP port 7777
	url := fmt.Sprintf(":%d", 7777)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile("../cert/server-cert.pem", "../cert/server-key.pem")
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(Unary()),
		grpc.StreamInterceptor(Stream()),
	)

	// attach the Event service to the server
	apiproto.RegisterEventServiceServer(grpcServer, s)

	confBytes, err := ioutil.ReadFile("./access-tokens-db.json")
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}

	err = json.Unmarshal(confBytes, &s.config.AccessTokenDB)
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}

	// start the server
	fmt.Printf("Starting grpc-server at => %s\n", url)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
