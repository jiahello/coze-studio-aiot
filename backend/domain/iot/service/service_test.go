package service

import (
	"context"
	"testing"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
	repo "github.com/coze-dev/coze-studio/backend/domain/iot/repository"
)

type mockSettingsRepo struct {
	app    *domain.AppTTSSettings
	hw     *domain.HardwareTTSSettings
}

func (m *mockSettingsRepo) UpsertAppTTS(ctx context.Context, cfg *domain.AppTTSSettings) error { m.app = cfg; return nil }
func (m *mockSettingsRepo) UpsertHardwareTTS(ctx context.Context, cfg *domain.HardwareTTSSettings) error { m.hw = cfg; return nil }
func (m *mockSettingsRepo) GetAppTTSByAppID(ctx context.Context, appID uint64) (*domain.AppTTSSettings, error) { return m.app, nil }
func (m *mockSettingsRepo) GetHardwareTTSByDeviceID(ctx context.Context, deviceID string) (*domain.HardwareTTSSettings, error) { return m.hw, nil }

// stubs to satisfy constructor
type noopDeviceRepo struct{}
func (noopDeviceRepo) ListDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*domain.HardwareDevice, int64, error) { return nil, 0, nil }
func (noopDeviceRepo) UpsertDevice(ctx context.Context, d *domain.HardwareDevice) error { return nil }

type noopVoiceRepo struct{}
func (noopVoiceRepo) ListVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*domain.TTSVoice, int64, error) { return nil, 0, nil }
func (noopVoiceRepo) GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error) { return "", nil }

var _ repo.TTSSettingsRepository = (*mockSettingsRepo)(nil)
var _ repo.DeviceRepository = (*noopDeviceRepo)(nil)
var _ repo.VoiceRepository = (*noopVoiceRepo)(nil)

func TestGetEffectiveTTS_Priority(t *testing.T) {
	ctx := context.Background()
	settings := &mockSettingsRepo{}
	svc := NewService(noopDeviceRepo{}, noopVoiceRepo{}, settings)

	// default
	res, err := svc.GetEffectiveTTS(ctx, "", nil)
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res.Source != "default" { t.Fatalf("expect default, got %s", res.Source) }

	// app-level
	settings.app = &domain.AppTTSSettings{Provider: "p1", Model: "m1", Voice: "v1"}
	res, err = svc.GetEffectiveTTS(ctx, "dev1", ptrU64(123))
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res.Source != "app" || res.Provider != "p1" { t.Fatalf("expect app p1, got %+v", res) }

	// device-level overrides app
	settings.hw = &domain.HardwareTTSSettings{Provider: "p2", Model: "m2", Voice: "v2"}
	res, err = svc.GetEffectiveTTS(ctx, "dev1", ptrU64(123))
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res.Source != "device" || res.Provider != "p2" { t.Fatalf("expect device p2, got %+v", res) }
}

func ptrU64(v uint64) *uint64 { return &v }