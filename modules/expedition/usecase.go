package expedition

import (
	"context"

	"github.com/rendyfutsuybase-case-courses/helpers/request"
	"github.com/rendyfutsuybase-case-courses/models"
	"github.com/rendyfutsuybase-case-courses/modules/expedition/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateExpedition, authId string) (*models.Expedition, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateExpedition, authId string) (*models.Expedition, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Expedition, error)
	GetContactsByExpeditionID(ctx context.Context, id string) ([]models.ExpeditionContact, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error)
	GetAll(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, error)
	Export(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]byte, error)
}
