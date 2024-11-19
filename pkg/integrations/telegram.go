package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type TelegramNotifier struct {
	Token  string
	ChatID int
}

func (tn *TelegramNotifier) Notify(text string) error {
	hostname, err := os.Hostname()

	if err != nil {
		return fmt.Errorf("could not get the hostname: %w", err)
	}

	body := struct {
		ChatID                int    `json:"chat_id"`
		Text                  string `json:"text"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	}{
		ChatID:                tn.ChatID,
		Text:                  fmt.Sprintf("[%s]\n\n%s", hostname, text),
		DisableWebPagePreview: true,
	}

	json, err := json.Marshal(body)

	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tn.Token),
		"application/json",
		bytes.NewBuffer(json),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	log.Println("Telegram response status code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not send Telegram message, status: %d", resp.StatusCode)
	}

	return nil
}
