package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "../../proto"
	"google.golang.org/grpc"
)

// PORT ...
const PORT = 9002

func printLists(cli pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := cli.List(context.Background(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	return nil
}

func printRecord(cli pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := cli.Record(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(cli pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := cli.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	stream.CloseSend()
	return nil
}

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	client := pb.NewStreamServiceClient(conn)
	err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 2019}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 2019}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2019}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}
