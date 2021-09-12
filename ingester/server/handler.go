package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/encryption"
)

// SayHello generates response to a Ping request
func (s *Server) SendEvent(stream apiproto.EventService_SendEventServer) error {
	var count int
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			log.Printf("--> Messages sent to kafka: %v", count)

			return stream.SendAndClose(&apiproto.EResponse{
				Status: apiproto.EventCode_SUCCESS,
			})
		}
		if err != nil {
			return err
		}

		key, ok := s.config.AccessTokenDB[event.AccessKey]
		if ok && s.config.Decrypt && key.EncryptionKey != "" {
			msg, err := encryption.Decrypt(string(key.EncryptionKey), event.Message)
			if err != nil {
				fmt.Printf("Failed to decrypt message %v/n", err)
			}
			event.Message = msg
		}

		s.producer.Produce([]byte(fmt.Sprintf("%v", event)))
		count++
	}
}

// GetFeature returns the feature at the given point.
func (s *Server) ExchangeConfig(ctx context.Context, accessKey *apiproto.AccessKey) (*apiproto.Config, error) {
	conf, ok := s.config.AccessTokenDB[accessKey.AccessKey]
	if !ok {
		return nil, errors.New("invalid Access Key")
	}

	return &conf, nil
}
