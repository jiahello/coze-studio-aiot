package iotadmin

import (
	"context"
	"errors"
	"time"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
)

// ListHardwareDevices returns paginated devices under a space.
func (s *ApplicationService) ListHardwareDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*domain.HardwareDevice, int64, error) {
	return s.DeviceRepo.ListDevices(ctx, spaceID, page, pageSize, keyword)
}

// UpsertHardwareDevice creates or updates a device by (space_id, device_id) or by id.
func (s *ApplicationService) UpsertHardwareDevice(ctx context.Context, d *domain.HardwareDevice) error {
	if d == nil { return errors.New("nil device") }
	now := uint64(time.Now().UnixMilli())
	if d.CreatedAtMs == 0 { d.CreatedAtMs = now }
	d.UpdatedAtMs = now
	return s.DeviceRepo.UpsertDevice(ctx, d)
}

// ListTTSVoices returns voices for a space (including global voices where space_id is NULL).
func (s *ApplicationService) ListTTSVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*domain.TTSVoice, int64, error) {
	return s.VoiceRepo.ListVoices(ctx, spaceID, provider, language, gender, page, pageSize)
}

// UpsertAppTTS saves app-level tts settings (unique by app_id).
func (s *ApplicationService) UpsertAppTTS(ctx context.Context, cfg *domain.AppTTSSettings) error {
	if cfg == nil { return errors.New("nil app tts") }
	now := uint64(time.Now().UnixMilli())
	if cfg.CreatedAtMs == 0 { cfg.CreatedAtMs = now }
	cfg.UpdatedAtMs = now
	return s.SettingsRepo.UpsertAppTTS(ctx, cfg)
}

// UpsertHardwareTTS saves device-level tts settings (unique by device_id).
func (s *ApplicationService) UpsertHardwareTTS(ctx context.Context, cfg *domain.HardwareTTSSettings) error {
	if cfg == nil { return errors.New("nil hardware tts") }
	now := uint64(time.Now().UnixMilli())
	if cfg.CreatedAtMs == 0 { cfg.CreatedAtMs = now }
	cfg.UpdatedAtMs = now
	return s.SettingsRepo.UpsertHardwareTTS(ctx, cfg)
}

// GetEffectiveDeviceTTS returns the effective provider/model/voice and source priority: device > app > default.
func (s *ApplicationService) GetEffectiveDeviceTTS(ctx context.Context, deviceID string, appID *uint64) (*domain.EffectiveTTS, error) {
	return s.Domain.GetEffectiveTTS(ctx, deviceID, appID)
}

// GetVoiceSampleURL finds a sample_url for provider+voice (space scoped first, then global)
func (s *ApplicationService) GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error) {
	return s.VoiceRepo.GetVoiceSampleURL(ctx, provider, voice, spaceID)
}