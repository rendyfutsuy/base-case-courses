package course

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuybase-case-courses/helpers/request"
	"github.com/rendyfutsuybase-case-courses/models"
	"github.com/rendyfutsuybase-case-courses/modules/course/dto"
)

type Repository interface {
	Create(ctx context.Context, createdBy uuid.UUID, data dto.ToDBCourse) (*models.Course, error)
	Update(ctx context.Context, id uuid.UUID, data dto.ToDBCourse) (*models.Course, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCourseIndexFilter) ([]models.Course, int, error)
	GetAll(ctx context.Context, filter dto.ReqCourseIndexFilter) ([]models.Course, error)
}
