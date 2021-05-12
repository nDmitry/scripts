package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramNotifier struct {
	Token  string
	ChatID int
}

func (tn *TelegramNotifier) Notify(text string) error {
	body := struct {
		ChatID                int    `json:"chat_id"`
		Text                  string `json:"text"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	}{
		ChatID:                tn.ChatID,
		Text:                  text,
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

	return nil
}
