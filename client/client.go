package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/horsleyltd/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Simple RPC
func RequestResponse(ctx context.Context, serviceClient service.ServiceClient) {
	// RequestResponse returns a response for a given request
	response, err := serviceClient.RequestResponse(ctx, &service.Request{})
	if err != nil {
		log.Fatalf("error calling function RequestResponse: %v", err)
	}

	log.Print(response.Message)
}

// Server-side streaming RPC
func RequestResponseStream(ctx context.Context, serviceClient service.ServiceClient) {
	// RequestResponseStream returns a stream of responses for a given request
	serverStream, err := serviceClient.RequestResponseStream(ctx, &service.Request{})
	if err != nil {
		log.Fatalf("error calling function RequestStreamResponse: %v", err)
	}

	// Recv is non blocking
	response, err := serverStream.Recv()
	if err != nil {
		log.Fatalf("error calling function Recv: %v", err)
	}

	log.Print(response.Message)

	// RecvMsg blocks while waiting for msgs
	for {
		response := &service.Response{}
		err := serverStream.RecvMsg(response)
		if err == io.EOF {
			log.Print("Stream is empty")
			break
		}
		if err != nil {
			log.Fatalf("error calling function RecvMsg: %v", err)
		}

		log.Print(response.Message)
	}
}

// Client-side streaming RPC
func StreamRequestResponse(ctx context.Context, serviceClient service.ServiceClient) {
	// StreamRequestResponse returns a stream for the client's requests and server's response
	clientStream, err := serviceClient.StreamRequestResponse(ctx)
	if err != nil {
		log.Fatalf("error calling function StreamRequestResponse: %v", err)
	}

	for i := 0; i < 3; i++ {
		err := clientStream.Send(&service.Request{})
		if err != nil {
			log.Fatalf("error calling function SendMsg: %v", err)
		}
	}

	// CloseAndRecv signals to server that request stream is completed, then waits for server to write response and close the stream
	response, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error calling function CloseAndRecv: %v", err)
	}

	log.Printf(response.Message)
}

// Bi-directional streaming RPC
func StreamRequestResponseStream(ctx context.Context, serviceClient service.ServiceClient) {
	// RequestResponseStream returns a bi-directional stream
	bidiStream, err := serviceClient.StreamRequestResponseStream(context.Background())
	if err != nil {
		log.Fatalf("error calling function StreamRequestResponseStream: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := bidiStream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("error calling function Recv: %v", err)
			}
			log.Printf(response.Message)
		}
	}()
	for i := 0; i < 3; i++ {
		err := bidiStream.Send(&service.Request{})
		if err != nil {
			log.Fatalf("error calling function Send: %v", err)
		}
	}

	bidiStream.CloseSend()
	<-waitc
}

func main() {
	// create gRPC channel to communicate with the gRPC server
	conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error calling function NewClient: %v", err)
	}
	defer conn.Close()

	// Create gRPC stub (client) to perform RPCs
	serviceClient := service.NewServiceClient(conn)

	// Set context with timeout (see deadlines)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Simple RPC
	RequestResponse(ctx, serviceClient)

	// Server-side streaming RPC
	RequestResponseStream(ctx, serviceClient)

	// Client-side streaming RPC
	StreamRequestResponse(ctx, serviceClient)

	// Bi-directional streaming RPC
	StreamRequestResponseStream(ctx, serviceClient)
}
