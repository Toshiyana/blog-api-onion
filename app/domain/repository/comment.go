package repository

import (
	"context"
	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/comment"
	"myblog/app/domain/model/user"
)

// Comment : コメントリポジトリインターフェース
type Comment interface {
	Save(ctx context.Context, comment *comment.Comment) error
	FindByID(ctx context.Context, id string) (*comment.Comment, error)
	FindByBlogID(ctx context.Context, blogID blog.ID) ([]*comment.Comment, error)
	FindByUserID(ctx context.Context, userID user.ID) ([]*comment.Comment, error)
	Update(ctx context.Context, comment *comment.Comment) error
	Delete(ctx context.Context, id string) error
}
