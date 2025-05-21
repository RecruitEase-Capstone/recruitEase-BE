package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/model"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	"github.com/jmoiron/sqlx"
)

const (
	InsertUserQuery = `INSERT INTO users(id, name, email, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6)`

	GetUserByIdQuery = `SELECT * FROM users WHERE id = $1`

	GetUserByEmailQuery = `SELECT * FROM users WHERE email = $1`
)

type IAuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) IAuthRepository {
	return &AuthRepository{db}
}

func (a *AuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	result, err := a.db.ExecContext(ctx, InsertUserQuery,
		user.ID, user.Name, user.Email, user.Password, user.UpdatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return customErr.ErrRowsAffected
	}

	return nil
}

func (a *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := a.db.GetContext(ctx, &user, GetUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}
