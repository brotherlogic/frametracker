package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/frametracker/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

const (
	// KEY - where we store sale info
	KEY = "/github.com/brotherlogic/frametracker/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config *pb.Config
	bad    int64
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		config:   &pb.Config{},
	}
	return s
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterFrameTrackerServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.save(ctx)
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		err := s.load(ctx)
		return err
	}

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	hashes := []string{}
	lastUpdate := make(map[string]string)
	for _, state := range s.config.States {
		hashes = append(hashes, state.TokenHash)
		lastUpdate[state.Origin] = fmt.Sprintf("%v", time.Unix(state.NewestFileDate, 0))
	}
	return []*pbg.State{
		&pbg.State{Key: "bad", Value: s.bad},
		&pbg.State{Key: "last", TimeValue: s.config.LastReceive},
		&pbg.State{Key: "states", Value: int64(len(s.config.States))},
		&pbg.State{Key: "hashes", Text: fmt.Sprintf("%v", hashes)},
		&pbg.State{Key: "updates", Text: fmt.Sprintf("%v", lastUpdate)},
	}
}

func (s *Server) checkTime(ctx context.Context) error {
	if time.Now().Sub(time.Unix(s.config.LastReceive, 0)) > time.Hour*24 {
		s.bad++
		//s.RaiseIssue(ctx, "Frame Tracker has not processed anything", fmt.Sprintf("No updates since %v", time.Unix(s.config.LastReceive, 0)), false)
	}
	return nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.DialMaster)
	server.PrepServer()
	server.Register = server
	err := server.RegisterServerV2("frametracker", false, false)
	if err != nil {
		log.Fatalf("Cannot register: %v", err)
	}

	server.RegisterRepeatingTask(server.checkTime, "check_time", time.Minute)

	fmt.Printf("%v", server.Serve())
}
