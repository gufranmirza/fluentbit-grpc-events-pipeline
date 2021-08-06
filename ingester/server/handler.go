package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/encryption"
)

// SayHello generates response to a Ping request
func (s *Server) SendEvent(stream apiproto.EventService_SendEventServer) error {
	// Read Public key for encryption of Events passed over wire
	pubKey, err := ioutil.ReadFile("../cert/encryption_aes.pub")
	if err != nil {
		log.Fatalf("Failed to read key %v \n", err)
	}

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&apiproto.EResponse{
				Status: apiproto.EventCode_SUCCESS,
			})
		}
		if err != nil {
			return err
		}

		fmt.Println("==============================================")
		if s.config.Decrypt {
			msg, err := encryption.Decrypt(string(pubKey), event.Message)
			if err != nil {
				fmt.Printf("Failed to decrypt message %v/n", err)
			}
			event.Message = msg
		}
		fmt.Println(event)

	}
}

// GetFeature returns the feature at the given point.
func (s *Server) ExchangeConfig(ctx context.Context, accessKey *apiproto.AccesKey) (*apiproto.Config, error) {
	conf, ok := s.config.AccessTokenDB[accessKey.AccesKey]
	if !ok {
		return nil, errors.New("invalid Access Key")
	}

	return &conf, nil
}
