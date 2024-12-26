package main

import (
	"log"
	"os"
	"strings"

	tgClient "github.com/gleblug/library-bot/clients/telegram"
	"github.com/gleblug/library-bot/consumer/event_consumer"
	"github.com/gleblug/library-bot/events/telegram"
	"github.com/gleblug/library-bot/lib/e"
	"github.com/gleblug/library-bot/storage/files"
)

const (
	batchSize = 100
	host      = "api.telegram.org"
)

func main() {
	token, storagePath := mustVariables()
	admins := strings.Split(os.Getenv("ADMIN_USERNAMES"), ";")

	storage, err := files.New(storagePath)
	if err != nil {
		log.Fatal(e.Wrap("Can't create storage", err))
	}

	eventsProcessor := telegram.New(
		tgClient.New(host, token),
		storage,
		admins,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
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
