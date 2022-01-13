# Deal - Go

[![test](https://github.com/faunists/deal-go/actions/workflows/test.yaml/badge.svg)](https://github.com/faunists/deal-go/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/faunists/deal-go/branch/main/graph/badge.svg?token=qFlORZnn09)](https://codecov.io/gh/faunists/deal-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/faunists/deal-go)](https://goreportcard.com/report/github.com/faunists/deal-go)

## Introduction

This plugin allows us to write [Consumer-Driver Contracts](https://martinfowler.com/articles/consumerDrivenContracts.html) tests!

__Deal__ will generate some code for us:
- A Client to be used in the client side to mock the responses based on the contract
- A Stub Server to be used in the client side as the Client above, but you should run it as another application
- Server Test Function, you should pass your server implementation to the function and all the contracts will be validated against it

You can check out an example project [here](https://github.com/faunists/deal-go-example).

## Installation

Assuming that you are using [Go Modules](https://github.com/golang/go/wiki/Modules), it's
recommended to use a [tool dependency](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
in order to track your tools version:

```go
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

### Generating code

If you're using [buf](https://buf.build) just add the following entry and execute `buf generate` passing your contract file path:
```yaml
version: v1beta1
plugins:
  - name: go-deal
    out: protogen
    opt: paths=source_relative,contract-file=contract.json
```

> Disclaimer: You must be using `go-grpc` in order to make the things work

To use the generated client you can just import it from the generated module:
```go
package main

import "YOUR_PACKAGE_HERE/example"

func main() {
	  contractClient := example.MyServiceContractClient{}

	  // TODO: Add the rest of the example here
}
```
