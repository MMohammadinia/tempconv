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

// gRPC: Celsius to Fahrenheit
func (s *server) CelsiusToFahrenheit(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	f := req.Value*9/5 + 32
	return &pb.TempResponse{Value: f}, nil
}

// gRPC: Fahrenheit to Celsius
func (s *server) FahrenheitToCelsius(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	c := (req.Value - 32) * 5 / 9
	return &pb.TempResponse{Value: c}, nil
}

func main() {
	// ---------------- gRPC Server ----------------
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

	// ---------------- HTTP Server ----------------
	// Root handler to avoid 404
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TempConv gRPC + HTTP server is running!")
	})

	// Celsius to Fahrenheit via HTTP
	http.HandleFunc("/celsius", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var c float64
		_, err := fmt.Sscanf(value, "%f", &c)
		if err != nil {
			http.Error(w, "Invalid value. Please provide a number.", http.StatusBadRequest)
			return
		}
		f := c*9/5 + 32
		fmt.Fprintf(w, "Fahrenheit: %.2f", f)
	})

	// Fahrenheit to Celsius via HTTP
	http.HandleFunc("/fahrenheit", func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")
		var f float64
		_, err := fmt.Sscanf(value, "%f", &f)
		if err != nil {
			http.Error(w, "Invalid value. Please provide a number.", http.StatusBadRequest)
			return
		}
		c := (f - 32) * 5 / 9
		fmt.Fprintf(w, "Celsius: %.2f", c)
	})

	log.Println("HTTP server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}