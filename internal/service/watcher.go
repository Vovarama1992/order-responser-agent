package service

import (
	"fmt"
	"log"
	"strings"

	"fl-agent/internal/gpt"
	"fl-agent/internal/source/fl"
	"fl-agent/internal/storage"
	"fl-agent/internal/telegram"
)

type Watcher struct {
	source   *fl.Source
	gpt      *gpt.Client
	telegram *telegram.Sender
	storage  *storage.Storage
}

func NewWatcher(
	source *fl.Source,
	gptClient *gpt.Client,
	telegramSender *telegram.Sender,
	storage *storage.Storage,
) *Watcher {
	return &Watcher{
		source:   source,
		gpt:      gptClient,
		telegram: telegramSender,
		storage:  storage,
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
		log.Printf("[WATCHER] order=%s source=%s title=%s\n", order.ID, order.Source, order.Title)

		exists, err := w.storage.Exists(order.Source, order.ID)
		if err != nil {
			log.Printf("[WATCHER] storage exists error=%v\n", err)
			continue
		}

		if exists {
			log.Printf("[WATCHER] skip already seen source=%s id=%s\n", order.Source, order.ID)
			continue
		}

		log.Printf("[WATCHER] parse order=%s source=%s\n", order.ID, order.Source)

		fullOrder, err := fl.ParseOrder(order)
		if err != nil {
			log.Printf("[WATCHER] parse error=%v\n", err)
			continue
		}

		log.Printf("[WATCHER] gpt request order=%s source=%s\n", order.ID, order.Source)

		result, err := w.gpt.Review(
			fullOrder.Title,
			fullOrder.Budget,
			fullOrder.Description,
		)
		if err != nil {
			log.Printf("[WATCHER] gpt error=%v\n", err)
			continue
		}

		if err := w.storage.Save(order.Source, order.ID, order.URL); err != nil {
			log.Printf("[WATCHER] storage save error=%v\n", err)
		}

		log.Printf(
			"[WATCHER] gpt response order=%s source=%s category=%s\n",
			order.ID,
			order.Source,
			result.Category,
		)

		if !isAllowedCategory(result.Category) {
			log.Printf("[WATCHER] filtered order=%s source=%s\n", order.ID, order.Source)
			continue
		}

		message := formatMessage(fullOrder.URL, result)

		log.Printf("[WATCHER] telegram send order=%s source=%s\n", order.ID, order.Source)

		if err := w.telegram.Send(message); err != nil {
			log.Printf("[WATCHER] telegram error=%v\n", err)
			continue
		}

		log.Printf("[WATCHER] telegram sent order=%s source=%s\n", order.ID, order.Source)
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
