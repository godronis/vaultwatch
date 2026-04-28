package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramAlerter sends alerts to a Telegram chat via the Bot API.
type TelegramAlerter struct {
	botToken string
	chatID   string
	client   *http.Client
}

// NewTelegramAlerter creates a new TelegramAlerter.
// botToken is the Telegram bot token and chatID is the target chat/channel ID.
func NewTelegramAlerter(botToken, chatID string) (*TelegramAlerter, error) {
	if botToken == "" {
		return nil, fmt.Errorf("telegram bot token must not be empty")
	}
	if chatID == "" {
		return nil, fmt.Errorf("telegram chat ID must not be empty")
	}
	return &TelegramAlerter{
		botToken: botToken,
		chatID:   chatID,
		client:   &http.Client{},
	}, nil
}

// Send dispatches an alert message to the configured Telegram chat.
func (t *TelegramAlerter) Send(secret vault.SecretMeta) error {
	payload := map[string]string{
		"chat_id":    t.chatID,
		"text":       fmt.Sprintf("⚠️ VaultWatch Alert\nSecret: %s\nExpires: %s", secret.Path, secret.Expiration.Format("2006-01-02 15:04:05 UTC")),
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBase, t.botToken)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
