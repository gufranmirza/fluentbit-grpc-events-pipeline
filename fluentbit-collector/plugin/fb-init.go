package main

import (
	"context"
	"fmt"
	"os"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func (plugin *Plugin) connectToIngest() error {
	var err error

	// read and set access key
	plugin.config.AccessKey = os.Getenv("ACCESS_KEY")
	plugin.config.AccessToken = os.Getenv("ACCESS_TOKEN")

	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile("/fluent-bit/bin/ca-cert.pem", "x.test.example.com")
	if err != nil {
		return fmt.Errorf("failed to load credentials: %v", err)
	}

	oauth := oauth.NewOauthAccess(&oauth2.Token{AccessToken: plugin.config.AccessToken})
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(oauth),
		// oauth.NewOauthAccess requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(creds),
	}

	// Dial
	conn, err := grpc.Dial("host.docker.internal:7777", opts...)
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
