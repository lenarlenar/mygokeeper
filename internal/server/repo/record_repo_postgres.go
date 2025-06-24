package repo

import (
	"context"
	"database/sql"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type PostgresRecordRepo struct {
	db *sql.DB
}

func NewPostgresRecordRepo(db *sql.DB) *PostgresRecordRepo {
	return &PostgresRecordRepo{db: db}
}

func (r *PostgresRecordRepo) Save(ctx context.Context, rec *entity.Record) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO records (user_id, type, data, meta) VALUES ($1, $2, $3, $4)`,
		rec.UserID, rec.Type, rec.Data, rec.Meta)
	return err
}

func (r *PostgresRecordRepo) GetAllByUser(ctx context.Context, userID int64) ([]entity.Record, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, type, data, meta FROM records WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Record
	for rows.Next() {
		var r entity.Record
		err := rows.Scan(&r.ID, &r.UserID, &r.Type, &r.Data, &r.Meta)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func (r *PostgresRecordRepo) DeleteByID(ctx context.Context, userID, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM records WHERE id = $1 AND user_id = $2`,
		id, userID)
	return err
}

func (r *PostgresRecordRepo) Update(ctx context.Context, rec *entity.Record) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE records SET type = $1, data = $2, meta = $3
		 WHERE id = $4 AND user_id = $5`,
		rec.Type, rec.Data, rec.Meta, rec.ID, rec.UserID,
	)
	return err
}
