// internal/service/comments.go
package logic

import (
	"errors"
	"fmt"

	"web-task/blog/internal/model"

	"gorm.io/gorm"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrInvalidInput    = errors.New("invalid input")
)

// CommentService 评论服务
type CommentService struct {
	DB *gorm.DB
}

// NewCommentService 构造函数
func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{DB: db}
}

// Create 创建评论（需要已认证用户）
// postID 为被评论的文章ID；parentID 为父评论ID（可选，为0表示顶级评论）
func (s *CommentService) Create(userID uint, postID uint, content string) (*model.Comment, error) {
	if content == "" {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidInput)
	}

	comment := model.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}

	if err := s.DB.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// 直接返回创建的评论（不含关联数据）
	return &comment, nil
}

// GetByID 根据ID获取评论（含作者、文章、父评论作者）
// GetByID 已经通过 Preload 加载了关联，无需再处理
func (s *CommentService) GetByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := s.DB.Preload("User").
		Preload("Post").
		Preload("Parent.User").
		First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

// Update 更新评论内容（仅作者可编辑）
func (s *CommentService) Update(userID uint, id uint, content string) (*model.Comment, error) {
	if content == "" {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidInput)
	}

	var comment model.Comment
	if err := s.DB.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	if comment.UserID != userID {
		return nil, errors.New("permission denied")
	}

	if err := s.DB.Model(&comment).Update("content", content).Error; err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	// 直接返回更新后的评论（不含关联数据）
	return &comment, nil
}

// Delete 删除评论（软删除；仅作者可删除）
// includeChildren 是否同时软删除其所有子评论（默认 false）
func (s *CommentService) Delete(userID uint, id uint, includeChildren bool) error {
	var comment model.Comment
	if err := s.DB.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// 权限校验：仅作者可删除
	if comment.UserID != userID {
		return errors.New("permission denied")
	}

	db := s.DB

	// 批量删除子评论（软删除）
	if includeChildren {
		subQuery := s.DB.Model(&model.Comment{}).
			Select("id").
			Where("parent_id = ?", comment.ID)

		if err := db.Where("id = ? OR parent_id IN (?)", comment.ID, subQuery).
			UpdateColumn("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
			return fmt.Errorf("failed to delete comment and children: %w", err)
		}
	} else {
		// 仅删除当前评论
		if err := db.Delete(&comment).Error; err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}
	}

	return nil
}

// Restore 恢复已软删的评论（仅作者可恢复）
// includeChildren 是否同时恢复其所有子评论（默认 false）
func (s *CommentService) Restore(userID uint, id uint, includeChildren bool) error {
	var comment model.Comment
	if err := s.DB.Unscoped().First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// 仅恢复未真正删除的记录
	if comment.DeletedAt == nil {
		return nil // 已恢复或不存在
	}

	// 权限校验：仅作者可恢复
	if comment.UserID != userID {
		return errors.New("permission denied")
	}

	db := s.DB

	// 批量恢复子评论
	if includeChildren {
		subQuery := s.DB.Model(&model.Comment{}).
			Select("id").
			Where("parent_id = ?", comment.ID).
			Unscoped() // 子评论可能已被软删，需 Unscoped 查询

		if err := db.Unscoped().Where("id = ? OR parent_id IN (?)", comment.ID, subQuery).
			UpdateColumn("deleted_at", nil).Error; err != nil {
			return fmt.Errorf("failed to restore comment and children: %w", err)
		}
	} else {
		// 仅恢复当前评论
		if err := db.Model(&comment).UpdateColumn("deleted_at", nil).Error; err != nil {
			return fmt.Errorf("failed to restore comment: %w", err)
		}
	}

	return nil
}
