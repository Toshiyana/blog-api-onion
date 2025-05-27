package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"
	"myblog/app/infra/db/rdb"
)

// blogDTO : ブログのデータ転送オブジェクト
type blogDTO struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toModel : DTOからドメインモデルへの変換
func (dto *blogDTO) toModel() (*blog.Blog, error) {
	userID, err := user.NewID(dto.UserID)
	if err != nil {
		return nil, err
	}

	return blog.Reconstruct(
		dto.ID,
		*userID,
		dto.Title,
		dto.Content,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

// BlogRepository : ブログリポジトリの実装
type BlogRepository struct {
	db *rdb.DB
}

// NewBlogRepository : BlogRepositoryの生成
func NewBlogRepository(db *rdb.DB) repository.Blog {
	return &BlogRepository{db: db}
}

// Save : ブログの保存
func (r *BlogRepository) Save(ctx context.Context, blog *blog.Blog) error {
	query := `
		INSERT INTO blogs (
			id, user_id, title, content, created_at, updated_at
		) VALUES (
			:id, :user_id, :title, :content, :created_at, :updated_at
		)
	`

	params := map[string]interface{}{
		"id":         blog.ID().String(),
		"user_id":    blog.UserID().String(),
		"title":      blog.Title(),
		"content":    blog.Content(),
		"created_at": blog.CreatedAt(),
		"updated_at": blog.UpdatedAt(),
	}

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		_, err := tx.NamedExec(query, params)
		return err
	}

	_, err := r.db.Write(ctx).NamedExecContext(ctx, query, params)
	return err
}

// FindByID : IDによるブログ検索
func (r *BlogRepository) FindByID(ctx context.Context, id string) (*blog.Blog, error) {
	query := `
		SELECT
			id, user_id, title, content, created_at, updated_at
		FROM
			blogs
		WHERE
			id = ?
	`

	var dto blogDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.QueryRowx(query, id).StructScan(&dto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("blog not found with id: %s", id)
			}
			return nil, err
		}
		return dto.toModel()
	}

	err := r.db.Read(ctx).QueryRowxContext(ctx, query, id).StructScan(&dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("blog not found with id: %s", id)
		}
		return nil, err
	}

	return dto.toModel()
}

// FindByUserID : ユーザーIDによるブログ検索
func (r *BlogRepository) FindByUserID(ctx context.Context, userID user.ID) ([]*blog.Blog, error) {
	query := `
		SELECT
			id, user_id, title, content, created_at, updated_at
		FROM
			blogs
		WHERE
			user_id = ?
		ORDER BY
			created_at DESC
	`

	var dtos []blogDTO

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

	blogs := make([]*blog.Blog, len(dtos))
	for i, dto := range dtos {
		blog, err := dto.toModel()
		if err != nil {
			return nil, err
		}
		blogs[i] = blog
	}

	return blogs, nil
}

// FindAll : 全ブログ検索（ページネーション付き）
func (r *BlogRepository) FindAll(ctx context.Context, offset, limit int) ([]*blog.Blog, error) {
	query := `
		SELECT
			id, user_id, title, content, created_at, updated_at
		FROM
			blogs
		ORDER BY
			created_at DESC
		LIMIT ? OFFSET ?
	`

	var dtos []blogDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.Select(&dtos, query, limit, offset)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Read(ctx).SelectContext(ctx, &dtos, query, limit, offset)
		if err != nil {
			return nil, err
		}
	}

	blogs := make([]*blog.Blog, len(dtos))
	for i, dto := range dtos {
		blog, err := dto.toModel()
		if err != nil {
			return nil, err
		}
		blogs[i] = blog
	}

	return blogs, nil
}

// Update : ブログの更新
func (r *BlogRepository) Update(ctx context.Context, blog *blog.Blog) error {
	query := `
		UPDATE blogs
		SET
			title = :title,
			content = :content,
			updated_at = :updated_at
		WHERE
			id = :id
	`

	params := map[string]interface{}{
		"id":         blog.ID().String(),
		"title":      blog.Title(),
		"content":    blog.Content(),
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
			return fmt.Errorf("blog not found with id: %s", blog.ID().String())
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
		return fmt.Errorf("blog not found with id: %s", blog.ID().String())
	}

	return nil
}

// Delete : ブログの削除
func (r *BlogRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM blogs
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
			return fmt.Errorf("blog not found with id: %s", id)
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
		return fmt.Errorf("blog not found with id: %s", id)
	}

	return nil
}
