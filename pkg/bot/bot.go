package bot

import (
	callbacks "camus/sanyavertolet/bot/pkg/bot/callbacks"
	"camus/sanyavertolet/bot/pkg/bot/inlinequeries"
	"camus/sanyavertolet/bot/pkg/bot/notifiers"
	services "camus/sanyavertolet/bot/pkg/bot/services"
	"camus/sanyavertolet/bot/pkg/bot/utils"
	"camus/sanyavertolet/bot/pkg/cron"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cronLib "github.com/robfig/cron"
	"log"
	"strings"
)

func InitBot(repo *database.Repository, cronService *cronLib.Cron, token string) error {
	bot, err := tgapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	if err := ScheduleRegistrationOpening(bot, repo, cronService); err != nil {
		return err
	}
	if err := ScheduleGameReminders(bot, repo, cronService); err != nil {
		return err
	}

	for update := range updates {
		if update.InlineQuery != nil {
			go ProcessInlineQuery(bot, update, repo)
		}
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

	return nil
}

func ProcessCommands(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository) {
	switch update.Message.Command() {
	case "start":
		services.AddUser(bot, repo, update.FromChat().ID, update.FromChat().UserName, update.FromChat().FirstName)
	case "add":
		services.AddGameAndNotifyEveryone(bot, repo, update.Message.CommandArguments())
	case "sign":
		services.SignViaCommand(repo, update.Message.CommandArguments(), update.FromChat().ID)
	case "menu":
		utils.ShowMenu(bot, update.FromChat().ID)
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

func ProcessInlineQuery(bot *tgapi.BotAPI, update tgapi.Update, repo *database.Repository) {
	if update.InlineQuery.Query == "" {
		inlinequeries.Hints(bot, update)
		return
	}
	data := strings.Fields(update.InlineQuery.Query)
	switch {
	case data[0] == "no":
		inlinequeries.AboutGame(bot, update, repo)
	}
}

func ScheduleRegistrationOpening(bot *tgapi.BotAPI, repo *database.Repository, cronService *cronLib.Cron) error {
	return cronService.AddFunc(cron.EverySundayEveningCronSpec, func() {
		games, err := repo.OpenRegistrationForGames()
		if err != nil {
			log.Printf("Could not open registration for games: %v", err)
		}
		notifiers.NotifyEverybodyGamesAdded(bot, repo, games)
		log.Printf("Opened %d games", len(games))
	})
}

func ScheduleGameReminders(bot *tgapi.BotAPI, repo *database.Repository, cronService *cronLib.Cron) error {
	return cronService.AddFunc(cron.EveryMiddayCronSpec, func() {
		userIdGameMap, err := repo.FindTomorrowUserIdsAssosiatedWithGames()
		if err != nil {
			log.Printf("Could not remind user about game: %v", err)
		}
		notifiers.NotifyUsersForTomorrowGames(bot, userIdGameMap)
		log.Printf("Notified %d users", len(userIdGameMap))
	})
}
