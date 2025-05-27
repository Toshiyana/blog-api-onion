package repository

import (
	"context"
	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/user"
)

// Blog : ブログリポジトリインターフェース
type Blog interface {
	Save(ctx context.Context, blog *blog.Blog) error
	FindByID(ctx context.Context, id string) (*blog.Blog, error)
	FindByUserID(ctx context.Context, userID user.ID) ([]*blog.Blog, error)
	FindAll(ctx context.Context, offset, limit int) ([]*blog.Blog, error)
	Update(ctx context.Context, blog *blog.Blog) error
	Delete(ctx context.Context, id string) error
}
