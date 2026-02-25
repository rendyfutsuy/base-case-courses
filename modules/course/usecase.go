package course

import (
	"context"

	"github.com/rendyfutsuy/base-case-courses/helpers/request"
	"github.com/rendyfutsuy/base-case-courses/models"
	"github.com/rendyfutsuy/base-case-courses/modules/course/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateCourse, authId string, thumbnailData []byte, thumbnailName string) (*models.Course, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateCourse, authId string, thumbnailData []byte, thumbnailName string) (*models.Course, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Course, error)
	GetParameterReferences(ctx context.Context, id string) (*dto.ReferenceObject, *dto.ReferenceObject, []dto.ReferenceObject, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCourseIndexFilter) ([]models.Course, int, error)
	GetAll(ctx context.Context, filter dto.ReqCourseIndexFilter) ([]models.Course, error)
}
