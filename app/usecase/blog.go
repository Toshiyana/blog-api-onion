package usecase

import (
	"context"
	"errors"
	"fmt"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"
)

// BlogUsecase : ブログユースケースインターフェース
type BlogUsecase interface {
	CreateBlog(ctx context.Context, userID, title, content string) (*blog.Blog, error)
	GetBlogByID(ctx context.Context, id string) (*blog.Blog, error)
	GetBlogsByUserID(ctx context.Context, userID string) ([]*blog.Blog, error)
	GetAllBlogs(ctx context.Context, page, perPage int) ([]*blog.Blog, error)
	UpdateBlog(ctx context.Context, id, userID, title, content string) (*blog.Blog, error)
	DeleteBlog(ctx context.Context, id, userID string) error
}

// blogUsecase : ブログユースケースの実装
type blogUsecase struct {
	blogRepo repository.Blog
	userRepo repository.User
}

// NewBlogUsecase : ブログユースケースの生成
func NewBlogUsecase(blogRepo repository.Blog, userRepo repository.User) BlogUsecase {
	return &blogUsecase{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}
}

// CreateBlog : ブログの作成
func (b *blogUsecase) CreateBlog(ctx context.Context, userID, title, content string) (*blog.Blog, error) {
	// ユーザーの検証
	user, err := b.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %w", err)
	}

	// ブログの生成
	newBlog, err := blog.NewBlog(user.ID(), title, content)
	if err != nil {
		return nil, fmt.Errorf("ブログ作成エラー: %w", err)
	}

	// ブログの保存
	if err := b.blogRepo.Save(ctx, newBlog); err != nil {
		return nil, fmt.Errorf("ブログ保存エラー: %w", err)
	}

	return newBlog, nil
}

// GetBlogByID : IDによるブログ取得
func (b *blogUsecase) GetBlogByID(ctx context.Context, id string) (*blog.Blog, error) {
	blog, err := b.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ブログ取得エラー: %w", err)
	}
	return blog, nil
}

// GetBlogsByUserID : ユーザーIDによるブログ一覧取得
func (b *blogUsecase) GetBlogsByUserID(ctx context.Context, userID string) ([]*blog.Blog, error) {
	// ユーザーIDの検証
	userIDObj, err := user.NewID(userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザーID検証エラー: %w", err)
	}

	blogs, err := b.blogRepo.FindByUserID(ctx, *userIDObj)
	if err != nil {
		return nil, fmt.Errorf("ブログ一覧取得エラー: %w", err)
	}

	return blogs, nil
}

// GetAllBlogs : 全ブログ取得（ページネーション付き）
func (b *blogUsecase) GetAllBlogs(ctx context.Context, page, perPage int) ([]*blog.Blog, error) {
	if page < 0 {
		page = 0
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := page * perPage
	blogs, err := b.blogRepo.FindAll(ctx, offset, perPage)
	if err != nil {
		return nil, fmt.Errorf("ブログ一覧取得エラー: %w", err)
	}

	return blogs, nil
}

// UpdateBlog : ブログの更新
func (b *blogUsecase) UpdateBlog(ctx context.Context, id, userID, title, content string) (*blog.Blog, error) {
	// ブログの検索
	existingBlog, err := b.blogRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ブログ取得エラー: %w", err)
	}

	// 所有権の検証
	if existingBlog.UserID().String() != userID {
		return nil, errors.New("このブログを更新する権限がありません")
	}

	// ブログの更新
	if title != "" {
		if err := existingBlog.UpdateTitle(title); err != nil {
			return nil, fmt.Errorf("タイトル更新エラー: %w", err)
		}
	}

	if content != "" {
		if err := existingBlog.UpdateContent(content); err != nil {
			return nil, fmt.Errorf("コンテンツ更新エラー: %w", err)
		}
	}

	// ブログの保存
	if err := b.blogRepo.Update(ctx, existingBlog); err != nil {
		return nil, fmt.Errorf("ブログ更新エラー: %w", err)
	}

	return existingBlog, nil
}

// DeleteBlog : ブログの削除
func (b *blogUsecase) DeleteBlog(ctx context.Context, id, userID string) error {
	// ブログの検索
	existingBlog, err := b.blogRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ブログ取得エラー: %w", err)
	}

	// 所有権の検証
	if existingBlog.UserID().String() != userID {
		return errors.New("このブログを削除する権限がありません")
	}

	// ブログの削除
	if err := b.blogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ブログ削除エラー: %w", err)
	}

	return nil
}
