package ranking

import (
	"time"

	"myblog/app/domain/model/blog"
)

// Ranking は人気記事ランキングを表すドメインモデル
type Ranking struct {
	BlogID          blog.ID
	RankingPosition int
	Score           int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewRanking はRankingのコンストラクタ
func NewRanking(blogID blog.ID, rankingPosition int, score int) *Ranking {
	return &Ranking{
		BlogID:          blogID,
		RankingPosition: rankingPosition,
		Score:           score,
	}
}
