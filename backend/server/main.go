package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

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
	// Use Railway's PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051" // fallback for local dev
	}

	// gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTempConvServer(s, &server{})

	go func() {
		log.Printf("gRPC running on :%s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Optional HTTP endpoint for testing
	http.HandleFunc("/celsius", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var c float64
		fmt.Sscanf(value, "%f", &c)
		f := c*9/5 + 32
		fmt.Fprintf(w, "Fahrenheit: %.2f", f)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // fallback for local dev
	}

	log.Printf("HTTP running on :%s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}