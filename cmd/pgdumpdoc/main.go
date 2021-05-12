package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nDmitry/scripts/pkg/integrations"
)

var container = os.Getenv("CONTAINER")
var user = os.Getenv("USER")
var database = os.Getenv("DATABASE")
var outfile = os.Getenv("OUTFILE")
var telegramToken = os.Getenv("TELEGRAM_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

type dumper interface {
	DumpDocker() error
}

type notifier interface {
	Notify(text string) error
}

func main() {
	Run(&integrations.PostgresDumper{
		Container: container,
		User:      user,
		Database:  database,
		Outfile:   outfile,
	}, &integrations.TelegramNotifier{
		Token:  telegramToken,
		ChatID: telegramChatID,
	})
}

func Run(d dumper, n notifier) {
	err := d.DumpDocker()

	if err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not check the disk space: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not dump the database: %v\n", err)
	}

	log.Println("Successfully dumped the database")
}
