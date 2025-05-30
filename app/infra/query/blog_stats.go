package query

import (
	"context"

	"myblog/app/infra/db/rdb"
)

// BlogStats はブログの統計情報を取得するクエリ
type BlogStats struct {
	db *rdb.DB
}

// NewBlogStats はBlogStatsのコンストラクタ
func NewBlogStats(db *rdb.DB) *BlogStats {
	return &BlogStats{
		db: db,
	}
}

// BlogRankingData はブログのランキングデータ
type BlogRankingData struct {
	BlogID       string
	CommentCount int
	TotalScore   int
}

// GetBlogRankingData は指定期間内のブログランキングデータを取得する
func (b *BlogStats) GetBlogRankingData(ctx context.Context, days int) ([]BlogRankingData, error) {
	query := `
		SELECT
			b.id as blog_id,
			COALESCE(c.comment_count, 0) as comment_count,
			COALESCE(c.comment_count, 0) as total_score
		FROM
			blogs b
		LEFT JOIN (
			SELECT blog_id, COUNT(*) as comment_count
			FROM comments
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
			GROUP BY blog_id
		) c ON b.id = c.blog_id
		ORDER BY total_score DESC
	`

	rows, err := b.db.Read(ctx).Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BlogRankingData
	for rows.Next() {
		var data BlogRankingData
		err := rows.Scan(
			&data.BlogID,
			&data.CommentCount,
			&data.TotalScore,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}
