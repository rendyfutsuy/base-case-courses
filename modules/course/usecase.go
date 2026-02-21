package course

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/course/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateCourse, authId string) (*models.Course, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateCourse, authId string) (*models.Course, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Course, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCourseIndexFilter) ([]models.Course, int, error)
	GetAll(ctx context.Context, filter dto.ReqCourseIndexFilter) ([]models.Course, error)
}
