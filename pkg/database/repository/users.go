package database

import model "camus/sanyavertolet/bot/pkg/database/model"

func (repo *Repository) CreateUser(id int64, telegramTag string) (model.User, error) {
	user := model.User{ID: id, Tag: telegramTag}
	err := repo.DB.Save(&user).Error
	return user, err
}

func (repo *Repository) FindAllUserIds() ([]int64, error) {
	var ids []int64
	err := repo.DB.Model(&model.User{}).Select("id").Find(&ids).Error
	return ids, err
}

func (repo *Repository) FindUserById(id int64) (model.User, error) {
	var user model.User
	err := repo.DB.Where(model.User{ID: id}).Find(&user).Error
	return user, err
}
