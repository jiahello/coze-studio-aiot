package iotadmin

import (
	"gorm.io/gorm"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
	mysqlrepo "github.com/coze-dev/coze-studio/backend/infra/impl/mysql/iot"
)

type ApplicationService struct {
	DeviceRepo  domain.DeviceRepository
	VoiceRepo   domain.VoiceRepository
	SettingsRepo domain.TTSSettingsRepository
}

var SVC *ApplicationService

// InitService 通过 gorm.DB 装配 mysql 仓储实现
func InitService(db *gorm.DB) *ApplicationService {
	bundle := mysqlrepo.NewRepositories(db)
	SVC = &ApplicationService{
		DeviceRepo:   bundle.Devices,
		VoiceRepo:    bundle.Voices,
		SettingsRepo: bundle.Settings,
	}
	return SVC
}