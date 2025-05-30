package dao

import (
	"context"
	"time"

	"myblog/app/domain/model/blog"
	"myblog/app/domain/model/ranking"
	"myblog/app/domain/repository"
	"myblog/app/infra/db/rdb"
)

type rankingRepository struct {
	db *rdb.DB
}

// NewRankingRepository はRankingRepositoryの実装を返す
func NewRankingRepository(db *rdb.DB) repository.RankingRepository {
	return &rankingRepository{
		db: db,
	}
}

// SaveRankings はランキングを保存する
func (r *rankingRepository) SaveRankings(ctx context.Context, rankings []*ranking.Ranking) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// 既存のランキングを削除
	_, err = tx.Exec("DELETE FROM rankings")
	if err != nil {
		return err
	}

	// 新しいランキングを挿入
	for _, rank := range rankings {
		_, err = tx.Exec(
			"INSERT INTO rankings (blog_id, ranking_position, score, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			rank.BlogID.String(), rank.RankingPosition, rank.Score, time.Now(), time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRankings はランキングを取得する
func (r *rankingRepository) GetRankings(ctx context.Context, limit int) ([]*ranking.Ranking, error) {
	rows, err := r.db.Read(ctx).Query(
		"SELECT blog_id, ranking_position, score, created_at, updated_at FROM rankings ORDER BY ranking_position ASC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*ranking.Ranking
	for rows.Next() {
		var rank ranking.Ranking
		var blogIDStr string

		err := rows.Scan(
			&blogIDStr,
			&rank.RankingPosition,
			&rank.Score,
			&rank.CreatedAt,
			&rank.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		blogID, err := blog.NewID(blogIDStr)
		if err != nil {
			return nil, err
		}
		rank.BlogID = *blogID

		rankings = append(rankings, &rank)
	}

	return rankings, nil
}
