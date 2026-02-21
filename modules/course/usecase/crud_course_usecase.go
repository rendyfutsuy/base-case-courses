package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/course"
	"github.com/rendyfutsuy/base-go/modules/course/dto"
	paramMod "github.com/rendyfutsuy/base-go/modules/parameter"
	"github.com/rendyfutsuy/base-go/utils"
)

type courseUsecase struct {
	repo      course.Repository
	paramRepo paramMod.Repository
}

func NewCourseUsecase(repo course.Repository, paramRepo paramMod.Repository) course.Usecase {
	return &courseUsecase{repo: repo, paramRepo: paramRepo}
}

func (u *courseUsecase) Create(ctx context.Context, req *dto.ReqCreateCourse, authId string) (*models.Course, error) {
	// Validate parameter types
	if err := u.validateParameterType(ctx, req.LevelID, "course_level"); err != nil {
		return nil, err
	}
	if err := u.validateParameterType(ctx, req.LangID, "lang"); err != nil {
		return nil, err
	}
	for _, tid := range req.TopicIDs {
		if err := u.validateParameterType(ctx, tid, "topic"); err != nil {
			return nil, err
		}
	}

	createdBy := uuid.Nil
	if authId != "" {
		if uid, err := utils.StringToUUID(authId); err == nil {
			createdBy = uid
		}
	}

	c, err := u.repo.Create(ctx, createdBy, req.Title, req.Description, req.ShortDescription, req.Price, req.DiscountRate, req.ThumbnailURL)
	if err != nil {
		return nil, err
	}

	// Assign relations via parameters_to_module
	if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, []uuid.UUID{req.LevelID}); err != nil {
		return nil, err
	}
	if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, []uuid.UUID{req.LangID}); err != nil {
		return nil, err
	}
	if len(req.TopicIDs) > 0 {
		if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, req.TopicIDs); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (u *courseUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdateCourse, authId string) (*models.Course, error) {
	// Validate parameter types
	if err := u.validateParameterType(ctx, req.LevelID, "course_level"); err != nil {
		return nil, err
	}
	if err := u.validateParameterType(ctx, req.LangID, "lang"); err != nil {
		return nil, err
	}
	for _, tid := range req.TopicIDs {
		if err := u.validateParameterType(ctx, tid, "topic"); err != nil {
			return nil, err
		}
	}

	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	c, err := u.repo.Update(ctx, cid, req.Title, req.Description, req.ShortDescription, req.Price, req.DiscountRate, req.ThumbnailURL)
	if err != nil {
		return nil, err
	}

	// Re-assign relations: for simplicity, append new assignments (idempotency relies on unique checks if needed)
	if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, []uuid.UUID{req.LevelID}); err != nil {
		return nil, err
	}
	if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, []uuid.UUID{req.LangID}); err != nil {
		return nil, err
	}
	if len(req.TopicIDs) > 0 {
		if err := u.paramRepo.AssignParametersToModule(ctx, "course", c.ID, req.TopicIDs); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (u *courseUsecase) Delete(ctx context.Context, id string, authId string) error {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, cid)
}

func (u *courseUsecase) GetByID(ctx context.Context, id string) (*models.Course, error) {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, cid)
}

func (u *courseUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCourseIndexFilter) ([]models.Course, int, error) {
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *courseUsecase) GetAll(ctx context.Context, filter dto.ReqCourseIndexFilter) ([]models.Course, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *courseUsecase) validateParameterType(ctx context.Context, id uuid.UUID, expectedType string) error {
	p, err := u.paramRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil || p.Type == nil || *p.Type != expectedType {
		return errors.New("invalid parameter type for " + expectedType)
	}
	return nil
}
