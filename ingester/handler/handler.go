package handler

import (
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	apiproto.UnimplementedEventServer
}

// SayHello generates response to a Ping request
func (s *Server) RecordEvents(ctx context.Context, in *apiproto.Record) (*apiproto.RecordSummary, error) {
	record := in.GetRecord()
	for k, v := range record {
		log.Printf("%s->%s", k, v)
	}
	log.Printf("%s %s %s", in.GetTimestamp(), in.GetTag(), in.GetRecord())
	return &apiproto.RecordSummary{EventCount: 99}, nil
}
