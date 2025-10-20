// internal/handler/comments.go
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"web-task/blog/internal/logic"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *logic.CommentService
}

func NewCommentHandler(commentService *logic.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment 创建评论
// @Summary 创建新评论
// @Description 为指定文章创建评论，需要用户认证
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param postID path int true "文章ID"
// @Param comment body CreateCommentRequest true "评论信息"
// @Success 201 {object} model.Comment
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /posts/{postID}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// 从请求头获取用户ID（实际项目中应从JWT解析）
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	// 获取文章ID
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的文章ID"})
		return
	}

	// 绑定请求体
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	// 调用服务层
	comment, err := h.commentService.Create(uint(userID), uint(postID), req.Content)
	if err != nil {
		if errors.Is(err, logic.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetCommentByID 获取评论详情
// @Summary 获取评论详情
// @Description 根据ID获取评论详情，包含作者、文章等关联信息
// @Tags comments
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} model.Comment
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id} [get]
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的评论ID"})
		return
	}

	comment, err := h.commentService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, logic.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "评论不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment 更新评论
// @Summary 更新评论内容
// @Description 更新指定评论的内容，仅评论作者可操作
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Param comment body UpdateCommentRequest true "更新的评论内容"
// @Success 200 {object} model.Comment
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的评论ID"})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := h.commentService.Update(uint(userID), uint(id), req.Content)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "评论不存在"})
		case errors.Is(err, logic.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		case err.Error() == "permission denied":
			c.JSON(http.StatusForbidden, ErrorResponse{Message: "没有权限更新此评论"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 软删除指定评论，仅作者可操作，可选择是否删除子评论
// @Tags comments
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Param includeChildren query bool false "是否同时删除子评论，默认false"
// @Success 204 {string} string ""
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的评论ID"})
		return
	}

	includeChildren, _ := strconv.ParseBool(c.DefaultQuery("includeChildren", "false"))

	err = h.commentService.Delete(uint(userID), uint(id), includeChildren)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "评论不存在"})
		case err.Error() == "permission denied":
			c.JSON(http.StatusForbidden, ErrorResponse{Message: "没有权限删除此评论"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// RestoreComment 恢复评论
// @Summary 恢复已删除的评论
// @Description 恢复被软删除的评论，仅作者可操作，可选择是否恢复子评论
// @Tags comments
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Param includeChildren query bool false "是否同时恢复子评论，默认false"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id}/restore [post]
func (h *CommentHandler) RestoreComment(c *gin.Context) {
	userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的用户ID"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "无效的评论ID"})
		return
	}

	includeChildren, _ := strconv.ParseBool(c.DefaultQuery("includeChildren", "false"))

	err = h.commentService.Restore(uint(userID), uint(id), includeChildren)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "评论不存在"})
		case err.Error() == "permission denied":
			c.JSON(http.StatusForbidden, ErrorResponse{Message: "没有权限恢复此评论"})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "评论已恢复"})
}

// 请求和响应结构体定义
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

// 注册路由
func (h *CommentHandler) RegisterRoutes(router *gin.RouterGroup) {
	// 文章下的评论路由
	posts := router.Group("/posts")
	{
		posts.POST("/:postID/comments", h.CreateComment)
	}

	// 评论自身的路由
	comments := router.Group("/comments")
	{
		comments.GET("/:id", h.GetCommentByID)
		comments.PUT("/:id", h.UpdateComment)
		comments.DELETE("/:id", h.DeleteComment)
		comments.POST("/:id/restore", h.RestoreComment)
	}
}
