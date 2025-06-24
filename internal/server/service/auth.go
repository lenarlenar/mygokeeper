package service

import (
	"context"
	"errors"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
	"github.com/lenarlenar/mygokeeper/internal/server/repo"
	"golang.org/x/crypto/bcrypt"
)

// AuthService предоставляет методы регистрации и логина пользователей.
type AuthService struct {
	Users repo.UserRepo
}

// NewAuthService создаёт новый AuthService с переданным UserRepository.
func NewAuthService(users repo.UserRepo) *AuthService {
	return &AuthService{Users: users}
}

// Register регистрирует нового пользователя с логином и паролем.
// Возвращает ошибку, если пользователь уже существует или при ошибке хеширования.
func (s *AuthService) Register(ctx context.Context, login, password string) error {
	_, err := s.Users.FindByLogin(ctx, login)
	if err == nil {
		return errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.Users.Create(ctx, &entity.User{
		Login:    login,
		Password: string(hashed),
	})
}

// Login аутентифицирует пользователя по логину и паролю.
// Возвращает ID пользователя при успехе или ошибку при неудаче.
func (s *AuthService) Login(ctx context.Context, login, password string) (int64, error) {
	user, err := s.Users.FindByLogin(ctx, login)
	if err != nil {
		return 0, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("invalid password")
	}

	return user.ID, nil
}
