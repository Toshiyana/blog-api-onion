package handler

import (
	"encoding/json"
	"net/http"

	"myblog/app/ui/http/middleware/auth"
	"myblog/app/usecase"

	"github.com/go-chi/chi/v5"
)

// UserHandler : ユーザーハンドラー
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler : UserHandlerの生成
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// RegisterRequest : ユーザー登録リクエスト
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest : ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest : ユーザー更新リクエスト
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse : ユーザーレスポンス
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Register : ユーザー登録
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := UserResponse{
		ID:        user.ID().String(),
		Username:  user.Username(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Login : ログイン
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.userUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetUser : ユーザー情報取得
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	authUserID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 自分以外のユーザー情報は取得できない
	if id != authUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := UserResponse{
		ID:        user.ID().String(),
		Username:  user.Username(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateUser : ユーザー情報更新
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	authUserID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 自分以外のユーザー情報は更新できない
	if id != authUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), id, req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := UserResponse{
		ID:        user.ID().String(),
		Username:  user.Username(),
		Email:     user.Email(),
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteUser : ユーザー削除
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	authUserID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 自分以外のユーザーは削除できない
	if id != authUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := h.userUsecase.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
