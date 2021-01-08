package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/frametracker/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	frameAge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "frametracker_age",
		Help: "ImageAge",
	}, []string{"frame"})
)

//RecordStatus records a frame status
func (s *Server) RecordStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	s.Log(fmt.Sprintf("%v -> %v (%v)", in.Status.Origin, time.Now().Sub(time.Unix(in.Status.NewestFileDate/1000, 0)), in.Status.NewestFile))
	frameAge.With(prometheus.Labels{"frame": in.Status.Origin}).Set(float64(time.Now().Sub(time.Unix(in.Status.NewestFileDate/1000, 0)).Seconds()))

	if time.Now().Sub(time.Unix(in.Status.NewestFileDate/1000, 0)) > time.Hour*24*7 {
		if in.Status.Origin != "natframe" {
			s.RaiseIssue("Picture Frame Behind", fmt.Sprintf("%v is behind -> %v", in.Status.Origin, time.Unix(in.Status.NewestFileDate/1000, 0)))
		}
	}

	for i, status := range config.States {
		if status.TokenHash == in.Status.TokenHash && status.Origin == in.Status.Origin {
			config.States[i] = in.Status
			return &pb.StatusResponse{}, nil
		}
	}

	config.States = append(config.States, in.Status)
	return &pb.StatusResponse{}, s.save(ctx, config)
}
