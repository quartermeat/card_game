package scratch

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Andrew-peng/go-dalle2/dalle2"
	"github.com/joho/godotenv"
)

func RunDalle() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("DALLE_API_KEY")

	client, err := dalle2.MakeNewClientV1(apiKey)
	if err != nil {
		log.Fatalf("Error initializing client: %s", err)
	}

	resp, err := client.Create(
		context.Background(),
		"A black and white cat, kinda lengthy, mostly black. Has four white paws",
		dalle2.WithNumImages(1),
		dalle2.WithSize(dalle2.SMALL),
		dalle2.WithFormat(dalle2.URL),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, img := range resp.Data {
		fmt.Println("%s", img.Url)
	}
}
