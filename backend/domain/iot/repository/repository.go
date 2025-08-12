package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/iot/internal/dal"
	"github.com/coze-dev/coze-studio/backend/domain/iot"
)

// 构造函数，对齐 user 仓储风格
func NewDeviceRepo(db *gorm.DB) iot.DeviceRepository { return dal.NewDeviceDAO(db) }
func NewVoiceRepo(db *gorm.DB) iot.VoiceRepository { return dal.NewVoiceDAO(db) }
func NewTTSSettingsRepo(db *gorm.DB) iot.TTSSettingsRepository { return dal.NewSettingsDAO(db) }

// 接口定义（与之前在 domain/iot/repository.go 中一致）

type DeviceRepository interface {
	ListDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*iot.HardwareDevice, int64, error)
	UpsertDevice(ctx context.Context, d *iot.HardwareDevice) error
}

type VoiceRepository interface {
	ListVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*iot.TTSVoice, int64, error)
	GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error)
}

type TTSSettingsRepository interface {
	UpsertAppTTS(ctx context.Context, cfg *iot.AppTTSSettings) error
	UpsertHardwareTTS(ctx context.Context, cfg *iot.HardwareTTSSettings) error
	GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*iot.EffectiveTTS, error)
}