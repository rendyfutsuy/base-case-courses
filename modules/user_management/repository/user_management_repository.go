package repository

import (
	"os"

	user "github.com/rendyfutsuybase-case-courses/modules/user_management"
	"github.com/rendyfutsuybase-case-courses/utils"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserManagementRepository(DB *gorm.DB) user.Repository {
	return &userRepository{
		DB: DB,
	}
}

func (repo *userRepository) CreateTable(sqlFilePath string) (err error) {

	sqlCommands, err := os.ReadFile(sqlFilePath)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// Get underlying SQL DB for raw SQL execution
	sqlDB, err := repo.DB.DB()
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	_, err = sqlDB.Exec(string(sqlCommands))
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return err
}
