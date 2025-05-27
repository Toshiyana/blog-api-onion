package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"
	"myblog/app/infra/db/rdb"
)

// userDTO : ユーザーのデータ転送オブジェクト
type userDTO struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  []byte    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// toModel : DTOからドメインモデルへの変換
func (dto *userDTO) toModel() (*user.User, error) {
	return user.Reconstruct(
		dto.ID,
		dto.Username,
		dto.Email,
		dto.Password,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
}

// UserRepository : ユーザーリポジトリの実装
type UserRepository struct {
	db *rdb.DB
}

// NewUserRepository : UserRepositoryの生成
func NewUserRepository(db *rdb.DB) repository.User {
	return &UserRepository{db: db}
}

// Save : ユーザーの保存
func (r *UserRepository) Save(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password, created_at, updated_at
		) VALUES (
			:id, :username, :email, :password, :created_at, :updated_at
		)
	`

	params := map[string]interface{}{
		"id":         user.ID().String(),
		"username":   user.Username(),
		"email":      user.Email(),
		"password":   user.Password(),
		"created_at": user.CreatedAt(),
		"updated_at": user.UpdatedAt(),
	}

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		_, err := tx.NamedExec(query, params)
		return err
	}

	_, err := r.db.Write(ctx).NamedExecContext(ctx, query, params)
	return err
}

// FindByID : IDによるユーザー検索
func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT
			id, username, email, password, created_at, updated_at
		FROM
			users
		WHERE
			id = ?
	`

	var dto userDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.QueryRowx(query, id).StructScan(&dto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("user not found with id: %s", id)
			}
			return nil, err
		}
		return dto.toModel()
	}

	err := r.db.Read(ctx).QueryRowxContext(ctx, query, id).StructScan(&dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with id: %s", id)
		}
		return nil, err
	}

	return dto.toModel()
}

// FindByEmail : メールアドレスによるユーザー検索
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT
			id, username, email, password, created_at, updated_at
		FROM
			users
		WHERE
			email = ?
	`

	var dto userDTO

	// トランザクションがあれば使用
	if tx, ok := rdb.GetTx(ctx); ok {
		err := tx.QueryRowx(query, email).StructScan(&dto)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("user not found with email: %s", email)
			}
			return nil, err
		}
		return dto.toModel()
	}

	err := r.db.Read(ctx).QueryRowxContext(ctx, query, email).StructScan(&dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, err
	}

	return dto.toModel()
}

// Update : ユーザーの更新
func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE users
		SET
			username = :username,
			email = :email,
			password = :password,
			updated_at = :updated_at
		WHERE
			id = :id
	`

	params := map[string]interface{}{
		"id":         user.ID().String(),
		"username":   user.Username(),
		"email":      user.Email(),
		"password":   user.Password(),
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
			return fmt.Errorf("user not found with id: %s", user.ID().String())
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
		return fmt.Errorf("user not found with id: %s", user.ID().String())
	}

	return nil
}

// Delete : ユーザーの削除
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
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
			return fmt.Errorf("user not found with id: %s", id)
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
		return fmt.Errorf("user not found with id: %s", id)
	}

	return nil
}
