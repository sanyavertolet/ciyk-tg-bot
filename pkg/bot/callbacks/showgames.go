package bot

import (
	keyboards "camus/sanyavertolet/bot/pkg/bot/keyboards"
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type GamesGetter func()([]model.Game, error)

func ShowGames(
	bot *tgapi.BotAPI,
	update tgapi.Update,
	gamesGetter GamesGetter,
	text string,
	args []string,
	) {
	DeleteMessage(bot, update.FromChat().ID, update.CallbackQuery.Message.MessageID)

	page, err := strconv.Atoi(args[1])
	if err != nil {
		log.Panic(err)
	}

	games, err := gamesGetter()
	if err != nil {
		log.Panic(err)
	}

	message := tgapi.NewMessage(update.SentFrom().ID, text)
	message.ReplyMarkup = keyboards.GamesMenuKeyboard(games, page, args[0])
	if _, err := bot.Send(message); err != nil {
		log.Panic(err)
		return
	}
}

func ShowFutureGames(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	gameGetter := func() ([]model.Game, error) { return repo.FindFutureGames() }
	ShowGames(bot, update, gameGetter, "Предстоящие игры:", args)
}

func ShowPastGames(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	gameGetter := func() ([]model.Game, error) { return repo.FindPastGames() }
	ShowGames(bot, update, gameGetter, "Прошедшие игры:", args)
}

func ShowUserFutureGames(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	gameGetter := func() ([]model.Game, error) { return repo.FindFutureGamesByUserId(update.FromChat().ID) }
	ShowGames(bot, update, gameGetter, "Ваши предстоящие игры:", args)
}

func ShowUserPastGames(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	gameGetter := func() ([]model.Game, error) { return repo.FindPastGamesByUserId(update.FromChat().ID) }
	ShowGames(bot, update, gameGetter, "Ваши прошедшие игры:", args)
}