package repository

import (
	"errors"
	"kairo-backend/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetAllUsers() ([]model.User, error)
	RegisterNewUser(user *model.User) error //POST-type
	LoginCheck(user *model.User) error      //POST-type
}

type authRepository struct {
	db *gorm.DB
}

func InitAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *authRepository) RegisterNewUser(user *model.User) error {
	var existingUser model.User

	// Check duplicate email
	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("email is already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Create(user).Error
}

func (r *authRepository) LoginCheck(loginCreds *model.User) error {
	err := r.db.Where(
		"email = ? AND password = ?",
		loginCreds.Email,
		loginCreds.Password,
	).First(&model.User{}).Error

	if err != nil {
		return err
	}

	return nil
}
