package main

import (
	"camus/sanyavertolet/bot/pkg/bot"
	configuration "camus/sanyavertolet/bot/pkg/config"
	"camus/sanyavertolet/bot/pkg/cron"
	"camus/sanyavertolet/bot/pkg/database"
	"camus/sanyavertolet/bot/pkg/sheets"
	"log"
)

func main() {
	config, err := configuration.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Config loading error: %v", err)
	}

	repo, err := database.InitDatabase(config.DatabaseFileName)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	sheetsServices, err := sheets.InitSheets(config.GoogleSheetsKeyFileName)
	if err != nil {
		log.Fatalf("Google Sheets initialization error: %v", err)
	}

	cronService, err := cron.InitCron(sheetsServices, repo)
	if err != nil {
		log.Fatalf("Cron initialization error: %v", err)
	}

	sheetsServices.SyncGames(repo)

	if err := bot.InitBot(repo, cronService, config.TelegramBotToken); err != nil {
		log.Fatalf("Telegram bot initialization error: %v", err)
	}
}
