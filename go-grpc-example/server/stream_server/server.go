package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "../../proto"
	"google.golang.org/grpc"
)

// StreamService ...
type StreamService struct{}

// PORT ...
const PORT = 9002

// List ...
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Record ...
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

// Route ...
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gRPC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err != nil {
			return err
		}
		n++
		log.Printf("Stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

func main() {
	server := grpc.NewServer()
	// pb.RegisterStreamServiceServer(server, &StreamService{})
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}
