package database

import (
	. "camus/sanyavertolet/bot/pkg/database/model"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) CreateUser(id int64, telegramTag string) (User, error) {
	user := User{ID: id, Tag: telegramTag}
	err := repo.DB.Save(&user).Error
	return user, err
}

func (repo *Repository) CreateGame(name string, place string, date time.Time, maxPlayers int) (Game, error) {
	game := Game{Name: name, Date: date, MaxPlayers: maxPlayers, Place: place}
	err := repo.DB.Create(&game).Error
	return game, err
}

func (repo *Repository) CreateRegistration(userId int64, gameId uint, isQueuing bool) (Registration, error) {
	registration := Registration{UserId: userId, GameId: gameId, IsQueuing: isQueuing}
	if err := repo.DB.FirstOrCreate(&registration).Error; err != nil {
		return registration, err
	}
	err := repo.DB.Preload("Game").First(&registration, registration.ID).Error
	return registration, err
}

func (repo *Repository) DeleteRegistration(userId int64, gameId uint) error {
	err := repo.DB.
		Where("user_id = ? AND game_id = ?", userId, gameId).
		Delete(&Registration{}).Error
	return err
}

func (repo *Repository) FindAllUserIds() ([]int64, error) {
	var ids []int64
	err := repo.DB.Model(&User{}).Select("id").Find(&ids).Error
	return ids, err
}

func (repo *Repository) FindUserById(id int64) (User, error) {
	var user User
	err := repo.DB.Where(User{ID: id}).Find(&user).Error
	return user, err
}

func (repo *Repository) FindFutureGames() ([]Game, error) {
	var games []Game
	currentTime := time.Now()
	err := repo.DB.
		Where("date > ?", currentTime).
		Order("date ASC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindPastGames() ([]Game, error) {
	var games []Game
	currentTime := time.Now()
	err := repo.DB.
		Where("date <= ?", currentTime).
		Order("date DESC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindFutureGamesByUserId(userId int64) ([]Game, error) {
	var games []Game
	currentTime := time.Now()
	err := repo.DB.
		Joins("JOIN registrations ON registrations.game_id = games.id").
		Where("registrations.user_id = ? AND games.date > ? AND registrations.deleted_at IS NULL", userId, currentTime).
		Group("games.id").
		Order("date ASC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindPastGamesByUserId(userId int64) ([]Game, error) {
	var games []Game
	currentTime := time.Now()
	err := repo.DB.
		Joins("JOIN registrations ON registrations.game_id = games.id AND registrations.deleted_at IS NULL").
		Where("registrations.user_id = ? AND games.date <= ?", userId, currentTime).
		Group("games.id").
		Order("date DESC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindGameById(id uint) (Game, error) {
	var game Game
	err := repo.DB.
		Preload("Registrations").
		Where("id = ?", id).
		First(&game).Error
	return game, err
}

func (repo *Repository) FindUsersByGameIdOrderedByRegistrationTime(gameId uint) ([]User, error) {
	var users []User
	err := repo.DB.
		Joins("JOIN registrations ON registrations.user_id = users.id AND registrations.game_id = ?", gameId).
		Where("registrations.deleted_at IS NULL").
		Order("registrations.updated_at ASC").
		Find(&users).Error
	return users, err
}

func (repo *Repository) FindGameByNameAndDate(name string, date time.Time) (Game, error) {
	var game Game
	err := repo.DB.Where(Game{Name: name, Date: date}).Find(&game).Error
	return game, err
}

func (repo *Repository) FindFirstQueuingUserByGameId(gameId uint) (User, Game, error) {
	var user User
	var game Game
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", gameId).First(&game).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&Registration{}).
			Where("game_id = ? AND is_queuing = ?", gameId, false).
			Count(&count).Error; err != nil {
			return err
		}

		if count < int64(game.MaxPlayers) {
			var registration Registration
			if err := tx.
				Preload("Game").
				Where(&Registration{GameId: gameId, IsQueuing: true}).
				Order("created_at ASC").
				First(&registration).Error; err != nil {
				return err
			}
			user = registration.User
			return nil
		}
		return fmt.Errorf("no queuing user found under the specified conditions")
	})
	return user, game, err
}

func (repo *Repository) ChangeIsQueuingByUserIdAndGameId(userId int64, gameId uint, isQueuing bool) error {
	if err := repo.DB.
		Model(&Registration{}).
		Where(&Registration{GameId: gameId, UserId: userId}).
		Update("is_queuing", isQueuing).
		Error; err != nil {
			log.Print(err)
			return err
	}
	return nil
}

func (repo *Repository) CountUsersForGame(gameId uint) (int64, error) {
	var counter int64
	err := repo.DB.
		Model(Registration{}).
		Where(Registration{GameId: gameId, IsQueuing: false}).
		Count(&counter).Error
	return counter, err
}
