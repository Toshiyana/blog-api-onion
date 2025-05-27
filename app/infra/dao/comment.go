package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/comment"
	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"
	"myblog/app/infra/db/rdb"
)

// commentDTO : コメントのデータ転送オブジェクト
type commentDTO struct {
	ID        string    `db:"id"`
	BlogID    string    `db:"blog_id"`
	UserID    string    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toModel : DTOからドメインモデルへの変換
func (dto *commentDTO) toModel() (*comment.Comment, error) {
	blogID, err := blog.NewID(dto.BlogID)
	if err != nil {
		return nil, err
	}

	userID, err := user.NewID(dto.UserID)
	if err != nil {
		return nil, err
	}

	return comment.Reconstruct(
		dto.ID,
		*blogID,
		*userID,
		dto.Content,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

// CommentRepository : コメントリポジトリの実装
type CommentRepository struct {
	db *rdb.DB
}

// NewCommentRepository : CommentRepositoryの生成
func NewCommentRepository(db *rdb.DB) repository.Comment {
	return &CommentRepository{db: db}
}

// Save : コメントの保存
func (r *CommentRepository) Save(ctx context.Context, comment *comment.Comment) error {
	query := `
		INSERT INTO comments (
			id, blog_id, user_id, content, created_at, updated_at
		) VALUES (
			:id, :blog_id, :user_id, :content, :created_at, :updated_at
		)
	`

	params := map[string]interface{}{
		"id":         comment.ID().String(),
		"blog_id":    comment.BlogID().String(),
		"user_id":    comment.UserID().String(),
		"content":    comment.Content(),
		"created_at": comment.CreatedAt(),
		"updated_at": comment.UpdatedAt(),
	}

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		_, err := tx.NamedExec(query, params)
		return err
	}

	_, err := r.db.Write(ctx).NamedExecContext(ctx, query, params)
	return err
}

// FindByID : IDによるコメント検索
func (r *CommentRepository) FindByID(ctx context.Context, id string) (*comment.Comment, error) {
	query := `
		SELECT
			id, blog_id, user_id, content, created_at, updated_at
		FROM
			comments
		WHERE
			id = ?
	`

	var dto commentDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.QueryRowx(query, id).StructScan(&dto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("comment not found with id: %s", id)
			}
			return nil, err
		}
		return dto.toModel()
	}

	err := r.db.Read(ctx).QueryRowxContext(ctx, query, id).StructScan(&dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("comment not found with id: %s", id)
		}
		return nil, err
	}

	return dto.toModel()
}

// FindByBlogID : ブログIDによるコメント検索
func (r *CommentRepository) FindByBlogID(ctx context.Context, blogID blog.ID) ([]*comment.Comment, error) {
	query := `
		SELECT
			id, blog_id, user_id, content, created_at, updated_at
		FROM
			comments
		WHERE
			blog_id = ?
		ORDER BY
			created_at ASC
	`

	var dtos []commentDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.Select(&dtos, query, blogID.String())
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Read(ctx).SelectContext(ctx, &dtos, query, blogID.String())
		if err != nil {
			return nil, err
		}
	}

	comments := make([]*comment.Comment, len(dtos))
	for i, dto := range dtos {
		comment, err := dto.toModel()
		if err != nil {
			return nil, err
		}
		comments[i] = comment
	}

	return comments, nil
}

// FindByUserID : ユーザーIDによるコメント検索
func (r *CommentRepository) FindByUserID(ctx context.Context, userID user.ID) ([]*comment.Comment, error) {
	query := `
		SELECT
			id, blog_id, user_id, content, created_at, updated_at
		FROM
			comments
		WHERE
			user_id = ?
		ORDER BY
			created_at DESC
	`

	var dtos []commentDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.Select(&dtos, query, userID.String())
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Read(ctx).SelectContext(ctx, &dtos, query, userID.String())
		if err != nil {
			return nil, err
		}
	}

	comments := make([]*comment.Comment, len(dtos))
	for i, dto := range dtos {
		comment, err := dto.toModel()
		if err != nil {
			return nil, err
		}
		comments[i] = comment
	}

	return comments, nil
}

// Update : コメントの更新
func (r *CommentRepository) Update(ctx context.Context, comment *comment.Comment) error {
	query := `
		UPDATE comments
		SET
			content = :content,
			updated_at = :updated_at
		WHERE
			id = :id
	`

	params := map[string]interface{}{
		"id":         comment.ID().String(),
		"content":    comment.Content(),
		"updated_at": time.Now(),
	}

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		result, err := tx.NamedExec(query, params)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("comment not found with id: %s", comment.ID().String())
		}

		return nil
	}

	result, err := r.db.Write(ctx).NamedExecContext(ctx, query, params)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found with id: %s", comment.ID().String())
	}

	return nil
}

// Delete : コメントの削除
func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM comments
		WHERE id = ?
	`

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		result, err := tx.Exec(query, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("comment not found with id: %s", id)
		}

		return nil
	}

	result, err := r.db.Write(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment not found with id: %s", id)
	}

	return nil
}
