package repository

import (
	"context"

	"myblog/app/domain/model/ranking"
)

// RankingRepository はランキングのリポジトリインターフェース
type RankingRepository interface {
	// SaveRankings はランキングを保存する
	SaveRankings(ctx context.Context, rankings []*ranking.Ranking) error

	// GetRankings はランキングを取得する
	GetRankings(ctx context.Context, limit int) ([]*ranking.Ranking, error)
}
