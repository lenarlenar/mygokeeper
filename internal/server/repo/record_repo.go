package repo

import (
	"context"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type RecordRepository interface {
	Save(ctx context.Context, r *entity.Record) error
	GetAllByUser(ctx context.Context, userID int64) ([]entity.Record, error)
	DeleteByID(ctx context.Context, userID, id int64) error
	Update(ctx context.Context, rec *entity.Record) error
}
