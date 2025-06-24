package service

import (
	"context"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
	"github.com/lenarlenar/mygokeeper/internal/server/repo"
)

// RecordService реализует бизнес-логику работы с пользовательскими записями.
type RecordService struct {
	Repo repo.RecordRepository
}

// NewRecordService создаёт новый RecordService с указанным репозиторием.
func NewRecordService(r repo.RecordRepository) *RecordService {
	return &RecordService{Repo: r}
}

// Save сохраняет новую запись пользователя.
func (s *RecordService) Save(ctx context.Context, rec *entity.Record) error {
	return s.Repo.Save(ctx, rec)
}

// GetAll возвращает все записи пользователя.
func (s *RecordService) GetAll(ctx context.Context, userID int64) ([]entity.Record, error) {
	return s.Repo.GetAllByUser(ctx, userID)
}

// Update обновляет запись пользователя по ID.
func (s *RecordService) Update(ctx context.Context, rec *entity.Record) error {
	return s.Repo.Update(ctx, rec)
}

// Delete удаляет запись пользователя по ID.
func (s *RecordService) Delete(ctx context.Context, userID, id int64) error {
	return s.Repo.DeleteByID(ctx, userID, id)
}
