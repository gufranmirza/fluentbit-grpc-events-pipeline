package server

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/api/apiproto"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/utils"
	"google.golang.org/protobuf/proto"
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

		buffer, err := proto.Marshal(event)
		if err != nil {
			log.Printf("Failed to marshal event %v", err)
		}

		if s.config.Print {
			utils.Print(event, s.config.Decrypt)
		}

		s.producer.Produce(buffer)
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
