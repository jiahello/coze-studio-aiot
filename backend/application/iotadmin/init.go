package iotadmin

import (
	"gorm.io/gorm"

	repo "github.com/coze-dev/coze-studio/backend/domain/iot/repository"
	domainSvc "github.com/coze-dev/coze-studio/backend/domain/iot/service"
	crossIOT "github.com/coze-dev/coze-studio/backend/crossdomain/contract/iot"
	implIOT "github.com/coze-dev/coze-studio/backend/crossdomain/impl/iot"
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
	domain := domainSvc.NewService(device, voice, settings)
	SVC = &ApplicationService{
		DeviceRepo:   device,
		VoiceRepo:    voice,
		SettingsRepo: settings,
		Domain:       domain,
	}
	// 注册跨域默认服务
	crossIOT.SetDefaultSVC(implIOT.NewAdapter(domain))
	return SVC
}