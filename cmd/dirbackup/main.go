package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nDmitry/scripts/pkg/archive"
	"github.com/nDmitry/scripts/pkg/integrations"
)

var directory = flag.String("dir", "", "Directory path to backup")
var outfile = flag.String("o", "", "Backup file output path, will be deleted in a successful scenario")
var endpoints = flag.String("s3-endpoints", "", "S3 API endpoint URL (can be multiple, separate with a comma)")
var accessKeyIDs = flag.String("s3-key-ids", "", "S3 access key ID (can be multiple, separate with a comma)")
var accessKeySecrets = flag.String("s3-key-secrets", "", "S3 access key secret (can be multiple, separate with a comma)")
var buckets = flag.String("s3-buckets", "", "S3 bucket name (can be multiple, separate with a comma)")
var object = flag.String("s3-object", "", "S3 object key (path to uploaded file)")
var telegramToken = flag.String("telegram-token", "", "Telegram Bot API token")
var telegramChatID = flag.Int("telegram-chat-id", 0, "ID of a Telegram group to send messages in")

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
	flag.Parse()

	Run(&archive.Tar{}, &integrations.TelegramNotifier{
		Token:  *telegramToken,
		ChatID: *telegramChatID,
	}, &integrations.S3{
		Endpoints:        *endpoints,
		AccessKeyIDs:     *accessKeyIDs,
		AccessKeySecrets: *accessKeySecrets,
	})
}

func Run(a archiver, n notifier, u uploader) {
	var err error

	if err = a.TarGzip(*directory, *outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not create the archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not create the archive: %v\n", err)
	}

	if err = u.Upload(*outfile, *buckets, *object); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not upload the archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not upload the archive: %v\n", err)
	}

	if err = os.Remove(*outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not remove uploaded archive: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not remove uploaded archive: %v\n", err)
	}

	log.Println("Successfully archived the directory")
}
