package usecase

import (
	"context"
	"fmt"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/ranking"
	"myblog/app/domain/repository"
	"myblog/app/infra/db/rdb"
	"myblog/app/infra/query"
)

// RankingUseCase はランキングに関するユースケース
type RankingUseCase struct {
	rankingRepository repository.RankingRepository
	blogStatsQuery    *query.BlogStats
	txManager         rdb.TransactionManager
}

// NewRankingUseCase はRankingUseCaseのコンストラクタ
func NewRankingUseCase(
	rankingRepository repository.RankingRepository,
	blogStatsQuery *query.BlogStats,
	txManager rdb.TransactionManager,
) *RankingUseCase {
	return &RankingUseCase{
		rankingRepository: rankingRepository,
		blogStatsQuery:    blogStatsQuery,
		txManager:         txManager,
	}
}

// CalculatePopularRanking は人気記事ランキングを集計する
func (u *RankingUseCase) CalculatePopularRanking(ctx context.Context, days int) error {
	// ブログの統計データを取得
	blogStats, err := u.blogStatsQuery.GetBlogRankingData(ctx, days)
	if err != nil {
		return fmt.Errorf("ブログ統計データの取得に失敗しました: %w", err)
	}

	// ランキングを作成
	var rankings []*ranking.Ranking
	for i, stat := range blogStats {
		blogID, err := blog.NewID(stat.BlogID)
		if err != nil {
			return fmt.Errorf("ブログIDのパースに失敗しました: %w", err)
		}

		rank := ranking.NewRanking(*blogID, i+1, stat.TotalScore)
		rankings = append(rankings, rank)
	}

	// トランザクション内でランキングを保存
	err = u.txManager.Transaction(ctx, func(ctx context.Context) error {
		// ランキングを保存
		if err := u.rankingRepository.SaveRankings(ctx, rankings); err != nil {
			return fmt.Errorf("ランキングの保存に失敗しました: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
