package utils

import model "camus/sanyavertolet/bot/pkg/database/model"

func IsIdIn(id int64, registrations []model.Registration) bool {
	for _, registration := range registrations {
		if registration.UserId == id {
			return true
		}
	}
	return false
}
