package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nDmitry/scripts/pkg/integrations"
	"github.com/nDmitry/scripts/pkg/tempmon"
)

var telegramToken = flag.String("telegram-token", "", "Telegram Bot API token")
var telegramChatID = flag.Int("telegram-chat-id", 0, "ID of a Telegram group to send messages in")
var threshold = flag.Float64("threshold", 60.0, "Temperature threshold considered not OK (not included)")
var sendOK = flag.Bool("send-ok", false, "Whether an OK result should be sent")

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
	cpuTemp, err := tempmon.GetAverageCPUTemp()

	if err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not check the CPU temperature: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Encountered and error while doing a syscall: %v\n", err)
	}

	ok := cpuTemp <= *threshold
	text := fmt.Sprintf("%.1f°C", cpuTemp)

	if ok {
		text = "CPU temperature is OK: " + text
	} else {
		text = "CPU temperature is running high: " + text
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
