Within the ROOT DIRECTORY i.e. gprc-auth-mongo

protoc --go_out=. proto/services.proto

protoc --go_out=. --go-grpc_out=. proto/services.proto


Specify the go_package path within the protofile.

option go_package = "./proto";