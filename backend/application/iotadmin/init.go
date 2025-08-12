package iotadmin

import (
	"gorm.io/gorm"

	repo "github.com/coze-dev/coze-studio/backend/domain/iot/repository"
	domainSvc "github.com/coze-dev/coze-studio/backend/domain/iot/service"
)

type ApplicationService struct {
	DeviceRepo   repo.DeviceRepository
	VoiceRepo    repo.VoiceRepository
	SettingsRepo repo.TTSSettingsRepository
	Domain       *domainSvc.Service
}

var SVC *ApplicationService

// InitService 通过 gorm.DB 装配仓储实现（对齐 domain/user 风格）
func InitService(db *gorm.DB) *ApplicationService {
	device := repo.NewDeviceRepo(db)
	voice := repo.NewVoiceRepo(db)
	settings := repo.NewTTSSettingsRepo(db)
	SVC = &ApplicationService{
		DeviceRepo:   device,
		VoiceRepo:    voice,
		SettingsRepo: settings,
		Domain:       domainSvc.NewService(device, voice, settings),
	}
	return SVC
}