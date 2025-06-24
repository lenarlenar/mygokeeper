package repo

import (
	"context"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type UserRepo interface {
	Create(ctx context.Context, u *entity.User) error
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
}
