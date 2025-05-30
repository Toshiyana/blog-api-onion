package batch

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"myblog/app/ui/http"
	"myblog/app/usecase"

	"github.com/spf13/cobra"
)

// errLocked はロックが既に取得されていることを示すエラー
var errLocked = errors.New("locked")

// Ranking はランキングバッチのハンドラー
type Ranking interface {
	CalculatePopularRanking(cmd *cobra.Command, args []string) error
}

type ranking struct {
	rankingUseCase *usecase.RankingUseCase
	mutex          http.Mutex
}

// NewRanking はRankingハンドラーのコンストラクタ
func NewRanking(rankingUseCase *usecase.RankingUseCase, mutex http.Mutex) Ranking {
	return &ranking{
		rankingUseCase: rankingUseCase,
		mutex:          mutex,
	}
}

// NewCalculatePopularRankingCmd は人気記事ランキング集計コマンドを生成する
func NewCalculatePopularRankingCmd(r Ranking) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate-popular-ranking [days]",
		Args:  cobra.ExactArgs(1),
		Short: "人気記事ランキングを集計する",
		Long:  "指定した日数の間のアクセス数とコメント数に基づいて人気記事ランキングを集計します",
		RunE: func(cmd *cobra.Command, args []string) error {
			return r.CalculatePopularRanking(cmd, args)
		},
		Example: "calculate-popular-ranking 7  # 過去7日間のデータでランキングを集計",
	}

	return cmd
}

// CalculatePopularRanking は人気記事ランキングを集計する
func (r *ranking) CalculatePopularRanking(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// 引数から日数を取得
	days, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("日数の指定が不正です: %w", err)
	}

	if days <= 0 {
		return errors.New("日数は1以上を指定してください")
	}

	// 多重実行を防ぐためロック
	lockID := "calculate-popular-ranking"
	unlock, err := r.mutex.Lock(ctx, lockID, 10*time.Minute)
	if err != nil {
		if errors.Is(err, errLocked) {
			return fmt.Errorf("人気記事ランキング集計が多重実行されています: Mutex.Lock(id: %s): %w", lockID, err)
		}
		return fmt.Errorf("ロック取得処理に失敗しました: Mutex.Lock(id: %s): %w", lockID, err)
	}
	defer unlock()

	// ランキング集計を実行
	if err := r.rankingUseCase.CalculatePopularRanking(ctx, days); err != nil {
		return fmt.Errorf("人気記事ランキング集計に失敗しました: %w", err)
	}

	return nil
}
