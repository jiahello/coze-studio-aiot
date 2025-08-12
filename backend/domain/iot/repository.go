package iot

import "context"

// DeviceRepository 负责设备实体的持久化
 type DeviceRepository interface {
	ListDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*HardwareDevice, int64, error)
	UpsertDevice(ctx context.Context, d *HardwareDevice) error
}

// VoiceRepository 负责音色数据的查询
 type VoiceRepository interface {
	ListVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*TTSVoice, int64, error)
	GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error)
}

// TTSSettingsRepository 负责 TTS 设置的读写与生效规则查询
 type TTSSettingsRepository interface {
	UpsertAppTTS(ctx context.Context, cfg *AppTTSSettings) error
	UpsertHardwareTTS(ctx context.Context, cfg *HardwareTTSSettings) error
	GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*EffectiveTTS, error)
}