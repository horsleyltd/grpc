package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/horsleyltd/grpc/service"
	"google.golang.org/grpc"
)

// Embedded service interface generated from service definition in service.proto
type ServiceServer struct {
	service.UnimplementedServiceServer
}

// Simple RPC
func (*ServiceServer) RequestResponse(ctx context.Context, req *service.Request) (*service.Response, error) {
	resp := &service.Response{Message: "A Response to your request!"}
	return resp, nil
}

// Server-side streaming RPC, sends a stream of 3 messages
func (s *ServiceServer) RequestResponseStream(request *service.Request, stream grpc.ServerStreamingServer[service.Response]) error {
	for i := 1; i < 4; i++ {
		if err := stream.Send(&service.Response{
			Message: fmt.Sprintf("Response %d", i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// Client-side streaming RPC
func (s *ServiceServer) StreamRequestResponse(clientStream grpc.ClientStreamingServer[service.Request, service.Response]) error {
	requestCount := 0
	for {
		_, err := clientStream.Recv()
		if err == io.EOF {
			return clientStream.SendAndClose(&service.Response{Message: fmt.Sprintf("Request stream count = %d", requestCount)})
		}

		if err != nil {
			log.Fatalf("error calling function Recv: %v", err)
		}

		requestCount++
		log.Printf("Request: %d", requestCount)
	}
}

// Bi-directional streaming RPC
func (s *ServiceServer) StreamRequestResponseStream(bidiStream grpc.BidiStreamingServer[service.Request, service.Response]) error {
	requestResponseCount := 0
	for {
		_, err := bidiStream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error calling function Recv: %v", err)
		}

		requestResponseCount++
		response := &service.Response{Message: fmt.Sprintf("Response %d", requestResponseCount)}

		err = bidiStream.Send(response)
		if err != nil {
			log.Fatalf("error calling function Send: %v", err)
		}

		log.Printf("Request: %d", requestResponseCount)
	}

	// for {
	// 	request, err := stream.Recv()
	// 	if err == io.EOF {
	// 		return nil
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// key := serialize(request)
	// 	// look for notes to be sent to client
	// 	for _, response := range s.routeNotes[key] {
	// 		if err := stream.Send(note); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
}

func main() {
	// Initialise listener for client requests
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen on port 50053: %v", err)
	}

	// Initialise gRPC server
	grpcServer := grpc.NewServer()

	// Register (bind) embedded service interface to grpc server
	service.RegisterServiceServer(grpcServer, &ServiceServer{})

	// Start server
	log.Printf("gRPC server listening at %v", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
