package handler

import (
	"encoding/json"
	"net/http"

	"myblog/app/ui/http/middleware/auth"
	"myblog/app/usecase"

	"github.com/go-chi/chi/v5"
)

// CommentHandler : コメントハンドラー
type CommentHandler struct {
	commentUsecase usecase.CommentUsecase
}

// NewCommentHandler : CommentHandlerの生成
func NewCommentHandler(commentUsecase usecase.CommentUsecase) *CommentHandler {
	return &CommentHandler{
		commentUsecase: commentUsecase,
	}
}

// CreateCommentRequest : コメント作成リクエスト
type CreateCommentRequest struct {
	Content string `json:"content"`
}

// UpdateCommentRequest : コメント更新リクエスト
type UpdateCommentRequest struct {
	Content string `json:"content"`
}

// CommentResponse : コメントレスポンス
type CommentResponse struct {
	ID        string `json:"id"`
	BlogID    string `json:"blog_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateComment : コメント作成
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentUsecase.CreateComment(r.Context(), blogID, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CommentResponse{
		ID:        comment.ID().String(),
		BlogID:    comment.BlogID().String(),
		UserID:    comment.UserID().String(),
		Content:   comment.Content(),
		CreatedAt: comment.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetBlogComments : ブログのコメント一覧取得
func (h *CommentHandler) GetBlogComments(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	comments, err := h.commentUsecase.GetCommentsByBlogID(r.Context(), blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []CommentResponse
	for _, comment := range comments {
		resp = append(resp, CommentResponse{
			ID:        comment.ID().String(),
			BlogID:    comment.BlogID().String(),
			UserID:    comment.UserID().String(),
			Content:   comment.Content(),
			CreatedAt: comment.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: comment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateComment : コメント更新
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := h.commentUsecase.UpdateComment(r.Context(), id, userID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CommentResponse{
		ID:        comment.ID().String(),
		BlogID:    comment.BlogID().String(),
		UserID:    comment.UserID().String(),
		Content:   comment.Content(),
		CreatedAt: comment.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteComment : コメント削除
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.commentUsecase.DeleteComment(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
