package main

import (
	"camus/sanyavertolet/bot/pkg/bot"
	configuration "camus/sanyavertolet/bot/pkg/config"
	"camus/sanyavertolet/bot/pkg/cron"
	"camus/sanyavertolet/bot/pkg/database"
	"camus/sanyavertolet/bot/pkg/sheets"
	"github.com/carlmjohnson/versioninfo"
	"log"
)

func main() {
	log.Printf("Version: %s", versioninfo.Short())
	config, err := configuration.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Config loading error: %v", err)
	}

	repo, err := database.InitDatabase(config.DatabaseFileName)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	sheetsServices, err := sheets.InitSheets(config.GoogleSheetsKeyFileName, config.GoogleSpreadsheetID)
	if err != nil {
		log.Fatalf("Google Sheets initialization error: %v", err)
	}

	sheetsServices.SyncGames(repo)

	cronService, err := cron.InitCron(sheetsServices, repo)
	if err != nil {
		log.Fatalf("Cron initialization error: %v", err)
	}

	if err := bot.InitBot(repo, cronService, config.TelegramBotToken); err != nil {
		log.Fatalf("Telegram bot initialization error: %v", err)
	}
}
