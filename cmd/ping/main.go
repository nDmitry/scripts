package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nDmitry/scripts/pkg/integrations"
)

var timeout = flag.Int("timeout", 10, "Timeout in seconds after which a tested service is considered down")
var host = flag.String("host", "", "Host to check")
var port = flag.String("port", "", "Port to check")
var telegramToken = flag.String("telegram-token", "", "Telegram Bot API token")
var telegramChatID = flag.Int("telegram-chat-id", 0, "ID of a Telegram group to send messages in")

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
	isUp := up(*host, *port)
	text := fmt.Sprintf("Network target %s:%s seems online", *host, *port)

	if !isUp {
		text = fmt.Sprintf("Network target %s:%s is offline", *host, *port)
	}

	log.Println(text)

	if isUp {
		return
	}

	if err := n.Notify(text); err != nil {
		log.Fatalf("Could not send the message: %v", err)
	}

	log.Println("Successfully sent a notification")
}

func up(host string, port string) bool {
	conn, err := net.DialTimeout(
		"tcp",
		net.JoinHostPort(host, port),
		time.Duration(*timeout)*time.Second,
	)

	if err != nil {
		log.Println("Unsuccessfull ping: %w", err)

		return false
	}

	if conn != nil {
		defer conn.Close()

		return true
	}

	log.Println("Unsuccessfull ping: could not obtain a connection")

	return false
}
