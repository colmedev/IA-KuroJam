package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/colmedev/IA-KuroJam/Backend/validator"
	"github.com/jmoiron/sqlx"
)

var (
	UserValidationError = errors.New("user data is not valid")
)

var AnonymousUser = User{}

func (u *User) IsAnonymous() bool {
	return u == &AnonymousUser
}

type Service interface {
	Insert(ctx context.Context, u *User) (error, map[string]string)
	Update(ctx context.Context, u *User) (error, map[string]string)
	FindById(ctx context.Context, userId int64) (*User, error)
	FindByGoogleId(ctx context.Context, userId string) (*User, error)
}

type UserService struct {
	store store
}

func NewService(db *sqlx.DB) (*UserService, error) {
	us := newUserStore(db)

	return &UserService{
		store: us,
	}, nil
}

// TODO: Finish user validation
func validateUser(v *validator.Validator, u *User) {
	// v.Check(u.Name != "", "name", "must be provided")
	// v.Check(len(u.Name) < 60, "name", "must not be more than 50 characters")
	// v.Check(validator.Matches(u.Email, validator.EmailRX), "email", "must be a valid email address")
}

func (us *UserService) Insert(ctx context.Context, u *User) (error, map[string]string) {
	v := validator.New()

	if validateUser(v, u); !v.Valid() {
		return UserValidationError, v.Errors
	}

	err := us.store.Insert(ctx, u)
	if err != nil {
		return fmt.Errorf("UserService: user creation failed: %w", err), nil
	}

	return nil, nil
}

func (us *UserService) Update(ctx context.Context, u *User) (error, map[string]string) {
	v := validator.New()

	if validateUser(v, u); !v.Valid() {
		return UserValidationError, v.Errors
	}

	err := us.store.Insert(ctx, u)
	if err != nil {
		return fmt.Errorf("UserService: user creation failed: %w", err), nil
	}

	return nil, nil
}

func (us *UserService) FindById(ctx context.Context, userId int64) (*User, error) {
	return nil, nil
}

func (us *UserService) FindByGoogleId(ctx context.Context, userId string) (*User, error) {
	return us.store.ReadByGoogleId(ctx, userId)
}
