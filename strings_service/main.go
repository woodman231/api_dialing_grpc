package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/woodman231/api_dialing_grpc/protos/stringspb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implment StringsServiceServer
type server struct {
	pb.UnimplementedStringServiceServer
}

// MakeUpperCase implments StringsServiceServer.MakeUpperCase
func (s *server) MakeUpperCase(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResult, error) {
	requestedString := in.GetInputString()

	log.Printf("Received: MakeUpperCase request for %v", requestedString)

	return &pb.OperationResult{
		OutputString: strings.ToUpper(requestedString),
	}, nil
}

// MakeLowerCase implements StringsSErviceServer.MakeLowerCase
func (s *server) MakeLowerCase(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResult, error) {
	requestedString := in.GetInputString()

	log.Printf("Received: MakeLowerCase request for %v", requestedString)

	return &pb.OperationResult{
		OutputString: strings.ToLower(requestedString),
	}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStringServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
