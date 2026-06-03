package telegram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Sender struct {
	token  string
	chatID string
}

func New() *Sender {
	return &Sender{
		token:  os.Getenv("TELEGRAM_TOKEN"),
		chatID: os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

func (s *Sender) Send(text string) error {
	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		s.token,
	)

	values := url.Values{}
	values.Set("chat_id", s.chatID)
	values.Set("text", text)

	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"telegram status: %s, body: %s",
			resp.Status,
			string(body),
		)
	}

	return fmt.Errorf(
		"TELEGRAM TEST %s body=%s",
		resp.Status,
		string(body),
	)
}
