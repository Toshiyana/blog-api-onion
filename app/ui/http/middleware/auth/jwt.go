package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	userIDContextKey contextKey = "user_id"
)

// JWTMiddleware : JWT認証ミドルウェア
func JWTMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authorization ヘッダーの取得
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "認証が必要です", http.StatusUnauthorized)
				return
			}

			// Bearer トークンの取得
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			// トークンの検証
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("無効な署名方式です")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
				return
			}

			// クレームの取得
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
				return
			}

			// ユーザーIDの取得
			userID, ok := claims["id"].(string)
			if !ok {
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
				return
			}

			// ユーザーIDをコンテキストに設定
			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ExtractUserID : コンテキストからユーザーIDを取得
func ExtractUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}

// RequireAuth : 認証を必須とするミドルウェア
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := ExtractUserID(r.Context())
		if !ok || userID == "" {
			http.Error(w, "認証が必要です", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireOwnership : リソースの所有権を検証するミドルウェア
func RequireOwnership(resourceUserIDFunc func(r *http.Request) (string, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := ExtractUserID(r.Context())
			if !ok || userID == "" {
				http.Error(w, "認証が必要です", http.StatusUnauthorized)
				return
			}

			resourceUserID, err := resourceUserIDFunc(r)
			if err != nil {
				http.Error(w, "リソースの取得に失敗しました", http.StatusInternalServerError)
				return
			}

			if userID != resourceUserID {
				http.Error(w, "このリソースにアクセスする権限がありません", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
