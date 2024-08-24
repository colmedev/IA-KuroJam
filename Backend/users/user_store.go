package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type store interface {
	Insert(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	ReadById(ctx context.Context, userId int64) (*User, error)
	ReadByGoogleId(ctx context.Context, googleId string) (*User, error)
}

type userStore struct {
	db *sqlx.DB
}

var (
	ErrInsertingUser  = errors.New("error inserting user")
	ErrUpdatingUser   = errors.New("error updating user")
	ErrEditConflict   = errors.New("error conflict user")
	ErrFetchingUser   = errors.New("error fetching user")
	ErrRecordNotFound = errors.New("error record not found")
)

func newUserStore(db *sqlx.DB) *userStore {
	return &userStore{
		db: db,
	}
}

func (us *userStore) Insert(ctx context.Context, u *User) error {
	query := `INSERT INTO users (clerk_id)
				VALUES ($1)
				RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := us.db.QueryRowxContext(ctx, query, u.ClerkId).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("%w: error executing query db:%w", ErrInsertingUser, err)
	}

	return nil
}

func (us *userStore) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET
		version = version + 1
		RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := us.db.QueryRowxContext(ctx, query, u).Scan(&u.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return fmt.Errorf("%w: error executing query db: %w", ErrUpdatingUser, err)

		}
	}

	return nil
}

func (us *userStore) ReadById(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, clerk_id, created_at, updated_at, version
			FROM users
			WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user User

	err := us.db.QueryRowxContext(ctx, query, id).Scan(&user.Id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: error fetching user (id: %v) db: %w", ErrRecordNotFound, id, err)
		default:
			return nil, fmt.Errorf("%w: error executing query db: %w", ErrFetchingUser, err)

		}
	}

	return &user, nil
}

func (us *userStore) ReadByGoogleId(ctx context.Context, googleId string) (*User, error) {
	query := `SELECT id, clerk_id, created_at, updated_at, version
			FROM users
			WHERE clerk_id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user User

	err := us.db.QueryRowxContext(ctx, query, googleId).StructScan(&user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: error fetching user (id: %v) db: %w", ErrRecordNotFound, googleId, err)
		default:
			return nil, fmt.Errorf("%w: error executing query db: %w", ErrFetchingUser, err)
		}
	}

	return &user, nil
}
