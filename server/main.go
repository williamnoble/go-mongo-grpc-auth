package main

import (
	"context"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"grpc-auth-mongo/proto"
	"grpc-auth-mongo/server/authService"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	timeout = "5s"
	dbURI = "mongodb://localhost:27017"
	dbName = "grpcDemo"
)


func main() {
	ctx, cancel := createDBContext(timeout)
	defer cancel()
	db, err := openDB(ctx)
	if err != nil {
		log.Fatal("failed to connect to mongodb")
	}

	service := &authService.AuthServer{
		DB: db,
	}

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, service)

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error creating a listener")
	}

	go func() {
		log.Fatal("starting gRPC server", server.Serve(listener).Error())
	}()

	log.Println("Initalized gRPC server")
	grpcWebServer := grpcweb.WrapServer(server)

	httpServer := &http.Server{
		Addr: ":9001",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Println("Initalized proxy server")
	log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())
}

func openDB(ctx context.Context) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		return &mongo.Database{}, nil
	}

	db := client.Database(dbName)
	return db, nil

}


func createDBContext(timeout string) (context.Context, context.CancelFunc) {
	t, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("unable to parse timeout")
	}
	return context.WithTimeout(context.Background(), t )
}