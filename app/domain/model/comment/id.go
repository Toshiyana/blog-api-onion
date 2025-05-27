package comment

import (
	"errors"
)

// ID : コメントID
type ID struct {
	value string
}

// NewID : IDの生成
func NewID(value string) (*ID, error) {
	if value == "" {
		return nil, errors.New("IDが空です")
	}
	return &ID{value: value}, nil
}

// String : 文字列表現を返す
func (id ID) String() string {
	return id.value
}
