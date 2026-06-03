package service

import (
	"fmt"
	"strings"

	"fl-agent/internal/gpt"
	"fl-agent/internal/source/fl"
	"fl-agent/internal/telegram"
)

type Watcher struct {
	source   *fl.Source
	gpt      *gpt.Client
	telegram *telegram.Sender

	seen map[string]bool
}

func NewWatcher(
	source *fl.Source,
	gptClient *gpt.Client,
	telegramSender *telegram.Sender,
) *Watcher {
	return &Watcher{
		source:   source,
		gpt:      gptClient,
		telegram: telegramSender,
		seen:     make(map[string]bool),
	}
}

func (w *Watcher) RunOnce() error {
	orders, err := w.source.Fetch()
	if err != nil {
		return err
	}

	limit := 5
	if len(orders) < limit {
		limit = len(orders)
	}

	for _, order := range orders[:limit] {
		if w.seen[order.ID] {
			continue
		}

		w.seen[order.ID] = true

		fullOrder, err := fl.ParseOrder(order)
		if err != nil {
			return err
		}

		result, err := w.gpt.Review(
			fullOrder.Title,
			fullOrder.Budget,
			fullOrder.Description,
		)
		if err != nil {
			return err
		}

		if !isAllowedCategory(result.Category) {
			continue
		}

		message := formatMessage(fullOrder.URL, result)

		if err := w.telegram.Send(message); err != nil {
			return err
		}
	}

	return nil
}

func isAllowedCategory(category string) bool {
	category = strings.ToLower(category)

	return strings.Contains(category, "благородного дона") ||
		strings.Contains(category, "с бодуна")
}

func formatMessage(orderURL string, result *gpt.ReviewResult) string {
	return fmt.Sprintf(
		"Заказ:\n%s\n\nКатегория:\n%s\n\nСрок:\n%d\n\nЦена:\n%d\n\nОтклик:\n%s",
		orderURL,
		result.Category,
		result.Days,
		result.Price,
		result.Reply,
	)
}
