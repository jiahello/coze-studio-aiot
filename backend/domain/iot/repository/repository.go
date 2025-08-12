package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/iot/internal/dal"
	"github.com/coze-dev/coze-studio/backend/domain/iot"
)

// 构造函数：返回实现了接口的方法集的 *DAO
func NewDeviceRepo(db *gorm.DB) DeviceRepository { return dal.NewDeviceDAO(db) }
func NewVoiceRepo(db *gorm.DB) VoiceRepository { return dal.NewVoiceDAO(db) }
func NewTTSSettingsRepo(db *gorm.DB) TTSSettingsRepository { return dal.NewSettingsDAO(db) }

// 接口定义

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
	GetAppTTSByAppID(ctx context.Context, appID uint64) (*iot.AppTTSSettings, error)
	GetHardwareTTSByDeviceID(ctx context.Context, deviceID string) (*iot.HardwareTTSSettings, error)
}