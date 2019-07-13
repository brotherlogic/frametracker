package main

import (
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/frametracker/proto"
)

//RecordStatus records a frame status
func (s *Server) RecordStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	defer s.save(ctx)

	for i, status := range s.config.States {
		if status.TokenHash == in.Status.TokenHash {
			s.config.States[i] = in.Status
			return &pb.StatusResponse{}, nil
		}
	}

	s.config.States = append(s.config.States, in.Status)
	return &pb.StatusResponse{}, nil
}
