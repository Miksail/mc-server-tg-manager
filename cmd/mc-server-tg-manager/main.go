package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"

	"mc-server-tg-manager/internal/client"
	"mc-server-tg-manager/internal/model"
	"mc-server-tg-manager/internal/service"
)

func main() {
	log.Println("Starting mc-server-tg-manager")
	token := os.Getenv("TELEGRAM_TOKEN")
	rconHost := os.Getenv("RCON_HOST")
	rconPort := os.Getenv("RCON_PORT")
	rconPassword := os.Getenv("RCON_PASSWORD")
	containerName := os.Getenv("CONTAINER_NAME")

	rconClient, err := client.WaitForRCON(fmt.Sprintf("%s:%s", rconHost, rconPort), rconPassword)
	if err != nil {
		log.Fatal("Cannot connect to RCON:", err)
	}

	dockerClient, err := client.NewDockerClient(containerName)
	if err != nil {
		log.Fatal("Docker error:", err)
	}

	monitor := service.NewServerMonitor(
		rconClient,
		dockerClient,
		1*time.Minute,  // check interval
		15*time.Minute, // timout for shutdown
	)
	log.Println("Starting monitoring")
	monitor.Start()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to connect to tg api: %w", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID

		switch update.Message.Text {
		case "/start":
			status, _ := dockerClient.Status()
			if status == model.ServerStatusRunning {
				bot.Send(tgbotapi.NewMessage(chatID, "Server already running"))
				continue
			}

			err = dockerClient.Start()
			if err != nil {
				log.Printf("failed to start server: %s\n", err)
				bot.Send(tgbotapi.NewMessage(chatID, "Failed to start server"))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Server starting..."))
			}

		case "/status":
			status, _ := dockerClient.Status()

			if status == model.ServerStatusStopped {
				bot.Send(tgbotapi.NewMessage(chatID, "Server stopped"))
				continue
			}

			_, err = rconClient.ListPlayers()
			if err != nil {
				log.Printf("failed to list players: %s\n", err)
				bot.Send(tgbotapi.NewMessage(chatID, "Server starting..."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Server running"))
			}

		case "/stop":
			status, err := dockerClient.Status()
			if err != nil {
				log.Printf("failed to get server status: %s\n", err)
				bot.Send(tgbotapi.NewMessage(chatID, "Failed to stop server"))
				continue
			}
			if status != model.ServerStatusRunning {
				bot.Send(tgbotapi.NewMessage(chatID, "Server already stopped"))
				continue
			}

			err = rconClient.StopServer()
			if err != nil {
				log.Printf("failed to stop server: %s\n", err)
				bot.Send(tgbotapi.NewMessage(chatID, "Failed to stop server"))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Stopping server..."))
			}
		}
	}
}
