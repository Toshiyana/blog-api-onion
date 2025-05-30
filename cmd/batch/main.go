package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"myblog/app/infra/dao"
	"myblog/app/infra/db/rdb"
	"myblog/app/infra/query"
	"myblog/app/ui/batch"
	"myblog/app/ui/http"
	"myblog/app/usecase"

	"github.com/spf13/cobra"
)

// RootCmd はルートコマンド
var RootCmd = &cobra.Command{
	Use:   "batch",
	Short: "MyBlogバッチ処理",
}

// run はメイン処理
func run() int {
	startTime := time.Now()
	ctx := context.Background()

	// データベース接続
	db, err := rdb.NewDB()
	if err != nil {
		fmt.Printf("データベース接続に失敗しました: %v\n", err)
		return 1
	}
	defer db.Close()

	// 依存関係の構築
	mutex := http.NewMutex(nil) // TODO: キャッシュクライアントを渡す
	txManager := rdb.NewDefaultTransactionManager(db)

	// ランキング関連の依存関係
	rankingRepository := dao.NewRankingRepository(db)
	blogStatsQuery := query.NewBlogStats(db)
	rankingUseCase := usecase.NewRankingUseCase(rankingRepository, blogStatsQuery, txManager)
	rankingHandler := batch.NewRanking(rankingUseCase, *mutex)
	calculatePopularRankingCmd := batch.NewCalculatePopularRankingCmd(rankingHandler)

	// コマンドの登録
	RootCmd.AddCommand(calculatePopularRankingCmd)

	// コマンドの実行
	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("バッチ実行中にエラーが発生しました: %v\n", err)
		return 1
	}

	executionTime := time.Since(startTime)
	cmdName := strings.Join(os.Args, " ")
	fmt.Printf("バッチ実行が完了しました。バッチ名: %s, 実行時間: %v\n", cmdName, executionTime)
	return 0
}

func main() {
	os.Exit(run())
}
