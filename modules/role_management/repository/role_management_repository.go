package repository

import (
	"os"

	role "github.com/rendyfutsuybase-case-courses/modules/role_management"
	"github.com/rendyfutsuybase-case-courses/utils"
	"gorm.io/gorm"
)

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleManagementRepository(DB *gorm.DB) role.Repository {
	return &roleRepository{DB}
}

func (repo *roleRepository) CreateTable(sqlFilePath string) (err error) {

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
