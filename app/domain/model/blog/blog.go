package blog

import (
	"errors"
	"time"

	"myblog/app/domain/model/user"

	"github.com/google/uuid"
)

// Blog : ブログエンティティ
type Blog struct {
	id        ID
	userID    user.ID
	title     string
	content   string
	createdAt time.Time
	updatedAt time.Time
}

// NewBlog : ブログの生成
func NewBlog(userID user.ID, title, content string) (*Blog, error) {
	if title == "" {
		return nil, errors.New("タイトルが空です")
	}
	if content == "" {
		return nil, errors.New("コンテンツが空です")
	}

	id, err := NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Blog{
		id:        *id,
		userID:    userID,
		title:     title,
		content:   content,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct : ブログの再構築（DBからの読み込み時など）
func Reconstruct(id string, userID user.ID, title, content string, createdAt, updatedAt time.Time) (*Blog, error) {
	blogID, err := NewID(id)
	if err != nil {
		return nil, err
	}

	return &Blog{
		id:        *blogID,
		userID:    userID,
		title:     title,
		content:   content,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// ID : IDの取得
func (b Blog) ID() ID {
	return b.id
}

// UserID : ユーザーIDの取得
func (b Blog) UserID() user.ID {
	return b.userID
}

// Title : タイトルの取得
func (b Blog) Title() string {
	return b.title
}

// Content : コンテンツの取得
func (b Blog) Content() string {
	return b.content
}

// CreatedAt : 作成日時の取得
func (b Blog) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt : 更新日時の取得
func (b Blog) UpdatedAt() time.Time {
	return b.updatedAt
}

// UpdateTitle : タイトルの更新
func (b *Blog) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("タイトルが空です")
	}
	b.title = title
	b.updatedAt = time.Now()
	return nil
}

// UpdateContent : コンテンツの更新
func (b *Blog) UpdateContent(content string) error {
	if content == "" {
		return errors.New("コンテンツが空です")
	}
	b.content = content
	b.updatedAt = time.Now()
	return nil
}
