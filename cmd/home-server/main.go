package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/RasmusLindroth/home/grpchome/golang"
	"github.com/RasmusLindroth/home/home"
	"github.com/eriklupander/tradfri-go/tradfri"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHomeServiceServer
}

func (s *server) RunAction(ctx context.Context, in *pb.RunActionRequest) (*pb.RunActionResponse, error) {
	err := home.ExecAction(client, in.GetRoom(), in.GetLamp(), in.GetAction(), in.GetValue())
	return &pb.RunActionResponse{}, err
}

var client *tradfri.TradfriClient

func main() {
	conf, err := home.ParseConfig()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	client = tradfri.NewTradfriClient(
		conf.IKEA.Gateway,
		conf.IKEA.ClientID,
		conf.IKEA.PSK,
	)

	lis, err := net.Listen("tcp", conf.GetAddress())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHomeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
