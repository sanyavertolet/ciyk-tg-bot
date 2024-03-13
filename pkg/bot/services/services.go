package bot

import (
	"camus/sanyavertolet/bot/pkg/bot/notifiers"
	"camus/sanyavertolet/bot/pkg/bot/utils"
	model "camus/sanyavertolet/bot/pkg/database/model"
	database "camus/sanyavertolet/bot/pkg/database/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func AddUser(bot *tgbotapi.BotAPI, repo *database.Repository, id int64, userTag string, userName string) {
	if _, err := repo.CreateUser(id, userTag, userName); err != nil {
		log.Panic(err)
		return
	}
	utils.ShowMenu(bot, id)
}

func AddGame(repo *database.Repository, message string) *model.Game {
	args := strings.Fields(message)
	if len(args) < 3 || len(args) > 4 {
		log.Panic("Not enough arguments")
		return nil
	}

	maxPlayers := 9
	if len(args) == 4 {
		atoi, err := strconv.Atoi(args[3])
		if err != nil {
			log.Panic("Wrong maxPlayers format")
			return nil
		}
		maxPlayers = atoi
	}

	name := args[0]
	place := args[1]

	date, err := time.Parse("02.01.06", args[2])
	if err != nil {
		log.Panic("Could not parse date, format: DD.MM.YY")
		return nil
	}

	game, err := repo.CreateGame(name, place, date, maxPlayers)

	if err != nil {
		log.Panic("Could not create game")
		return nil
	}
	return &game
}

func AddGameAndNotifyEveryone(bot *tgbotapi.BotAPI, repo *database.Repository, message string) *model.Game {
	game := AddGame(repo, message)
	notifiers.NotifyGameAdded(bot, repo, *game)
	return game
}

func sign(repo *database.Repository, userId int64, game model.Game) *model.Registration {
	counter, err := repo.CountUsersForGame(game.ID)
	if err != nil {
		log.Printf("Couldn't count users for game %s: %v", game.Name, err)
		return nil
	}

	log.Print("Creating registration")
	registration, err := repo.CreateRegistration(userId, game.ID, counter >= int64(game.MaxPlayers))
	if err != nil {
		log.Panicf("Couldn't create registration for game %s: %v", game.Name, err)
		return nil
	}

	return &registration
}

func SignViaCommand(repo *database.Repository, message string, userId int64) *model.Registration {
	args := strings.Fields(message)
	if len(args) < 2 {
		log.Panic("Not enough arguments: need game name and date (DD.MM.YY)")
		return nil
	} else if len(args) > 2 {
		log.Panic("Too many arguments: need game name and date (DD.MM.YY)")
		return nil
	}

	date, err := time.Parse("02.01.06", args[1])
	if err != nil {
		log.Panic("Wrong date format: DD.MM.YY is required")
		return nil
	}

	game, err := repo.FindGameByNameAndDate(args[0], date)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return sign(repo, userId, game)
}

func SignByGameId(repo *database.Repository, userId int64, gameId uint) *model.Registration {
	game, err := repo.FindGameById(gameId)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return sign(repo, userId, game)
}

func UnsignByGameId(repo *database.Repository, userId int64, gameId uint) {
	if err := repo.DeleteRegistration(userId, gameId); err != nil {
		log.Panic(err)
		return
	}
}

func SignFromReserve(repo *database.Repository, userId int64, gameId uint) bool {
	return repo.ChangeIsQueuingByUserIdAndGameId(userId, gameId, true) == nil
}
