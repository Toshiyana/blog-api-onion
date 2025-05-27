package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myblog/app/infra/dao"
	"myblog/app/infra/db/rdb"
	"myblog/app/ui/http/handler"
	"myblog/app/ui/http/middleware/auth"
	"myblog/app/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// データベース接続
	db, err := rdb.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリ
	userRepo := dao.NewUserRepository(db)
	blogRepo := dao.NewBlogRepository(db)
	commentRepo := dao.NewCommentRepository(db)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_for_development"
		log.Println("Warning: Using default JWT secret")
	}

	// ユースケース
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	blogUsecase := usecase.NewBlogUsecase(blogRepo, userRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, blogRepo, userRepo)

	// ハンドラー
	userHandler := handler.NewUserHandler(userUsecase)
	blogHandler := handler.NewBlogHandler(blogUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)

	// ルーター
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ヘルスチェック
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API ルート
	r.Route("/api", func(r chi.Router) {
		// 認証不要のエンドポイント
		r.Post("/users/register", userHandler.Register)
		r.Post("/users/login", userHandler.Login)

		// 認証が必要なエンドポイント
		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware(jwtSecret))
			r.Use(auth.RequireAuth)

			// ユーザー関連
			r.Get("/users/{id}", userHandler.GetUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Delete("/users/{id}", userHandler.DeleteUser)

			// ブログ関連
			r.Post("/blogs", blogHandler.CreateBlog)
			r.Get("/blogs", blogHandler.GetAllBlogs)
			r.Get("/blogs/{id}", blogHandler.GetBlog)
			r.Get("/users/{id}/blogs", blogHandler.GetUserBlogs)
			r.Put("/blogs/{id}", blogHandler.UpdateBlog)
			r.Delete("/blogs/{id}", blogHandler.DeleteBlog)

			// コメント関連
			r.Post("/blogs/{id}/comments", commentHandler.CreateComment)
			r.Get("/blogs/{id}/comments", commentHandler.GetBlogComments)
			r.Put("/comments/{id}", commentHandler.UpdateComment)
			r.Delete("/comments/{id}", commentHandler.DeleteComment)
		})
	})

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// グレースフルシャットダウン
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Printf("Server is running on port %s", port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
