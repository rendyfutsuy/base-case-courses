package type_module

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuybase-case-courses/helpers/request"
	"github.com/rendyfutsuybase-case-courses/models"
	"github.com/rendyfutsuybase-case-courses/modules/type/dto"
)

type Repository interface {
	Create(ctx context.Context, subgroupID uuid.UUID, name string, createdBy string) (*models.Type, error)
	Update(ctx context.Context, id uuid.UUID, subgroupID uuid.UUID, name string, updatedBy string) (*models.Type, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Type, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqTypeIndexFilter) ([]models.Type, int, error)
	GetAll(ctx context.Context, filter dto.ReqTypeIndexFilter) ([]models.Type, error)
	ExistsByNameInSubgroup(ctx context.Context, subgroupID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)
	ExistsInBackings(ctx context.Context, typeID uuid.UUID) (bool, error)
}
