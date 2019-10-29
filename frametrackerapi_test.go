package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/frametracker/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")

	return s
}

func TestAddState(t *testing.T) {
	s := InitTest()

	_, err := s.RecordStatus(context.Background(), &pb.StatusRequest{Status: &pb.Status{TokenHash: "blah"}})

	if err != nil {
		t.Fatalf("Error recording state: %v", err)
	}

	if len(s.config.States) != 1 {
		t.Errorf("State was not recorded: %v", s.config.States)
	}

	_, err = s.RecordStatus(context.Background(), &pb.StatusRequest{Status: &pb.Status{TokenHash: "blah"}})

	if err != nil {
		t.Fatalf("Error recording state: %v", err)
	}

	if len(s.config.States) != 1 {
		t.Errorf("State was not recorded correctly: %v", s.config.States)
	}

}
