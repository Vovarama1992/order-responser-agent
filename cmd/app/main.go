package main

import (
	"log"
	"time"

	"fl-agent/config"
	"fl-agent/internal/gpt"
	"fl-agent/internal/service"
	"fl-agent/internal/source/fl"
	"fl-agent/internal/telegram"
)

func main() {
	config.LoadEnv()

	source := fl.NewSource()
	gptClient := gpt.New()
	telegramSender := telegram.New()

	watcher := service.NewWatcher(
		source,
		gptClient,
		telegramSender,
	)

	if err := watcher.RunOnce(); err != nil {
		log.Println(err)
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := watcher.RunOnce(); err != nil {
			log.Println(err)
		}
	}
}
