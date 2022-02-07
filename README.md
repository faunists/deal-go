# Deal - Go

[![test](https://github.com/faunists/deal-go/actions/workflows/test.yaml/badge.svg)](https://github.com/faunists/deal-go/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/faunists/deal-go/branch/main/graph/badge.svg?token=qFlORZnn09)](https://codecov.io/gh/faunists/deal-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/faunists/deal-go)](https://goreportcard.com/report/github.com/faunists/deal-go)

## Introduction

This plugin allows us to write [Consumer-Driver Contracts](https://martinfowler.com/articles/consumerDrivenContracts.html) tests!

__Deal__ generates some code for us:
- A Client to be used in the client side to mock the responses based on the contract
- A Stub Server to be used in the client side as the Client above, but you should run it as another application
- Server Test Function, where you should pass your server implementation to the function and all the contracts will be validated against it

You can check out an example project [here](https://github.com/faunists/deal-go-example).

## Installation

Assuming that you are using [Go Modules](https://github.com/golang/go/wiki/Modules), it's
recommended to use a [tool dependency](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
in order to track your tools version:

```go
//go:build tools
// +build tools

package tools

import (
    _ "github.com/faunists/deal-go/protoc-gen-go-deal"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

Once you have added the required packages run `go mod tidy` to resolve the versions and then
install them by running:

```shell
go install \
    github.com/faunists/deal-go/protoc-gen-go-deal \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## Usage example

### Proto service

First you need a proto service:
```protobuf
syntax = "proto3";

import "google/protobuf/struct.proto";
import "deal/v1/contract/annotations.proto";

option go_package = "YOUR_PACKAGE_HERE/example";

message RequestMessage {
  string requestField = 1;
}

message ResponseMessage {
  int64 responseField = 1;
}

service MyService {
  rpc MyMethod(RequestMessage) returns (ResponseMessage);
}
```

### Contract file

After that you need to write the contract that should be respected, the contract file can be written using a JSON or YAML file.
You can set both, Success and Failures cases:

<details>
  <summary>JSON Contract</summary>

```json
{
  "name": "Some Name Here",
  "services": {
    "MyService": {
      "MyMethod": {
        "successCases": [
          {
            "description": "Should do something",
            "request": {
              "requestField": "VALUE"
            },
            "response": {
              "responseField": 42
            }
          }
        ],
        "failureCases": [
          {
            "description": "Some description here",
            "request": {
              "requestField": "ANOTHER_VALUE"
            },
            "error": {
              "errorCode": "NotFound",
              "message": "ANOTHER_VALUE NotFound"
            }
          }
        ]
      }
    }
  }
}
```
</details>

<details>
  <summary>YAML Contract</summary>

```yaml
name: Some Name Here
services:
  MyService:
    MyMethod:
      successCases:
        - description: Should do something
          request:
            requestField: VALUE
            someMessage:
              firstField: FIRST_FIELD_VALUE
            someEnum: TWO
          response:
            responseField: 42
            myList:
              - firstField: first
              - firstField: second
            intList: [1, 2, 3]
      failureCases:
        - description: Some description here
          request:
            requestField: ANOTHER_VALUE
          error:
            errorCode: NotFound
            message: ANOTHER_VALUE NotFound"
```
</details>

### Generating code

If you're using [buf](https://buf.build) just add the following entries to `buf.gen.yaml` and execute `buf generate` passing your contract file path:
```yaml
version: v1beta1
plugins:
  - name: go
    out: protogen
    opt: paths=source_relative
  - name: go-grpc
    out: protogen
    opt: paths=source_relative
  - name: go-deal
    out: protogen
    opt:
      - paths=source_relative
      - contract-file=contract.yml
```

> Disclaimer: You must be using `go-grpc` in order to make the things work

### Using generated client on tests

Here is an example using the generated client, in the example we're using it inside
a test, but it could be used anywhere!

```go
package main_test

import (
	"context"
	"testing"

	"github.com/faunists/deal-go-example/protogen/proto/example"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestClient(t *testing.T) {
	t.Run("should return a response", func(t *testing.T) {
		ctx := context.Background()
		expectedResp := &example.ResponseMessage{ResponseField: 42}
		// Generated client
		client := example.MyServiceContractClient{}

		actualResp, err := client.MyMethod(ctx, &example.RequestMessage{
			RequestField: "VALUE",
		})

		require.NoError(t, err)
		require.True(t, proto.Equal(expectedResp, actualResp))
	})

	t.Run("should return an error", func(t *testing.T) {
		ctx := context.Background()
		expectedError := status.Error(codes.NotFound, "ANOTHER_VALUE NotFound")
		// Generated client
		client := example.MyServiceContractClient{}

		_, err := client.MyMethod(ctx, &example.RequestMessage{
			RequestField: "ANOTHER_VALUE",
		})

		require.EqualError(t, err, expectedError.Error())
	})
}
```

### Stub Server (Mock Server)

Deal generates a stub server that you can run it a test against it.

```go
package main

import (
	"log"
	"net"

	"github.com/faunists/deal-go-example/protogen/proto/example"
	"google.golang.org/grpc"
)

func main() {
	// Generated stub server
	stubServer := example.MyServiceStubServer{}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error opening the listener: %v", err)
	}
	defer func() { _ = listener.Close() }()

	grpcServer := grpc.NewServer()
	example.RegisterMyServiceServer(grpcServer, &stubServer)

	log.Printf("grpc server listening at %v", listener.Addr())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	grpcServer.GracefulStop()
}
```

### Validating contract with server

The first step is to implement our server, the below implementation is compliant with the presented contract:

```go
package server

import (
	"context"

	"github.com/faunists/deal-go-example/protogen/proto/example"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MyServer struct {
	example.UnimplementedMyServiceServer
}

func (svc *MyServer) MyMethod(
	_ context.Context,
	req *example.RequestMessage,
) (*example.ResponseMessage, error) {
	if req.RequestField == "ANOTHER_VALUE" {
		return nil, status.Error(codes.NotFound, "ANOTHER_VALUE NotFound")
	}

	return &example.ResponseMessage{ResponseField: 42}, nil
}
```

Now we can use the generated test function that will validate our implementation:

```go
package main_test

import (
	"context"
	"testing"

	"github.com/faunists/deal-go-example/protogen/proto/example"

	"github.com/faunists/deal-go-example/api/server"
	"google.golang.org/grpc"
)

func TestMyServiceContract(t *testing.T) {
	grpcServer := grpc.NewServer()
	example.RegisterMyServiceServer(grpcServer, &server.MyServer{})

	// Testing the implementation
	example.MyServiceContractTest(t, context.Background(), grpcServer)
}
```
