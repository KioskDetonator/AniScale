package main

import (
	"context"
	"log"
	"time"

	pb "github.com/KioskDetonator/AniScale/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotifierClient(conn)

	// 2. Prepare a fake manga update
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SendNotification(ctx, &pb.MangaUpdate{
		Title:   "One Piece",
		Chapter: "1111",
		Url:     "https://example.com/manga/one-piece-1111",
	})

	if err != nil {
		log.Fatalf("could not send notification: %v", err)
	}

	// 3. Print the server's response
	log.Printf("Server Response: %s (Success: %v)", r.GetMessage(), r.GetSuccess())
}