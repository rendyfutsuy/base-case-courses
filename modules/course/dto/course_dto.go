package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-case-courses/models"
	utilsServices "github.com/rendyfutsuy/base-case-courses/utils/services"
)

type ReqCreateCourse struct {
	Title            string      `json:"title" validate:"required,max=255" form:"title"`
	Description      string      `json:"description" validate:"required" form:"description"`
	ShortDescription string      `json:"short_description" validate:"required,max=255" form:"short_description"`
	Price            float64     `json:"price" validate:"required" form:"price"`
	DiscountRate     float64     `json:"discount_rate" validate:"required" form:"discount_rate"`
	LevelID          uuid.UUID   `json:"level_id" form:"level_id"`
	LangID           uuid.UUID   `json:"lang_id" form:"lang_id"`
	TopicIDs         []uuid.UUID `json:"topic_ids" form:"topic_ids"`
	ThumbnailURL     *string     // mutate from thumbnail file
}

type ReqUpdateCourse struct {
	Title            string      `json:"title" validate:"required,max=255" form:"title"`
	Description      string      `json:"description" validate:"required" form:"description"`
	ShortDescription string      `json:"short_description" validate:"required,max=255" form:"short_description"`
	Price            float64     `json:"price" validate:"required" form:"price"`
	DiscountRate     float64     `json:"discount_rate" validate:"required" form:"discount_rate"`
	RemoveThumbnail  bool        `json:"remove_thumbnail" form:"remove_thumbnail" default:"false"`
	LevelID          uuid.UUID   `json:"level_id" form:"level_id"`
	LangID           uuid.UUID   `json:"lang_id" form:"lang_id"`
	TopicIDs         []uuid.UUID `json:"topic_ids" form:"topic_ids"`
	ThumbnailURL     *string     // mutate from thumbnail file
}

type RespCourseIndex struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Price            float64   `json:"price"`
	DiscountRate     float64   `json:"discount_rate"`
	DiscountedPrice  float64   `json:"discounted_price"`
	ThumbnailURL     *string   `json:"thumbnail_url"`
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
	discounted := m.Price - (m.Price * (m.DiscountRate / 100.0))
	return RespCourseIndex{
		ID:               m.ID,
		Title:            m.Title,
		ShortDescription: m.ShortDescription,
		Price:            m.Price,
		DiscountRate:     m.DiscountRate,
		DiscountedPrice:  discounted,
		ThumbnailURL:     m.ThumbnailURL,
		CreatedAt:        m.CreatedAt,
	}
}

type ReferenceObject struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RespCourse struct {
	ID               uuid.UUID         `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	ShortDescription string            `json:"short_description"`
	Price            float64           `json:"price"`
	DiscountRate     float64           `json:"discount_rate"`
	DiscountedPrice  float64           `json:"discounted_price"`
	Level            *ReferenceObject  `json:"level"`
	Lang             *ReferenceObject  `json:"lang"`
	Topics           []ReferenceObject `json:"topics"`
	ThumbnailURL     *string           `json:"thumbnail_url"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

func ToRespCourse(m models.Course) RespCourse {
	discounted := m.Price - (m.Price * (m.DiscountRate / 100.0))

	var presignedURL string
	if m.ThumbnailURL != nil {
		presignedURL, _ = utilsServices.GeneratePresignedURL(*m.ThumbnailURL)
	}

	return RespCourse{
		ID:               m.ID,
		Title:            m.Title,
		Description:      m.Description,
		ShortDescription: m.ShortDescription,
		Price:            m.Price,
		DiscountRate:     m.DiscountRate,
		DiscountedPrice:  discounted,
		ThumbnailURL:     &presignedURL,
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
