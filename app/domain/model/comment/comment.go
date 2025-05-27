package comment

import (
	"errors"
	"time"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/user"

	"github.com/google/uuid"
)

// Comment : コメントエンティティ
type Comment struct {
	id        ID
	blogID    blog.ID
	userID    user.ID
	content   string
	createdAt time.Time
	updatedAt time.Time
}

// NewComment : コメントの生成
func NewComment(blogID blog.ID, userID user.ID, content string) (*Comment, error) {
	if content == "" {
		return nil, errors.New("コンテンツが空です")
	}

	id, err := NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Comment{
		id:        *id,
		blogID:    blogID,
		userID:    userID,
		content:   content,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct : コメントの再構築（DBからの読み込み時など）
func Reconstruct(id string, blogID blog.ID, userID user.ID, content string, createdAt, updatedAt time.Time) (*Comment, error) {
	commentID, err := NewID(id)
	if err != nil {
		return nil, err
	}

	return &Comment{
		id:        *commentID,
		blogID:    blogID,
		userID:    userID,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// ID : IDの取得
func (c Comment) ID() ID {
	return c.id
}

// BlogID : ブログIDの取得
func (c Comment) BlogID() blog.ID {
	return c.blogID
}

// UserID : ユーザーIDの取得
func (c Comment) UserID() user.ID {
	return c.userID
}

// Content : コンテンツの取得
func (c Comment) Content() string {
	return c.content
}

// CreatedAt : 作成日時の取得
func (c Comment) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt : 更新日時の取得
func (c Comment) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateContent : コンテンツの更新
func (c *Comment) UpdateContent(content string) error {
	if content == "" {
		return errors.New("コンテンツが空です")
	}
	c.content = content
	c.updatedAt = time.Now()
	return nil
}
