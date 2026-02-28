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

// gRPC server implementation
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
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTempConvServer(grpcServer, &server{})

	go func() {
		log.Println("gRPC running on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	http.HandleFunc("/celsius", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var c float64
		fmt.Sscanf(value, "%f", &c)
		f := c*9/5 + 32
		fmt.Fprintf(w, "Fahrenheit: %.2f", f)
	})

	http.HandleFunc("/fahrenheit", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var f float64
		fmt.Sscanf(value, "%f", &f)
		c := (f - 32) * 5 / 9
		fmt.Fprintf(w, "Celsius: %.2f", c)
	})

	log.Println("HTTP running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}