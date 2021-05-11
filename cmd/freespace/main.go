package main

import (
	"fmt"
	"log"

	"github.com/nDmitry/scripts/pkg/diskstat"
	"github.com/nDmitry/scripts/pkg/notifier"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

const Threshold = 10

func main() {
	stat, err := diskstat.Get()

	if err != nil {
		log.Fatalln("Encountered and error while doing a syscall: %w", err)
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

	n := notifier.NotifierImpl{
		Notifier: &notifier.TelegramNotifier{},
	}

	if err = n.Notify(text); err != nil {
		log.Fatalln("Could not send the message: %w", err)
	}

	log.Println("Successfully sent a notification")
}
