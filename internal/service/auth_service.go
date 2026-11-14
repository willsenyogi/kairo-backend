package service

import (
	"kairo-backend/internal/dto"
	"kairo-backend/internal/model"
	"kairo-backend/internal/repository"
)

type AuthService interface {
	GetAllUsers() ([]model.User, error)
	RegisterNewUser(user *model.User) error
	LoginValidateCreds(req *dto.LoginRequest) error
}

type authService struct {
	authRepo repository.AuthRepository
}

func InitAuthService(
	authRepo repository.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) GetAllUsers() ([]model.User, error) {
	users, err := s.authRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *authService) RegisterNewUser(user *model.User) error {
	return s.authRepo.RegisterNewUser(user)
}

func (s *authService) LoginValidateCreds(req *dto.LoginRequest) error {
	user := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	return s.authRepo.LoginCheck(&user)
}
