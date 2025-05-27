package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"myblog/app/ui/http/middleware/auth"
	"myblog/app/usecase"

	"github.com/go-chi/chi/v5"
)

// BlogHandler : ブログハンドラー
type BlogHandler struct {
	blogUsecase usecase.BlogUsecase
}

// NewBlogHandler : BlogHandlerの生成
func NewBlogHandler(blogUsecase usecase.BlogUsecase) *BlogHandler {
	return &BlogHandler{
		blogUsecase: blogUsecase,
	}
}

// CreateBlogRequest : ブログ作成リクエスト
type CreateBlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateBlogRequest : ブログ更新リクエスト
type UpdateBlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// BlogResponse : ブログレスポンス
type BlogResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateBlog : ブログ作成
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := h.blogUsecase.CreateBlog(r.Context(), userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := BlogResponse{
		ID:        blog.ID().String(),
		UserID:    blog.UserID().String(),
		Title:     blog.Title(),
		Content:   blog.Content(),
		CreatedAt: blog.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: blog.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetBlog : ブログ取得
func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	blog, err := h.blogUsecase.GetBlogByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := BlogResponse{
		ID:        blog.ID().String(),
		UserID:    blog.UserID().String(),
		Title:     blog.Title(),
		Content:   blog.Content(),
		CreatedAt: blog.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: blog.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetAllBlogs : 全ブログ取得
func (h *BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータの取得
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := 0
	perPage := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p >= 0 {
			page = p
		}
	}

	if perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	blogs, err := h.blogUsecase.GetAllBlogs(r.Context(), page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []BlogResponse
	for _, blog := range blogs {
		resp = append(resp, BlogResponse{
			ID:        blog.ID().String(),
			UserID:    blog.UserID().String(),
			Title:     blog.Title(),
			Content:   blog.Content(),
			CreatedAt: blog.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: blog.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetUserBlogs : ユーザーのブログ一覧取得
func (h *BlogHandler) GetUserBlogs(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	blogs, err := h.blogUsecase.GetBlogsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []BlogResponse
	for _, blog := range blogs {
		resp = append(resp, BlogResponse{
			ID:        blog.ID().String(),
			UserID:    blog.UserID().String(),
			Title:     blog.Title(),
			Content:   blog.Content(),
			CreatedAt: blog.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: blog.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateBlog : ブログ更新
func (h *BlogHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := h.blogUsecase.UpdateBlog(r.Context(), id, userID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := BlogResponse{
		ID:        blog.ID().String(),
		UserID:    blog.UserID().String(),
		Title:     blog.Title(),
		Content:   blog.Content(),
		CreatedAt: blog.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: blog.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteBlog : ブログ削除
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Blog ID is required", http.StatusBadRequest)
		return
	}

	// 認証済みユーザーIDの取得
	userID, ok := auth.ExtractUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.blogUsecase.DeleteBlog(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
