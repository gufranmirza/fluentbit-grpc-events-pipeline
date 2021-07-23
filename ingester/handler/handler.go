package handler

import (
	"io"
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
)

// Server represents the gRPC server
type Server struct {
	apiproto.UnimplementedEventServiceServer
}

// SayHello generates response to a Ping request
func (s *Server) SendEvent(stream apiproto.EventService_SendEventServer) error {
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

		log.Printf("%+v \n\n", event)
	}
}
