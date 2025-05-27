package repository

import (
	"context"
	"myblog/app/domain/model/user"
)

// User : ユーザーリポジトリインターフェース
type User interface {
	Save(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id string) error
}
