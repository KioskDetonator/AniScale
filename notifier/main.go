package main

import (
	"context"
	"log"
	"net"

	// Replace 'your-username' with your actual GitHub username
	pb "github.com/KioskDetonator/AniScale/proto"
	"google.golang.org/grpc"
)

// server is used to implement proto.NotifierServer
type server struct {
	pb.UnimplementedNotifierServer
}

// SendNotification implements proto.NotifierServer
func (s *server) SendNotification(ctx context.Context, in *pb.MangaUpdate) (*pb.NotificationResponse, error) {
	log.Printf("Received Manga Update: %s (Chapter %s) - URL: %s", in.GetTitle(), in.GetChapter(), in.GetUrl())
	
	// For now, we just return a successful response
	return &pb.NotificationResponse{
		Success: true, 
		Message: "Notification received and logged!",
	}, nil
}

func main() {
	// 1. Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Create a new gRPC server instance
	s := grpc.NewServer()

	// 3. Register our service implementation with the server
	pb.RegisterNotifierServer(s, &server{})

	log.Printf("Notifier server listening at %v", lis.Addr())

	// 4. Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}