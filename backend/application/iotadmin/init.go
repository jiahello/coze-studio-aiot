package iotadmin

import (
	"gorm.io/gorm"
)

type ApplicationService struct {
	DB *gorm.DB
}

var SVC *ApplicationService

func InitService(db *gorm.DB) *ApplicationService {
	SVC = &ApplicationService{DB: db}
	return SVC
}