package main

import (
	"log"
	"os"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/frametracker/proto"

	//Needed to pull in gzip encoding init
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("frametracker-cli", "frametracker")
	defer cancel()

	switch os.Args[1] {
	case "get":
		conn, err := grpc.Dial(os.Args[2], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer conn.Close()

		client := pb.NewFrameTrackerServiceClient(conn)

		_, err = client.RecordStatus(ctx, &pb.StatusRequest{
			Status: &pb.Status{
				Origin: "test",
			},
		})
		if err != nil {
			log.Fatalf("Err: %v")
		}
	}
}
