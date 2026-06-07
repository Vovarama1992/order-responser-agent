package main

import (
	"log"
	"time"

	"fl-agent/config"
	"fl-agent/internal/gpt"
	"fl-agent/internal/service"
	"fl-agent/internal/source/fl"
	"fl-agent/internal/source/freelance"
	"fl-agent/internal/source/upwork"
	"fl-agent/internal/storage"
	"fl-agent/internal/telegram"
)

func main() {
	config.LoadEnv()

	gptClient := gpt.New()
	telegramSender := telegram.New()

	store, err := storage.New("orders.db")
	if err != nil {
		log.Fatal(err)
	}

	watchers := []*service.Watcher{
		service.NewWatcher(
			fl.NewSource(),
			gptClient,
			telegramSender,
			store,
		),
		service.NewWatcher(
			freelance.NewSource(),
			gptClient,
			telegramSender,
			store,
		),
		service.NewWatcher(
			upwork.NewSource(),
			gptClient,
			telegramSender,
			store,
		),
	}
	runAll := func() {
		for _, watcher := range watchers {
			if err := watcher.RunOnce(); err != nil {
				log.Println(err)
			}
		}
	}

	runAll()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		runAll()
	}
}
