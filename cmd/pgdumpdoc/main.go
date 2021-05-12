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
var endpoint = os.Getenv("S3_ENDPOINT")
var accessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
var accessKeySecret = os.Getenv("S3_ACCESS_KEY_SECRET")
var bucket = os.Getenv("S3_BUCKET")
var object = os.Getenv("S3_OBJECT")
var telegramToken = os.Getenv("TELEGRAM_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

type dumper interface {
	DumpDocker() error
}

type uploader interface {
	Upload(filePath string, bucketName string, objectName string) error
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
	}, &integrations.S3{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	})
}

func Run(d dumper, n notifier, u uploader) {
	var err error

	if err = d.DumpDocker(); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not backup the database: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not backup the database: %v\n", err)
	}

	if err = u.Upload(outfile, bucket, object); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not upload the database backup: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not upload the database backup: %v\n", err)
	}

	if err = os.Remove(outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not remove uploaded backup: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not remove uploaded backup: %v\n", err)
	}

	log.Println("Successfully dumped the database")
}
