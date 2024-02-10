# Rest and gRPC Gateway Setup

## Tools

- protoc
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway

---

## Basic gRPC Steps

### 1. Install tools

[tools.go](/tools/tools.go)

```sh
# install tools from tools.go
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/envoyproxy/protoc-gen-validate
```

### 2. Write new proto

[`message.proto`](/pb/message.proto)

### 3. Generate gRPC

```sh
protoc -I ./pb --go_out=./pkg/protobuf/ --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/protobuf/ --go-grpc_opt=paths=source_relative \
	pb/*.proto
```

### 4. Implement grpc stubs

Implement proto server interface type (pb.QuoteServiceServer)

[`service/quote.go`](/pkg/service/quote.go)

### 5. Create gRPC Server

[`main.go`](/main.go#L75-L94)

---

## Add REST Steps

### 6. Get proto libs

- [google/api](/pb/google/api)
- [validate](/pb/validate/validate.proto)

or get from source

- [google: googleapi](https://github.com/googleapis/googleapis/tree/master/google/api)
- [bufbuild: validate](https://github.com/bufbuild/protoc-gen-validate/tree/main/validate)

### 7. Add http tag for REST in proto

Add a `google.api.http` annotation

[`message.proto`](/pb/message.proto#L21-L31):

```diff
syntax = "proto3";
package messaging;
option go_package = "restgrpc/pb/messaging;messaging";

+import "google/api/annotations.proto";

...

service QuoteService {
- rpc Echo(StringMessage) returns (StringMessage) {}
+ rpc Echo(StringMessage) returns (StringMessage) {
+   option (google.api.http) = {
+     post: "/test/echo"
+     body: "*"
+   };
+ }
- rpc GetQuote(QuoteRequest) returns (QuoteReply) {}
+ rpc GetQuote(QuoteRequest) returns (QuoteReply) {
+   option (google.api.http) = {
+     get: "/quote/{author}"
+   };
+ }
}
```

### 8. Generate REST stubs

```sh
protoc -I ./pb --grpc-gateway_out ./pkg/protobuf \
    --grpc-gateway_opt paths=source_relative \
    pb/*.proto
```

### 9. Create HTTP server

[`main.go`](/main.go#L49-L73)

---
