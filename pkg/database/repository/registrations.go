package database

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func (repo *Repository) CreateRegistration(userId int64, gameId uint, isQueuing bool) (model.Registration, error) {
	registration := model.Registration{UserId: userId, GameId: gameId, IsQueuing: isQueuing}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if err := repo.DB.FirstOrCreate(&registration).Error; err != nil {
		return registration, err
	}
	err := repo.DB.Preload("Game").First(&registration, registration.ID).Error
	return registration, err
}

func (repo *Repository) DeleteRegistration(userId int64, gameId uint) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	err := repo.DB.
		Where("user_id = ? AND game_id = ?", userId, gameId).
		Delete(&model.Registration{}).Error
	return err
}

func (repo *Repository) FindFirstQueuingUserByGameId(gameId uint) (model.User, model.Game, error) {
	var user model.User
	var game model.Game
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", gameId).First(&game).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&model.Registration{}).
			Where("game_id = ? AND is_queuing = ?", gameId, false).
			Count(&count).Error; err != nil {
			return err
		}

		if count < int64(game.MaxPlayers) {
			var registration model.Registration
			if err := tx.
				Preload("Game").
				Where(&model.Registration{GameId: gameId, IsQueuing: true}).
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

func (repo *Repository) FindUsersByGameIdOrderedByRegistrationTime(gameId uint) ([]model.User, error) {
	var users []model.User
	err := repo.DB.
		Joins("JOIN registrations ON registrations.user_id = users.id AND registrations.game_id = ?", gameId).
		Where("registrations.deleted_at IS NULL").
		Order("registrations.updated_at ASC").
		Find(&users).Error
	return users, err
}

func (repo *Repository) ChangeIsQueuingByUserIdAndGameId(userId int64, gameId uint, isQueuing bool) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if err := repo.DB.
		Model(&model.Registration{}).
		Where(&model.Registration{GameId: gameId, UserId: userId}).
		Update("is_queuing", isQueuing).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *Repository) CountUsersForGame(gameId uint) (int64, error) {
	var counter int64
	err := repo.DB.
		Model(model.Registration{}).
		Where(model.Registration{GameId: gameId, IsQueuing: false}).
		Count(&counter).Error
	return counter, err
}

func (repo *Repository) FindTomorrowUserIdsAssosiatedWithGames() (map[int64]model.Game, error) {
	startOfNextDay := time.Now().AddDate(0, 0, 1)
	startOfNextDayLocal := time.Date(
		startOfNextDay.Year(),
		startOfNextDay.Month(),
		startOfNextDay.Day(),
		0,
		0,
		0,
		0,
		time.Local,
	)
	endOfNextDayLocal := startOfNextDayLocal.AddDate(0, 0, 1)

	var registrations []model.Registration
	err := repo.DB.
		Model(&model.Registration{}).
		Preload("Game").
		Joins("JOIN games ON games.id = registrations.game_id").
		Where("games.date >= ? AND games.date < ?", startOfNextDayLocal, endOfNextDayLocal).
		Where("is_queuing = ?", false).
		Find(&registrations).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int64]model.Game)
	for _, registration := range registrations {
		result[registration.UserId] = registration.Game
	}

	return result, nil
}
