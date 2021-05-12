package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nDmitry/scripts/pkg/diskstat"
	"github.com/nDmitry/scripts/pkg/integrations"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

const Threshold = 10

var telegramToken = os.Getenv("TELEGRAM_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

type notifier interface {
	Notify(text string) error
}

func main() {
	Run(&integrations.TelegramNotifier{
		Token:  telegramToken,
		ChatID: telegramChatID,
	})
}

func Run(n notifier) {
	stat, err := diskstat.Get()

	if err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not check the disk space: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Encountered and error while doing a syscall: %v\n", err)
	}

	all := stat.All / float64(GB)
	avail := stat.Avail / float64(GB)
	left := avail / all * 100

	if left > Threshold {
		log.Printf("Disk space is OK: %.1f%%\n", left)
		return
	}

	text := fmt.Sprintf(
		"Disk space is running low: %.1f%% available, which is %.1f GB out of %.1f GB",
		left,
		avail,
		all,
	)

	if err = n.Notify(text); err != nil {
		log.Fatalf("Could not send the message: %v", err)
	}

	log.Println("Successfully sent a notification")
}
