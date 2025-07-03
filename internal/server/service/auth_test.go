package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type fakeUserRepo struct {
	users  map[string]*entity.User
	autoID int64
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]*entity.User)}
}

func (r *fakeUserRepo) Create(_ context.Context, u *entity.User) error {
	if _, exists := r.users[u.Login]; exists {
		return errors.New("already exists")
	}
	r.autoID++
	u.ID = r.autoID
	r.users[u.Login] = u
	return nil
}

func (r *fakeUserRepo) FindByLogin(_ context.Context, login string) (*entity.User, error) {
	u, ok := r.users[login]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	repo := newFakeUserRepo()
	auth := NewAuthService(repo)

	// Регистрация
	err := auth.Register(ctx, "testuser", "secret")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Повторная регистрация → ошибка
	err = auth.Register(ctx, "testuser", "secret")
	if err == nil {
		t.Fatal("expected duplicate registration to fail")
	}

	// Успешный логин
	id, err := auth.Login(ctx, "testuser", "secret")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if id != 1 {
		t.Fatalf("unexpected user id: got %d, want 1", id)
	}

	// Неверный пароль
	_, err = auth.Login(ctx, "testuser", "wrong")
	if err == nil {
		t.Fatal("expected login to fail with wrong password")
	}

	// Несуществующий пользователь
	_, err = auth.Login(ctx, "nouser", "pass")
	if err == nil {
		t.Fatal("expected login to fail for unknown user")
	}
}
