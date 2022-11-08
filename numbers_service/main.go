package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/woodman231/api_dialing_grpc/protos/numberspb"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

// server is used to implement NumbersServiceServer
type server struct {
	pb.UnimplementedNumberServiceServer
}

// AddTwoNumbers implements NumbersServiceServer.AddTwoNumbers
func (s *server) AddTwoNumbers(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResult, error) {
	numberOne := in.GetInputNumberOne()
	numberTwo := in.GetInputNumberTwo()

	log.Printf("Received: AddTwoNumbers Request %v, %v", numberOne, numberTwo)

	result := numberOne + numberTwo

	return &pb.OperationResult{
		OutputNumber: result,
	}, nil
}

// SubtractTwoNumbers implments NumbersServiceServer.SubtractTwoNumbers
func (s *server) SubtractTwoNumbers(ctx context.Context, in *pb.OperationRequest) (*pb.OperationResult, error) {
	numberOne := in.GetInputNumberOne()
	numberTwo := in.GetInputNumberTwo()

	log.Printf("Received: SubtractTwoNumbers Request %v, %v", numberOne, numberTwo)

	result := numberOne - numberTwo

	return &pb.OperationResult{
		OutputNumber: result,
	}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNumberServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
