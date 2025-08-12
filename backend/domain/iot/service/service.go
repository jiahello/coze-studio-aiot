package service

import (
	"context"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
	repo "github.com/coze-dev/coze-studio/backend/domain/iot/repository"
)

type Service struct {
	Devices   repo.DeviceRepository
	Voices    repo.VoiceRepository
	Settings  repo.TTSSettingsRepository
}

func NewService(d repo.DeviceRepository, v repo.VoiceRepository, s repo.TTSSettingsRepository) *Service {
	return &Service{Devices: d, Voices: v, Settings: s}
}

// GetEffectiveTTS 组合设备与应用级设置，返回最终生效项（device > app > default）
func (s *Service) GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*domain.EffectiveTTS, error) {
	if deviceID != "" {
		if h, err := s.Settings.GetHardwareTTSByDeviceID(ctx, deviceID); err == nil && h != nil {
			return &domain.EffectiveTTS{Provider: h.Provider, Model: h.Model, Voice: h.Voice, Source: "device"}, nil
		}
	}
	if appID != nil {
		if a, err := s.Settings.GetAppTTSByAppID(ctx, *appID); err == nil && a != nil {
			return &domain.EffectiveTTS{Provider: a.Provider, Model: a.Model, Voice: a.Voice, Source: "app"}, nil
		}
	}
	return &domain.EffectiveTTS{Provider: "doubao", Model: "speech-1", Voice: "doubao-standard", Source: "default"}, nil
}