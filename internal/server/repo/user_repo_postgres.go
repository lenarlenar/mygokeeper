package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, u *entity.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2)`,
		u.Login, u.Password,
	)
	return err
}

func (r *PostgresUserRepo) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, login, password FROM users WHERE login = $1`,
		login,
	)

	var user entity.User
	if err := row.Scan(&user.ID, &user.Login, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
