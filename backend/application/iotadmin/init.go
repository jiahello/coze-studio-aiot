package iotadmin

import (
	"gorm.io/gorm"

	repo "github.com/coze-dev/coze-studio/backend/domain/iot/repository"
)

type ApplicationService struct {
	DeviceRepo   repo.DeviceRepository
	VoiceRepo    repo.VoiceRepository
	SettingsRepo repo.TTSSettingsRepository
}

var SVC *ApplicationService

// InitService 通过 gorm.DB 装配仓储实现（对齐 domain/user 风格）
func InitService(db *gorm.DB) *ApplicationService {
	SVC = &ApplicationService{
		DeviceRepo:   repo.NewDeviceRepo(db),
		VoiceRepo:    repo.NewVoiceRepo(db),
		SettingsRepo: repo.NewTTSSettingsRepo(db),
	}
	return SVC
}