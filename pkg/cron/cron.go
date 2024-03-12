package cron

import (
	database "camus/sanyavertolet/bot/pkg/database/repository"
	"camus/sanyavertolet/bot/pkg/sheets"

	"github.com/robfig/cron"
)

const (
	EverySundayMidnightCronSpec = "0 0 0 * * SUN"
	EverySundayEveningCronSpec  = "0 0 22 * * SUN"

	EveryMiddayCronSpec = "0 0 12 * * *"
)

func InitCron(sheets *sheets.Sheets, repo *database.Repository) (*cron.Cron, error) {
	c := cron.New()

	err := c.AddFunc(EverySundayMidnightCronSpec, func() { sheets.SyncGames(repo) })
	if err != nil {
		return nil, err
	}

	c.Start()

	return c, nil
}
