package service

import (
	"fmt"
	"log"
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
	log.Println("[WATCHER] cycle started")

	orders, err := w.source.Fetch()
	if err != nil {
		return err
	}

	log.Printf("[WATCHER] fetched=%d\n", len(orders))

	for _, order := range orders {
		log.Printf("[WATCHER] order=%s title=%s\n", order.ID, order.Title)

		if w.seen[order.ID] {
			log.Printf("[WATCHER] skip already seen=%s\n", order.ID)
			continue
		}

		w.seen[order.ID] = true

		log.Printf("[WATCHER] parse order=%s\n", order.ID)

		fullOrder, err := fl.ParseOrder(order)
		if err != nil {
			log.Printf("[WATCHER] parse error=%v\n", err)
			continue
		}

		log.Printf("[WATCHER] gpt request order=%s\n", order.ID)

		result, err := w.gpt.Review(
			fullOrder.Title,
			fullOrder.Budget,
			fullOrder.Description,
		)
		if err != nil {
			log.Printf("[WATCHER] gpt error=%v\n", err)
			continue
		}

		log.Printf(
			"[WATCHER] gpt response order=%s category=%s\n",
			order.ID,
			result.Category,
		)

		if !isAllowedCategory(result.Category) {
			log.Printf("[WATCHER] filtered order=%s\n", order.ID)
			continue
		}

		message := formatMessage(fullOrder.URL, result)

		log.Printf("[WATCHER] telegram send order=%s\n", order.ID)

		if err := w.telegram.Send(message); err != nil {
			log.Printf("[WATCHER] telegram error=%v\n", err)
			continue
		}

		log.Printf("[WATCHER] telegram sent order=%s\n", order.ID)
	}

	log.Println("[WATCHER] cycle finished")

	return nil
}

func isAllowedCategory(category string) bool {
	category = strings.ToLower(category)

	return strings.Contains(category, "благородного дона") ||
		strings.Contains(category, "дон с бодуна") ||
		strings.Contains(category, "дона с бодуна")
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
