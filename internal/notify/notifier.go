package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Notifier is a small interface for sending notifications.
type Notifier interface {
	Notify(ctx context.Context, title, message string) error
}

// TelegramNotifier sends messages to a Telegram chat.
type TelegramNotifier struct {
	botToken string
	chatID   string
	client   *http.Client
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *TelegramNotifier) Notify(ctx context.Context, title, message string) error {
	// keep short per-notify timeout
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)
	payload := map[string]string{
		"chat_id": t.chatID,
		"text":    fmt.Sprintf("[%s] %s", title, message),
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("telegram notify failed: status %d", resp.StatusCode)
	}
	return nil
}
