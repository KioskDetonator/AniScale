package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// The same structure as the Scraper
type MangaUpdate struct {
	Title   string `json:"title"`
	Chapter string `json:"chapter"`
	Url     string `json:"url"`
}

var discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")

func sendToDiscord(u MangaUpdate) error {
	msg := map[string]string{
		"content": fmt.Sprintf("📢 **New Manga Update!**\n**%s** - Chapter %s\nRead here: %s", u.Title, u.Chapter, u.Url),
	}
	body, _ := json.Marshal(msg)

	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT_ID")

	// 1. Setup the Client using the same key file
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("aniscale-key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 2. Point to your subscription
	sub := client.Subscription("notifier-sub")

	log.Println("📥 Notifier is waiting for Pub/Sub messages...")

	// 3. Receive messages indefinitely
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var update MangaUpdate
		if err := json.Unmarshal(msg.Data, &update); err != nil {
			log.Printf("Could not decode message: %v", err)
			msg.Ack() // Ack anyway to remove bad message from queue
			return
		}

		log.Printf("Processing: %s", update.Title)

		// 4. Send to Discord
		if err := sendToDiscord(update); err != nil {
			log.Printf("Discord error: %v", err)
			msg.Nack() // Nack means "I failed, try again later"
		} else {
			log.Printf("Successfully notified Discord for %s", update.Title)
			msg.Ack() // Ack means "Success, delete message from cloud"
		}
	})

	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}
}