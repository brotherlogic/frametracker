package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/frametracker/proto"
)

//RecordStatus records a frame status
func (s *Server) RecordStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	defer s.save(ctx)

	s.Log(fmt.Sprintf("%v", time.Now().Sub(time.Unix(in.Status.NewestFileDate/1000, 0))))
	if time.Now().Sub(time.Unix(in.Status.NewestFileDate/1000, 0)) > time.Hour*24*7 {
		s.RaiseIssue(ctx, "Picture Frame Behind", fmt.Sprintf("%v is behind", in.Status.TokenHash), false)
	}

	for i, status := range s.config.States {
		if status.TokenHash == in.Status.TokenHash {
			s.config.States[i] = in.Status
			return &pb.StatusResponse{}, nil
		}
	}

	s.config.States = append(s.config.States, in.Status)
	return &pb.StatusResponse{}, nil
}
