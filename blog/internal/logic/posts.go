// internal/service/posts.go
package logic

import (
	"errors"
	"fmt"

	"web-task/blog/internal/model"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (s *PostService) Create(UserID uint, title, content string) (*model.Post, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("%w: title and content are required", ErrInvalidInput)
	}

	post := model.Post{
		Title:   title,
		Content: content,
		UserID:  UserID,
	}

	if err := s.DB.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &post, nil
}

func (s *PostService) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := s.DB.Preload("User").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

func (s *PostService) List(page, pageSize int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var posts []model.Post
	var total int64

	if err := s.DB.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	if err := s.DB.Preload("User").
		Limit(pageSize).
		Offset(offset).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, total, nil
}

func (s *PostService) Update(userID uint, id uint, title, content string) (*model.Post, error) {
	var post model.Post
	if err := s.DB.Preload("User").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post.UserID != userID {
		return nil, errors.New("permission denied")
	}

	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if content != "" {
		updates["content"] = content
	}

	if err := s.DB.Model(&post).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return &post, nil
}

func (s *PostService) Delete(userID uint, id uint) error {
	var post model.Post
	if err := s.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to get post: %w", err)
	}

	if post.UserID != userID {
		return errors.New("permission denied")
	}

	if err := s.DB.Delete(&post).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
