package main

import (
	"log"
	"os"
)

const (
	batchSize = 100
	host      = "api.telegram.org"
)

func main() {
	// token, storagePath := mustVariables()

	// eventsProcessor := telegram.New(
	// 	tgClient.New(host, token),
	// 	files.New(storagePath),
	// )

	log.Print("service started")

	// consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	// if err := consumer.Start(); err != nil {
	// 	log.Fatal("service is stopped", err)
	// }
}

func mustVariables() (string, string) {
	token := os.Getenv("TELEGRAM_API_KEY")
	storagePath := os.Getenv("LIBRARY_BOT_STORAGE")

	if token == "" {
		log.Fatal("TELEGRAM_API_KEY is not specified")
	}
	if storagePath == "" {
		log.Fatal("LIBRARY_BOT_STORAGE is not specified")
	}

	return token, storagePath
}
