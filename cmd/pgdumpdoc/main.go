package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nDmitry/scripts/pkg/integrations"
)

var container = flag.String("container", "", "Docker container name running PostgreSQL database")
var user = flag.String("u", "", "User to run pg_dump under")
var database = flag.String("db", "", "Database name")
var outfile = flag.String("o", "", "Backup file output path, will be deleted in a successful scenario")
var endpoint = flag.String("s3-endpoint", "", "S3 API endpoint URL")
var accessKeyID = flag.String("s3-key-id", "", "S3 access key ID")
var accessKeySecret = flag.String("s3-key-secret", "", "S3 access key secret")
var bucket = flag.String("s3-bucket", "", "S3 bucket name")
var object = flag.String("s3-object", "", "S3 object key (path to uploaded file)")
var telegramToken = flag.String("telegram-token", "", "Telegram Bot API token")
var telegramChatID = flag.Int("telegram-chat-id", 0, "ID of a Telegram group to send messages in")

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
	flag.Parse()

	Run(&integrations.PostgresDumper{
		Container: *container,
		User:      *user,
		Database:  *database,
		Outfile:   *outfile,
	}, &integrations.TelegramNotifier{
		Token:  *telegramToken,
		ChatID: *telegramChatID,
	}, &integrations.S3{
		Endpoint:        *endpoint,
		AccessKeyID:     *accessKeyID,
		AccessKeySecret: *accessKeySecret,
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

	if err = u.Upload(*outfile, *bucket, *object); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not upload the database backup: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not upload the database backup: %v\n", err)
	}

	if err = os.Remove(*outfile); err != nil {
		if telegramErr := n.Notify(
			fmt.Sprintf("Could not remove uploaded backup: %v\n", err),
		); telegramErr != nil {
			log.Printf("Could not send an error notification: %v\n", err)
		}

		log.Fatalf("Could not remove uploaded backup: %v\n", err)
	}

	log.Println("Successfully dumped the database")
}
