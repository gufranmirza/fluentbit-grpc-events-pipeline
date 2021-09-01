package main

import (
	"context"
	"fmt"
	"os"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (plugin *Plugin) connectToIngest() error {
	var err error

	// read and set access key
	plugin.config.AccessKey = os.Getenv("ACCESS_KEY")

	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile("/fluent-bit/bin/ca-cert.pem", "x.test.example.com")
	if err != nil {
		return fmt.Errorf("failed to load credentials: %v", err)
	}

	// Dial
	conn, err := grpc.Dial("host.docker.internal:7777", grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("did not connect: %s", err)
	}
	plugin.conn = conn

	// setup streaming
	plugin.eventClient = apiproto.NewEventServiceClient(conn)

	return nil
}

func (plugin *Plugin) exchangeConfig(accessKey string) error {
	config, err := plugin.eventClient.ExchangeConfig(context.Background(), &apiproto.AccessKey{AccessKey: accessKey})
	if err != nil {
		return fmt.Errorf("failed to exchange config with error %v", err)
	}
	plugin.config = config

	return nil
}
