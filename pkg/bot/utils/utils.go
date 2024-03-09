package utils

import model "camus/sanyavertolet/bot/pkg/database/model"

func IsIdIn(id int64, users *[]model.User) bool {
	for _, user := range *users {
		if user.ID == id {
			return true
		}
	}
	return false
}