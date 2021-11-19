package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/api/apiproto"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config holds the server specific config
type Config struct {
	Decrypt       bool                       // Decrypt messages received from fluentbit-agent
	Print         bool                       // Print messages to console
	AccessTokenDB map[string]apiproto.Config // List of access tokens and their config
}

// Server represents the gRPC server
type Server struct {
	apiproto.UnimplementedEventServiceServer
	config   *Config
	producer *kafka.Producer
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
	defer grpcServer.Stop()

	// attach the Event service to the server
	apiproto.RegisterEventServiceServer(grpcServer, s)

	// read local database of access tokens
	confBytes, err := ioutil.ReadFile("../access-tokens-db.json")
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}
	err = json.Unmarshal(confBytes, &s.config.AccessTokenDB)
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}

	// start kafka producer
	producer, err := kafka.NewProducer("fb-kafka", []string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatalf("Failed to connect to kafka %v \n", err)
	}
	s.producer = producer
	defer s.producer.Close()

	// start the server
	fmt.Printf("Starting grpc-server at => %s\n", url)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
