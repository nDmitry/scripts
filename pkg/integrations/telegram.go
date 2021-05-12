package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var telegramToken = os.Getenv("TELEGRAM_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

type TelegramNotifier struct{}

func (tn *TelegramNotifier) Notify(text string) error {
	chatID, err := strconv.Atoi(telegramChatID)

	if err != nil {
		return err
	}

	body := struct {
		ChatID                int    `json:"chat_id"`
		Text                  string `json:"text"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	}{
		ChatID:                chatID,
		Text:                  text,
		DisableWebPagePreview: true,
	}

	json, err := json.Marshal(body)

	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken),
		"application/json",
		bytes.NewBuffer(json),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
