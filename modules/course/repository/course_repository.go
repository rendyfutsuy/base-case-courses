package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/course"
	"github.com/rendyfutsuy/base-go/modules/course/dto"
	csearch "github.com/rendyfutsuy/base-go/modules/course/repository/searches"
	"gorm.io/gorm"
)

type courseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) course.Repository {
	return &courseRepository{DB: db}
}

func (r *courseRepository) Create(ctx context.Context, createdBy uuid.UUID, title, description, shortDescription string, price, discountRate float64, thumbnailURL *string) (*models.Course, error) {
	now := time.Now().UTC()
	c := &models.Course{
		CreatedBy:        createdBy,
		Title:            title,
		Description:      description,
		ShortDescription: shortDescription,
		Price:            price,
		DiscountRate:     discountRate,
		ThumbnailURL:     thumbnailURL,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := r.DB.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *courseRepository) Update(ctx context.Context, id uuid.UUID, title, description, shortDescription string, price, discountRate float64, thumbnailURL *string) (*models.Course, error) {
	updates := map[string]interface{}{
		"title":             title,
		"description":       description,
		"short_description": shortDescription,
		"price":             price,
		"discount_rate":     discountRate,
		"updated_at":        time.Now().UTC(),
	}
	if thumbnailURL != nil {
		updates["thumbnail_url"] = *thumbnailURL
	} else {
		updates["thumbnail_url"] = nil
	}
	c := &models.Course{}
	err := r.DB.WithContext(ctx).Model(&models.Course{}).
		Where("id = ?", id).
		Updates(updates).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *courseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Course{}).Error
}

func (r *courseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	c := &models.Course{}
	if err := r.DB.WithContext(ctx).
		Table("courses c").
		Select("c.id, c.created_by, c.title, c.description, c.short_description, c.price, c.discount_rate, c.thumbnail_url, c.created_at, c.updated_at").
		Where("c.id = ?", id).
		First(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *courseRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCourseIndexFilter) ([]models.Course, int, error) {
	var courses []models.Course
	query := r.DB.WithContext(ctx).
		Table("courses c").
		Select("c.id, c.title, c.short_description, c.price, c.discount_rate, c.thumbnail_url, c.created_at").
		Where("1=1")

	// Search support
	query = request.ApplySearchConditionFromInterface(query, req.Search, csearch.NewCourseSearchHelper())

	// Filters by parameter relations
	if len(filter.LevelIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'course_level'
				  AND p.id IN (?)
			)
		`, filter.LevelIDs)
	}
	if len(filter.LangIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'lang'
				  AND p.id IN (?)
			)
		`, filter.LangIDs)
	}
	if len(filter.TopicIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'topic'
				  AND p.id IN (?)
			)
		`, filter.TopicIDs)
	}

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:    "c.created_at",
		DefaultSortOrder: "DESC",
		MaxPerPage:       100,
		SortMapping: func(s string) string {
			switch s {
			case "id":
				return "c.id"
			case "title":
				return "c.title"
			case "short_description":
				return "c.short_description"
			case "price":
				return "c.price"
			case "discount_rate":
				return "c.discount_rate"
			case "created_at":
				return "c.created_at"
			default:
				return ""
			}
		},
		NaturalSortColumns: []string{"c.title"},
	}, &courses)
	if err != nil {
		return nil, 0, err
	}
	return courses, total, nil
}

func (r *courseRepository) GetAll(ctx context.Context, filter dto.ReqCourseIndexFilter) ([]models.Course, error) {
	var courses []models.Course
	query := r.DB.WithContext(ctx).
		Table("courses c").
		Select("c.id, c.title, c.short_description, c.price, c.discount_rate, c.thumbnail_url, c.created_at").
		Where("1=1")

	// Search support
	query = request.ApplySearchConditionFromInterface(query, filter.Search, csearch.NewCourseSearchHelper())

	// Filters (same as index)
	if len(filter.LevelIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'course_level'
				  AND p.id IN (?)
			)
		`, filter.LevelIDs)
	}
	if len(filter.LangIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'lang'
				  AND p.id IN (?)
			)
		`, filter.LangIDs)
	}
	if len(filter.TopicIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = 'course'
				  AND ptm.module_id = c.id
				  AND p.type = 'topic'
				  AND p.id IN (?)
			)
		`, filter.TopicIDs)
	}

	if err := query.Order("c.created_at DESC").Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}
