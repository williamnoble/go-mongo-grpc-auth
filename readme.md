# Google gRPC with MongoDB and simple JWT Authentication

<br/>

### How to setup protocol buffers and gRPC

> I include detailed information as this is something I struggled with when I first encountered protocol buffers!

Within the **root directory** i.e. `gprc-auth-mongo`

```shell
# Generate the Proto file alone
$ protoc --go_out=. proto/services.proto

# Generate the GRPC file
$ protoc --go_out=. --go-grpc_out=. proto/services.proto
```

Within the `.proto` file itself, we set the **package** and **optional go_package**

```go
// services.proto
option go_package = "./proto" // ignore the error
package proto
```

### How to run this example

##### Server

```shell
$ go run server/main.go
```

##### Client

```shell
$ go run client/main.go
```

The client will output two JWT tokens `[AuthResponse]`
- for when the _User_ registers
- for when the _User_ logs in.