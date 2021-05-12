package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nDmitry/scripts/pkg/archive"
	"github.com/nDmitry/scripts/pkg/integrations"
)

var directory = os.Getenv("DIRECTORY")
var outfile = os.Getenv("OUTFILE")
var endpoint = os.Getenv("S3_ENDPOINT")
var accessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
var accessKeySecret = os.Getenv("S3_ACCESS_KEY_SECRET")
var bucket = os.Getenv("S3_BUCKET")
var object = os.Getenv("S3_OBJECT")
var telegramToken = os.Getenv("TELEGRAM_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

type archiver interface {
	TarGzip(src string, out string) error
}

type uploader interface {
	Upload(filePath string, bucketName string, objectName string) error
}

type notifier interface {
	Notify(text string) error
}

func main() {
	Run(&archive.Tar{}, &integrations.TelegramNotifier{
		Token:  telegramToken,
		ChatID: telegramChatID,
	}, &integrations.S3{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	})
}

func Run(a archiver, n notifier, u uploader) {
	var err error

	if err = a.TarGzip(directory, outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not create the archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not create the archive: %v\n", err)
	}

	if err = u.Upload(outfile, bucket, object); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not upload the archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not upload the archive: %v\n", err)
	}

	if err = os.Remove(outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not remove uploaded archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not remove uploaded archive: %v\n", err)
	}

	log.Println("Successfully archived the directory")
}
