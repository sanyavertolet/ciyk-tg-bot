package database

import (
	model "camus/sanyavertolet/bot/pkg/database/model"

	"time"
)

func (repo *Repository) CreateGame(name string, place string, date time.Time, maxPlayers int) (model.Game, error) {
	game := model.Game{Name: name, Date: date, MaxPlayers: maxPlayers, Place: place, IsRegistrationOpen: true}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	err := repo.DB.Create(&game).Error
	return game, err
}

func (repo *Repository) CreateGames(games []model.Game) error {
	if len(games) == 0 {
		return nil
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	return repo.DB.Create(games).Error
}

func (repo *Repository) FindFutureGames() ([]model.Game, error) {
	var games []model.Game
	currentTime := time.Now()
	err := repo.DB.
		Where("date > ?", currentTime).
		Order("date ASC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindPastGames() ([]model.Game, error) {
	var games []model.Game
	currentTime := time.Now()
	err := repo.DB.
		Where("date <= ?", currentTime).
		Order("date DESC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindFutureGamesByUserId(userId int64) ([]model.Game, error) {
	var games []model.Game
	currentTime := time.Now()
	err := repo.DB.
		Joins("JOIN registrations ON registrations.game_id = games.id").
		Where("registrations.user_id = ? AND games.date > ?", userId, currentTime).
		Group("games.id").
		Order("date ASC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindPastGamesByUserId(userId int64) ([]model.Game, error) {
	var games []model.Game
	currentTime := time.Now()
	err := repo.DB.
		Joins("JOIN registrations ON registrations.game_id = games.id").
		Where("registrations.user_id = ? AND games.date <= ?", userId, currentTime).
		Group("games.id").
		Order("date DESC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindGameById(id uint) (model.Game, error) {
	var game model.Game
	err := repo.DB.
		Preload("Registrations").
		Where("id = ?", id).
		First(&game).Error
	return game, err
}

func (repo *Repository) FindGameByNameAndDate(name string, date time.Time) (model.Game, error) {
	var game model.Game
	err := repo.DB.Where(model.Game{Name: name, Date: date}).Find(&game).Error
	return game, err
}

func (repo *Repository) OpenRegistrationForGames() ([]model.Game, error) {
	var games []model.Game
	now := time.Now()
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	err := repo.DB.
		Model(&model.Game{}).
		Where(&model.Game{IsRegistrationOpen: false}).
		Where("date >= ? AND date <= ?", now, now.AddDate(0, 0, 7)).
		Update("is_registration_open", true).
		Order("date ASC").
		Find(&games).Error
	return games, err
}

func (repo *Repository) FindNextWeekGames() ([]model.Game, error) {
	var games []model.Game
	now := time.Now()
	err := repo.DB.
		Model(&model.Game{}).
		Preload("Users").
		Where("date >= ? AND date <= ?", now, now.AddDate(0, 0, 7)).
		Where("is_registration_open", true).
		Order("date ASC").
		Find(&games).Error

	return games, err
}
