package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {// TODO: replace this
	var result model.User
	err := r.db.Where("email = ?", email).First(&result).Error
	if err != nil {
		return model.User{}, nil
	}
	return result, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	result := []model.UserTaskCategory{}
	r.db.Table("tasks").Select("users.id, users.fullname, users.email, tasks.title AS task, tasks.deadline, tasks.priority, tasks.status, categories.name AS category").Joins("inner join categories on tasks.category_id = categories.id").Joins("inner join users on tasks.user_id = users.id").Scan(&result)
	return result, nil 
}
