package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type MangaDexResponse struct {
	Data []struct {
		Attributes    struct{ Chapter string `json:"chapter"` } `json:"attributes"`
		Relationships []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"relationships"`
	} `json:"data"`
}

type MangaUpdate struct {
	Title   string `json:"title"`
	Chapter string `json:"chapter"`
	Url     string `json:"url"`
}

func main() {
	// Load .env for local development; ignore error in production/Docker
	_ = godotenv.Load()

	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is required")
	}

	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("aniscale-key.json"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer client.Close()

	topic := client.Topic("manga-updates")
	log.Println("🚀 Scraper started. Publishing to topic: manga-updates")

	for {
		resp, err := http.Get("https://api.mangadex.org/chapter?limit=5&order[readableAt]=desc&includes[]=manga")
		if err != nil {
			log.Printf("MangaDex error: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		var result MangaDexResponse
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		for _, item := range result.Data {
			var mID string
			for _, r := range item.Relationships {
				if r.Type == "manga" { mID = r.ID; break }
			}

			update := MangaUpdate{
				Title:   fmt.Sprintf("Manga ID: %s", mID),
				Chapter: item.Attributes.Chapter,
				Url:     fmt.Sprintf("https://mangadex.org/title/%s", mID),
			}

			data, _ := json.Marshal(update)
			res := topic.Publish(ctx, &pubsub.Message{Data: data})
			id, err := res.Get(ctx)
			if err != nil {
				log.Printf("Publish failed: %v", err)
			} else {
				log.Printf("✅ Published: %s (ID: %s)", update.Title, id)
			}
		}
		time.Sleep(2 * time.Minute)
	}
}