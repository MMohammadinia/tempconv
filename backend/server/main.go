package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "tempconv/backend/server/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTempConvServer
}

func (s *server) CelsiusToFahrenheit(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	f := req.Value*9/5 + 32
	return &pb.TempResponse{Value: f}, nil
}

func (s *server) FahrenheitToCelsius(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	c := (req.Value - 32) * 5 / 9
	return &pb.TempResponse{Value: c}, nil
}

func main() {
	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTempConvServer(grpcServer, &server{})

	go func() {
		log.Println("gRPC server running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Simple HTTP endpoint for testing
	http.HandleFunc("/celsius", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var c float64
		fmt.Sscanf(value, "%f", &c)
		f := c*9/5 + 32
		fmt.Fprintf(w, "Fahrenheit: %.2f", f)
	})

	log.Println("HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}