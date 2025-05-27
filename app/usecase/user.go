package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"myblog/app/domain/model/user"
	"myblog/app/domain/repository"

	"github.com/golang-jwt/jwt/v4"
)

// UserUsecase : ユーザーユースケースインターフェース
type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*user.User, error)
	Login(ctx context.Context, email, password string) (string, error) // JWTトークンを返す
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	UpdateUser(ctx context.Context, id, username, email, password string) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// userUsecase : ユーザーユースケースの実装
type userUsecase struct {
	userRepo  repository.User
	jwtSecret string
}

// NewUserUsecase : ユーザーユースケースの生成
func NewUserUsecase(userRepo repository.User, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register : ユーザー登録
func (u *userUsecase) Register(ctx context.Context, username, email, password string) (*user.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("このメールアドレスは既に使用されています")
	}

	// ユーザーの生成
	newUser, err := user.NewUser(username, email, password)
	if err != nil {
		return nil, fmt.Errorf("ユーザー作成エラー: %w", err)
	}

	// ユーザーの保存
	if err := u.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("ユーザー保存エラー: %w", err)
	}

	return newUser, nil
}

// Login : ログイン
func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	// ユーザーの検索
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// パスワードの検証
	if !user.VerifyPassword(password) {
		return "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID().String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("トークン生成エラー: %w", err)
	}

	return tokenString, nil
}

// GetUserByID : IDによるユーザー取得
func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %w", err)
	}
	return user, nil
}

// UpdateUser : ユーザー情報の更新
func (u *userUsecase) UpdateUser(ctx context.Context, id, username, email, password string) (*user.User, error) {
	// ユーザーの検索
	existingUser, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ユーザー取得エラー: %w", err)
	}

	// ユーザー名の更新
	if username != "" && username != existingUser.Username() {
		if err := existingUser.UpdateUsername(username); err != nil {
			return nil, fmt.Errorf("ユーザー名更新エラー: %w", err)
		}
	}

	// メールアドレスの更新
	if email != "" && email != existingUser.Email() {
		// メールアドレスの重複チェック
		duplicateUser, err := u.userRepo.FindByEmail(ctx, email)
		if err == nil && duplicateUser != nil && duplicateUser.ID().String() != id {
			return nil, errors.New("このメールアドレスは既に使用されています")
		}

		if err := existingUser.UpdateEmail(email); err != nil {
			return nil, fmt.Errorf("メールアドレス更新エラー: %w", err)
		}
	}

	// パスワードの更新
	if password != "" {
		if err := existingUser.UpdatePassword(password); err != nil {
			return nil, fmt.Errorf("パスワード更新エラー: %w", err)
		}
	}

	// ユーザーの保存
	if err := u.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("ユーザー更新エラー: %w", err)
	}

	return existingUser, nil
}

// DeleteUser : ユーザーの削除
func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	// ユーザーの存在確認
	_, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ユーザー取得エラー: %w", err)
	}

	// ユーザーの削除
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("ユーザー削除エラー: %w", err)
	}

	return nil
}
