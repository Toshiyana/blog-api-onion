package usecase

import (
	"context"
	"errors"
	"fmt"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/comment"
	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"
)

// CommentUsecase : コメントユースケースインターフェース
type CommentUsecase interface {
	CreateComment(ctx context.Context, blogID, userID, content string) (*comment.Comment, error)
	GetCommentByID(ctx context.Context, id string) (*comment.Comment, error)
	GetCommentsByBlogID(ctx context.Context, blogID string) ([]*comment.Comment, error)
	UpdateComment(ctx context.Context, id, userID, content string) (*comment.Comment, error)
	DeleteComment(ctx context.Context, id, userID string) error
}

// commentUsecase : コメントユースケースの実装
type commentUsecase struct {
	commentRepo repository.Comment
	blogRepo    repository.Blog
	userRepo    repository.User
}

// NewCommentUsecase : コメントユースケースの生成
func NewCommentUsecase(commentRepo repository.Comment, blogRepo repository.Blog, userRepo repository.User) CommentUsecase {
	return &commentUsecase{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
		userRepo:    userRepo,
	}
}

// CreateComment : コメントの作成
func (c *commentUsecase) CreateComment(ctx context.Context, blogID, userID, content string) (*comment.Comment, error) {
	// ブログの検証
	blogIDObj, err := blog.NewID(blogID)
	if err != nil {
		return nil, fmt.Errorf("ブログID検証エラー: %w", err)
	}

	// ブログの存在確認
	_, err = c.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, fmt.Errorf("ブログ取得エラー: %w", err)
	}

	// ユーザーの検証
	_, err = user.NewID(userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザーID検証エラー: %w", err)
	}

	existingUser, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %w", err)
	}

	// コメントの生成
	newComment, err := comment.NewComment(*blogIDObj, existingUser.ID(), content)
	if err != nil {
		return nil, fmt.Errorf("コメント作成エラー: %w", err)
	}

	// コメントの保存
	if err := c.commentRepo.Save(ctx, newComment); err != nil {
		return nil, fmt.Errorf("コメント保存エラー: %w", err)
	}

	return newComment, nil
}

// GetCommentByID : IDによるコメント取得
func (c *commentUsecase) GetCommentByID(ctx context.Context, id string) (*comment.Comment, error) {
	comment, err := c.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("コメント取得エラー: %w", err)
	}
	return comment, nil
}

// GetCommentsByBlogID : ブログIDによるコメント一覧取得
func (c *commentUsecase) GetCommentsByBlogID(ctx context.Context, blogID string) ([]*comment.Comment, error) {
	// ブログIDの検証
	blogIDObj, err := blog.NewID(blogID)
	if err != nil {
		return nil, fmt.Errorf("ブログID検証エラー: %w", err)
	}

	// ブログの存在確認
	_, err = c.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, fmt.Errorf("ブログ取得エラー: %w", err)
	}

	comments, err := c.commentRepo.FindByBlogID(ctx, *blogIDObj)
	if err != nil {
		return nil, fmt.Errorf("コメント一覧取得エラー: %w", err)
	}

	return comments, nil
}

// UpdateComment : コメントの更新
func (c *commentUsecase) UpdateComment(ctx context.Context, id, userID, content string) (*comment.Comment, error) {
	// コメントの検索
	existingComment, err := c.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("コメント取得エラー: %w", err)
	}

	// 所有権の検証
	if existingComment.UserID().String() != userID {
		return nil, errors.New("このコメントを更新する権限がありません")
	}

	// コメントの更新
	if err := existingComment.UpdateContent(content); err != nil {
		return nil, fmt.Errorf("コンテンツ更新エラー: %w", err)
	}

	// コメントの保存
	if err := c.commentRepo.Update(ctx, existingComment); err != nil {
		return nil, fmt.Errorf("コメント更新エラー: %w", err)
	}

	return existingComment, nil
}

// DeleteComment : コメントの削除
func (c *commentUsecase) DeleteComment(ctx context.Context, id, userID string) error {
	// コメントの検索
	existingComment, err := c.commentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("コメント取得エラー: %w", err)
	}

	// 所有権の検証
	if existingComment.UserID().String() != userID {
		return errors.New("このコメントを削除する権限がありません")
	}

	// コメントの削除
	if err := c.commentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("コメント削除エラー: %w", err)
	}

	return nil
}
