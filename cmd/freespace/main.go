package main

import (
	"flag"
	"fmt"
	"log"

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

var telegramToken = flag.String("telegram-token", "", "Telegram Bot API token")
var telegramChatID = flag.Int("telegram-chat-id", 0, "ID of a Telegram group to send messages in")
var sendOK = flag.Bool("send-ok", true, "Whether an OK result should be sent")

type notifier interface {
	Notify(text string) error
}

func main() {
	flag.Parse()

	Run(&integrations.TelegramNotifier{
		Token:  *telegramToken,
		ChatID: *telegramChatID,
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
	ok := left > Threshold

	text := fmt.Sprintf(
		"%.1f%% available, which is %.1f GB out of %.1f GB",
		left,
		avail,
		all,
	)

	if ok {
		text = "Disk space is OK: " + text
	} else {
		text = "Disk space is running low: " + text
	}

	log.Println(text)

	if ok && !*sendOK {
		return
	}

	if err = n.Notify(text); err != nil {
		log.Fatalf("Could not send the message: %v", err)
	}

	log.Println("Successfully sent a notification")
}
