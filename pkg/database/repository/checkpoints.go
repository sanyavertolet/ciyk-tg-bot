package database

import (
	model "camus/sanyavertolet/bot/pkg/database/model"
	"time"
)

func (repo *Repository) SaveCheckpoint(lastProcessedLineNumber int) (model.Checkpoint, error) {
	checkpoint := model.Checkpoint{Line: lastProcessedLineNumber, CreatedAt: time.Now()}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	err := repo.DB.Create(&checkpoint).Error
	return checkpoint, err
}

func (repo *Repository) GetLastCheckpoint() (model.Checkpoint, error) {
	var checkpoint = model.Checkpoint{Line: 1}
	err := repo.DB.
		Order("created_at DESC").
		FirstOrCreate(&checkpoint).Error
	return checkpoint, err
}
