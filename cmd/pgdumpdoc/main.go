package main

import (
	"fmt"
	"log"

	"github.com/nDmitry/scripts/pkg/integrations"
)

type dumper interface {
	DumpDocker() error
}

type notifier interface {
	Notify(text string) error
}

func main() {
	Run(&integrations.PostgresDumper{}, &integrations.TelegramNotifier{})
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
