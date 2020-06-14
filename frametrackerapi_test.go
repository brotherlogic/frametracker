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
	s.KSclient.Save(context.Background(), KEY, &pb.Config{})

	return s
}

func TestFail(t *testing.T) {
	s := InitTest()
	s.GoServer.KSclient.Fail = true

	_, err := s.RecordStatus(context.Background(), &pb.StatusRequest{Status: &pb.Status{TokenHash: "blah"}})
	if err == nil {
		t.Errorf("Did not fail")
	}

}

func TestAddState(t *testing.T) {
	s := InitTest()

	_, err := s.RecordStatus(context.Background(), &pb.StatusRequest{Status: &pb.Status{TokenHash: "blah"}})

	if err != nil {
		t.Fatalf("Error recording state: %v", err)
	}

	_, err = s.RecordStatus(context.Background(), &pb.StatusRequest{Status: &pb.Status{TokenHash: "blah"}})

	if err != nil {
		t.Fatalf("Error recording state: %v", err)
	}
}
