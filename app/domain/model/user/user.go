package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User : ユーザーエンティティ
type User struct {
	id        ID
	username  string
	email     string
	password  []byte
	createdAt time.Time
	updatedAt time.Time
}

// NewUser : ユーザーの生成
func NewUser(username, email, password string) (*User, error) {
	if username == "" {
		return nil, errors.New("ユーザー名が空です")
	}
	if email == "" {
		return nil, errors.New("メールアドレスが空です")
	}
	if password == "" {
		return nil, errors.New("パスワードが空です")
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		id:        *id,
		username:  username,
		email:     email,
		password:  hashedPassword,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct : ユーザーの再構築（DBからの読み込み時など）
func Reconstruct(id, username, email string, password []byte, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := NewID(id)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        *userID,
		username:  username,
		email:     email,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// VerifyPassword : パスワードの検証
func (u User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.password, []byte(password))
	return err == nil
}

// ID : IDの取得
func (u User) ID() ID {
	return u.id
}

// Username : ユーザー名の取得
func (u User) Username() string {
	return u.username
}

// Email : メールアドレスの取得
func (u User) Email() string {
	return u.email
}

// Password : パスワードの取得
func (u User) Password() []byte {
	return u.password
}

// CreatedAt : 作成日時の取得
func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt : 更新日時の取得
func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

// UpdateUsername : ユーザー名の更新
func (u *User) UpdateUsername(username string) error {
	if username == "" {
		return errors.New("ユーザー名が空です")
	}
	u.username = username
	u.updatedAt = time.Now()
	return nil
}

// UpdateEmail : メールアドレスの更新
func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return errors.New("メールアドレスが空です")
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// UpdatePassword : パスワードの更新
func (u *User) UpdatePassword(password string) error {
	if password == "" {
		return errors.New("パスワードが空です")
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.password = hashedPassword
	u.updatedAt = time.Now()
	return nil
}
