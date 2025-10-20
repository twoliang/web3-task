package controller

import (
	"net/http"
	"strconv"

	"web-task/blog/internal/logic"
	"web-task/blog/internal/model"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *logic.PostService
}

func NewPostHandler(postService *logic.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost 创建文章
// @Summary 创建新文章
// @Description 创建新的博客文章，需要用户ID、标题和内容
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param post body CreatePostRequest true "文章信息"
// @Success 201 {object} model.Post
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	// 实际项目中应从JWT令牌中解析用户ID，这里简化处理
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	post, err := h.postService.Create(uint(userID), req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPostByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详情，包含作者信息
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.Post
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的文章ID"})
		return
	}

	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		if err == logic.ErrPostNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "文章不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// ListPosts 文章列表
// @Summary 获取文章列表
// @Description 分页获取文章列表，包含作者信息，按创建时间倒序
// @Tags posts
// @Produce json
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页条数，默认10"
// @Success 200 {object} PostListResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	posts, total, err := h.postService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PostListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Posts:    posts,
	})
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Description 更新指定ID的文章，只能更新自己的文章
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章ID"
// @Param post body UpdatePostRequest true "更新的文章信息"
// @Success 200 {object} model.Post
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的文章ID"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	post, err := h.postService.Update(uint(userID), uint(id), req.Title, req.Content)
	if err != nil {
		switch err {
		case logic.ErrPostNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "文章不存在"})
		case logic.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		default:
			if err.Error() == "permission denied" {
				c.JSON(http.StatusForbidden, ErrorResponse{Message: "没有权限更新此文章"})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost 删除文章
// @Summary 删除文章
// @Description 删除指定ID的文章，只能删除自己的文章
// @Tags posts
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章ID"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的文章ID"})
		return
	}

	err = h.postService.Delete(uint(userID), uint(id))
	if err != nil {
		switch err {
		case logic.ErrPostNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "文章不存在"})
		default:
			if err.Error() == "permission denied" {
				c.JSON(http.StatusForbidden, ErrorResponse{Message: "没有权限删除此文章"})
				return
			}
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// 请求和响应结构体定义
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type PostListResponse struct {
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Posts    []model.Post `json:"posts"`
}

// 注册路由
func (h *PostHandler) RegisterRoutes(router *gin.RouterGroup) {
	posts := router.Group("/posts")
	{
		posts.POST("", h.CreatePost)
		posts.GET("", h.ListPosts)
		posts.GET("/:id", h.GetPostByID)
		posts.PUT("/:id", h.UpdatePost)
		posts.DELETE("/:id", h.DeletePost)
	}
}
