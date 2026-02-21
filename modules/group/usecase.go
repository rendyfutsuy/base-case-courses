package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuybase-case-courses/helpers/request"
	"github.com/rendyfutsuybase-case-courses/models"
	"github.com/rendyfutsuybase-case-courses/modules/group/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateGroup, authId string) (*models.Group, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateGroup, authId string) (*models.Group, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Group, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.Group, int, error)
	GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.Group, error)
	Export(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]byte, error)
	ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error)
}
