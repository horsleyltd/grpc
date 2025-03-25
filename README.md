# Grpc

## Overview

gRPC, which stands for "Google Remote Procedure Call," is a high-performance, language-agnostic remote procedure call (RPC) framework. It was originally developed by Google and is now open-source. gRPC is designed to simplify the development of efficient and robust microservices by providing a standardized and efficient way for services to communicate with each other.

Key features of gRPC include:

Efficient Communication: gRPC uses Protocol Buffers (Protobuf) for message serialization, which results in compact and efficient data transmission. It also utilizes HTTP/2 for multiplexing and reducing latency.
Language Agnostic: gRPC supports multiple programming languages, making it suitable for building polyglot microservices.
Code Generation: It generates client and server code, reducing the boilerplate code required for communication.
Streaming: gRPC supports streaming, enabling various communication patterns, including unary (single request, single response), server streaming, client streaming, and bidirectional streaming.

## Install grpc

```bash
# install the protocol buffer compiler, protocol buffer plugins for Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# verify `~/go/bin` is in PATH
```

## Initialise go project

```bash
go mod init github.com/horsleyltd/grpc
go mod tidy
```

## Create service

```bash
mkdir service && touch service/service.proto
```

Append the following to `service/service.proto`:

```proto
// Here we define the gRPC service and the request and response types that use protocol buffers. 
// This file is used to generate our gRPC server and stub interfaces.

syntax = "proto3";

option go_package = "github.com/horsleyltd/grpc/service";

package service;

message Request {}

message Response {
  string message = 1;
}

service Service {
  // Simple RPC 
  rpc RequestResponse(Request) returns (Response) {}
  // Server-side streaming RPC
  rpc RequestResponseStream(Request) returns (stream Response) {}
  // Client-side streaming RPC
  rpc StreamRequestResponse(stream Request) returns (Response) {}
}
```

Run the following commands:

```bash
cd service
protoc --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       service.proto

# This will create 2 files in the service directory:
# service.pb.go:
#   contains protocol buffer code to populate, serialize, and retrieve request and response message types
# service_grpc.pb.go:
#   contains an interface type (stub) for clients to call with the methods defined in the serviceService service, 
#   contains an interface type for servers to implement, also with the methods defined in the serviceService service
```

## Create server & client application

```bash
mkdir server client && touch server/server.go client/client.go
```

As demonstrated in `server.go` & `client.go` functions are defined that utilise the protocol buffer code in `service/` to create:

- Simple RPC
- Server-side streaming RPC
- Client-side streaming RPC
- Bi-directional streaming RPC
