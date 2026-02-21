package sub_group

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuybase-case-courses/helpers/request"
	"github.com/rendyfutsuybase-case-courses/models"
	"github.com/rendyfutsuybase-case-courses/modules/sub-group/dto"
)

type Repository interface {
	Create(ctx context.Context, goodsGroupID uuid.UUID, name string, createdBy string) (*models.SubGroup, error)
	Update(ctx context.Context, id uuid.UUID, goodsGroupID uuid.UUID, name string, updatedBy string) (*models.SubGroup, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SubGroup, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error)
	GetAll(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error)
	ExistsByName(ctx context.Context, goodsGroupID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)
	ExistsInTypes(ctx context.Context, subGroupID uuid.UUID) (bool, error)
}
