package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateCourse struct {
	Title            string      `json:"title" validate:"required,max=255"`
	Description      string      `json:"description" validate:"required"`
	ShortDescription string      `json:"short_description" validate:"required,max=255"`
	Price            float64     `json:"price" validate:"required"`
	DiscountRate     float64     `json:"discount_rate" validate:"required"`
	LevelID          uuid.UUID   `json:"level_id"`
	LangID           uuid.UUID   `json:"lang_id"`
	TopicIDs         []uuid.UUID `json:"topic_ids"`
	ThumbnailURL     *string
}

type ReqUpdateCourse struct {
	Title            string      `json:"title" validate:"required,max=255"`
	Description      string      `json:"description" validate:"required"`
	ShortDescription string      `json:"short_description" validate:"required,max=255"`
	Price            float64     `json:"price" validate:"required"`
	DiscountRate     float64     `json:"discount_rate" validate:"required"`
	RemoveThumbnail  bool        `json:"remove_thumbnail" form:"remove_thumbnail"`
	LevelID          uuid.UUID   `json:"level_id"`
	LangID           uuid.UUID   `json:"lang_id"`
	TopicIDs         []uuid.UUID `json:"topic_ids"`
	ThumbnailURL     *string
}

type RespCourseIndex struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Price            float64   `json:"price"`
	DiscountRate     float64   `json:"discount_rate"`
	ThumbnailURL     *string   `json:"thumbnail_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ToDBCourse struct {
	Title            string
	Description      string
	ShortDescription string
	Price            float64
	DiscountRate     float64
	RemoveThumbnail  bool
	LevelID          uuid.UUID
	LangID           uuid.UUID
	TopicIDs         []uuid.UUID
	ThumbnailURL     *string
}

func ToRespCourseIndex(m models.Course) RespCourseIndex {
	return RespCourseIndex{
		ID:               m.ID,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Price:            m.Price,
		DiscountRate:     m.DiscountRate,
		ThumbnailURL:     m.ThumbnailURL,
		CreatedAt:        m.CreatedAt,
	}
}

type RespCourse struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	Price            float64   `json:"price"`
	DiscountRate     float64   `json:"discount_rate"`
	ThumbnailURL     *string   `json:"thumbnail_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func ToRespCourse(m models.Course) RespCourse {
	return RespCourse{
		ID:               m.ID,
		Title:            m.Title,
		Description:      m.Description,
		ShortDescription: m.ShortDescription,
		Price:            m.Price,
		DiscountRate:     m.DiscountRate,
		ThumbnailURL:     m.ThumbnailURL,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

type ReqCourseIndexFilter struct {
	Search    string      `query:"search" json:"search"`
	LevelIDs  []uuid.UUID `query:"level_ids" json:"level_ids"`
	TopicIDs  []uuid.UUID `query:"topic_ids" json:"topic_ids"`
	LangIDs   []uuid.UUID `query:"lang_ids" json:"lang_ids"`
	SortBy    string      `query:"sort_by" json:"sort_by"`
	SortOrder string      `query:"sort_order" json:"sort_order"`
}
