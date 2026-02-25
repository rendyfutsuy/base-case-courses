package usecase

import (
	"context"

	"github.com/rendyfutsuy/base-case-courses/utils"
	"github.com/rendyfutsuy/base-case-courses/utils/token_storage"
)

func (u *authUsecase) SignOut(ctx context.Context, token string) error {

	// destroy requested jwt token
	err := token_storage.DestroySession(ctx, token)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
