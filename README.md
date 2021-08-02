# Deal - Go

[![test](https://github.com/faunists/deal-go/actions/workflows/test.yaml/badge.svg)](https://github.com/faunists/deal-go/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/faunists/deal-go/branch/main/graph/badge.svg?token=qFlORZnn09)](https://codecov.io/gh/faunists/deal-go)

## Introduction

__WE DO NOT SUPPORT THE SERVER SIDE YET__

This plugin allows us to write [Consumer-Driver Contracts](https://martinfowler.com/articles/consumerDrivenContracts.html) tests!

## Usage example

### Proto service

First you need a proto service and add an option to it using our plugin providing the contract path:
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

  option(deal.v1.contract.deal_contract) = {
    contract_file: "contract.json"
  };
}
```

### Contract file

After that you need to write the contract that should be respected, the contract is written as a JSON file.
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
              "responseField": "RETURN_VALUE"
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

If you're using [buf](https://buf.build) just add the following entry and execute `buf generate`:
```yaml
version: v1beta1
plugins:
  - name: go-deal
    out: protogen
    opt: paths=source_relative
```

> Disclaimer: You must be using `go-grpc` in order to make the things work

To use the generated client you can just import it from the generated module:
```go
import "YOUR_PACKAGE_HERE/example"

func main() {
	  contractClient := example.MyServiceContractClient{}

	  // TODO: Add the rest of the example here
}
```
