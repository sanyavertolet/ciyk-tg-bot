package bot

import (
	callbacks "camus/sanyavertolet/bot/pkg/bot/callbacks"
	services "camus/sanyavertolet/bot/pkg/bot/services"
	"camus/sanyavertolet/bot/pkg/bot/utils"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func InitBot(token string, repo *database.Repository) {
	bot, err := tgapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			go ProcessCallbacks(bot, update, repo)
		}
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.IsCommand() {
				go ProcessCommands(bot, update, repo)
			}
		}
	}
}

func ProcessCommands(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository) {
	switch update.Message.Command() {
	case "start":
		services.AddUser(bot, repo, update.Message.From.ID, update.Message.From.UserName)
	case "add":
		services.AddGameAndNotifyEveryone(bot, repo, update.Message.CommandArguments())
	case "sign":
		services.SignViaCommand(repo, update.Message.CommandArguments(), update.Message.From.ID)
	case "menu":
		utils.ShowMenu(bot, update.Message.From.ID)
	}
}

func ProcessCallbacks(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository) {
	data := strings.Fields(update.CallbackQuery.Data)
	switch data[0] {
	case "sign":
		callbacks.SignCallback(bot, update, repo, data[1:])
	case "unsign":
		callbacks.UnsignCallback(bot, update, repo, data[1:])
	case "tomenu":
		callbacks.ShowMenuCallback(bot, update)
	case "futureGames":
		callbacks.ShowFutureGames(bot, update, repo, data)
	case "pastGames":
		callbacks.ShowPastGames(bot, update, repo, data)
	case "myFutureGames":
		callbacks.ShowUserFutureGames(bot, update, repo, data)
	case "myPastGames":
		callbacks.ShowUserPastGames(bot, update, repo, data)
	case "game":
		callbacks.ShowGame(bot, update, repo, data[1:])
	case "willCome":
		callbacks.WillCome(bot, update, repo, data[1:])
	case "wontCome":
		callbacks.WontCome(bot, update, repo, data[1:])
	}
}
