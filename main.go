package main

import (
	"camus/sanyavertolet/bot/pkg/bot"
	configuration "camus/sanyavertolet/bot/pkg/config"
	"camus/sanyavertolet/bot/pkg/database"
	"log"
)

func main() {
    config, err := configuration.LoadConfig("config.json")
    if err != nil {
        log.Fatal(err)
    }
    
    repo, err := database.InitDatabase(config.DatabaseFileName)
    if err != nil {
        log.Fatal(err)
    }
    
    bot.InitBot(config.TelegramBotToken, repo)
}
