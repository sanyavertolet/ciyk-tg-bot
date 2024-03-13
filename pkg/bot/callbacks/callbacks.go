package bot

import (
	keyboards "camus/sanyavertolet/bot/pkg/bot/keyboards"
	"camus/sanyavertolet/bot/pkg/bot/notifiers"
	services "camus/sanyavertolet/bot/pkg/bot/services"
	"camus/sanyavertolet/bot/pkg/bot/utils"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"fmt"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func SignCallback(
	bot *tgapi.BotAPI,
	update tgapi.Update,
	repo *database.Repository,
	args []string) {
	DeleteMessage(bot, update.FromChat().ID, update.CallbackQuery.Message.MessageID)

	userId := update.CallbackQuery.From.ID
	gameId, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		log.Panic(err)
	}

	registration := services.SignByGameId(repo, userId, uint(gameId))

	messageText := fmt.Sprintf("Вы успешно зарегистрировались на игру!\n\n%s", registration.Game.String())
	message := tgapi.NewMessage(userId, messageText)
	message.ReplyMarkup = keyboards.GoToMainMenuKeyboard()

	if _, err = bot.Send(message); err != nil {
		log.Print("Could not send message to user")
		return
	}
}

func UnsignCallback(
	bot *tgapi.BotAPI,
	update tgapi.Update,
	repo *database.Repository,
	args []string,
) {
	DeleteMessage(bot, update.FromChat().ID, update.CallbackQuery.Message.MessageID)

	userId := update.FromChat().ID
	gameId, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		log.Panic(err)
	}

	services.UnsignByGameId(repo, userId, uint(gameId))
	notifiers.NotifyHeadOfQueue(bot, repo, uint(gameId))
	message := tgapi.NewMessage(userId, "Вы выписались с игры!")
	message.ReplyMarkup = keyboards.GoToMainMenuKeyboard()

	if _, err = bot.Send(message); err != nil {
		log.Print("Could not send message to user")
		return
	}
}

func ShowMenuCallback(bot *tgapi.BotAPI, update tgapi.Update) {
	DeleteMessage(bot, update.FromChat().ID, update.CallbackQuery.Message.MessageID)
	utils.ShowMenu(bot, update.CallbackQuery.From.ID)
}

func ShowGame(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	DeleteMessage(bot, update.FromChat().ID, update.CallbackQuery.Message.MessageID)
	gameId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Panic(err)
	}

	game, err := repo.FindGameById(uint(gameId))
	if err != nil {
		log.Panic(err)
	}

	users, err := repo.FindUsersByGameIdOrderedByRegistrationTime(game.ID)
	if err != nil {
		log.Panic(err)
	}

	game.Users = users

	message := tgapi.NewMessage(update.FromChat().ID, game.StringWithUsers())
	if utils.IsIdIn(update.FromChat().ID, &users) {
		message.ReplyMarkup = keyboards.UnsignKeyboard(game.ID)
	} else {
		message.ReplyMarkup = keyboards.SignKeyboard(game)
	}
	if _, err := bot.Send(message); err != nil {
		log.Panic(err)
	}
}

func DeleteMessage(bot *tgapi.BotAPI, chatId int64, messageId int) {
	deleteMessage := tgapi.NewDeleteMessage(chatId, messageId)
	if _, err := bot.Request(deleteMessage); err != nil {
		log.Print(err)
	}
}

func WillCome(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	userId := update.FromChat().ID
	DeleteMessage(bot, userId, update.CallbackQuery.Message.MessageID)

	gameId, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		log.Panic(err)
	}

	var messageText string
	if isOk := services.SignFromReserve(repo, userId, uint(gameId)); isOk {
		messageText = "Теперь вы в основе!"
	} else {
		messageText = "Что-то пошло не так!"
	}

	message := tgapi.NewMessage(userId, messageText)
	message.ReplyMarkup = keyboards.GoToMainMenuKeyboard()

	if _, err = bot.Send(message); err != nil {
		log.Print("Could not send message to user")
		return
	}
}

func WontCome(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository, args []string) {
	userId := update.FromChat().ID
	DeleteMessage(bot, userId, update.CallbackQuery.Message.MessageID)

	gameId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Panic(err)
	}
	if err := repo.DeleteRegistration(userId, uint(gameId)); err != nil {
		return
	}

	notifiers.NotifyHeadOfQueue(bot, repo, uint(gameId))
	message := tgapi.NewMessage(userId, "Вы выписались из игры.")
	message.ReplyMarkup = keyboards.GoToMainMenuKeyboard()
	if _, err := bot.Send(message); err != nil {
		log.Print("Could not send message to user")
		return
	}
}
