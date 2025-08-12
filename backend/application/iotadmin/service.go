package iotadmin

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ListHardwareDevices returns paginated devices under a space.
func (s *ApplicationService) ListHardwareDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*HardwareDevice, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	var (
		list  []*HardwareDevice
		total int64
	)
	db := s.DB.WithContext(ctx).Model(&HardwareDevice{}).Where("space_id = ?", spaceID)
	if keyword != "" {
		db = db.Where("device_id LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpsertHardwareDevice creates or updates a device by (space_id, device_id) or by id.
func (s *ApplicationService) UpsertHardwareDevice(ctx context.Context, d *HardwareDevice) error {
	if d == nil {
		return errors.New("nil device")
	}
	now := uint64(time.Now().UnixMilli())
	if d.ID == 0 {
		// try find by (space_id, device_id)
		existing := &HardwareDevice{}
		err := s.DB.WithContext(ctx).Where("space_id = ? AND device_id = ?", d.SpaceID, d.DeviceID).First(existing).Error
		if err == nil {
			// update
			d.ID = existing.ID
			d.CreatedAtMs = existing.CreatedAtMs
			if d.CreatedAtMs == 0 {
				d.CreatedAtMs = now
			}
			d.UpdatedAtMs = now
			return s.DB.WithContext(ctx).Model(&HardwareDevice{}).Where("id = ?", d.ID).Updates(d).Error
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			d.CreatedAtMs = now
			d.UpdatedAtMs = now
			return s.DB.WithContext(ctx).Create(d).Error
		}
		return err
	}
	// update by id
	d.UpdatedAtMs = now
	return s.DB.WithContext(ctx).Model(&HardwareDevice{}).Where("id = ?", d.ID).Updates(d).Error
}

// ListTTSVoices returns voices for a space (including global voices where space_id is NULL).
func (s *ApplicationService) ListTTSVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*TTSVoice, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	var (
		list  []*TTSVoice
		total int64
	)
	db := s.DB.WithContext(ctx).Model(&TTSVoice{})
	if spaceID != nil {
		db = db.Where("space_id = ? OR space_id IS NULL", *spaceID)
	}
	if provider != "" {
		db = db.Where("provider = ?", provider)
	}
	if language != "" {
		db = db.Where("language = ?", language)
	}
	if gender != "" {
		db = db.Where("gender = ?", gender)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpsertAppTTS saves app-level tts settings (unique by app_id).
func (s *ApplicationService) UpsertAppTTS(ctx context.Context, cfg *AppTTSSettings) error {
	if cfg == nil {
		return errors.New("nil app tts")
	}
	now := uint64(time.Now().UnixMilli())
	existing := &AppTTSSettings{}
	err := s.DB.WithContext(ctx).Where("app_id = ?", cfg.AppID).First(existing).Error
	if err == nil {
		cfg.ID = existing.ID
		cfg.CreatedAtMs = existing.CreatedAtMs
		cfg.UpdatedAtMs = now
		return s.DB.WithContext(ctx).Model(&AppTTSSettings{}).Where("id = ?", cfg.ID).Updates(cfg).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg.CreatedAtMs = now
		cfg.UpdatedAtMs = now
		return s.DB.WithContext(ctx).Create(cfg).Error
	}
	return err
}

// UpsertHardwareTTS saves device-level tts settings (unique by device_id).
func (s *ApplicationService) UpsertHardwareTTS(ctx context.Context, cfg *HardwareTTSSettings) error {
	if cfg == nil {
		return errors.New("nil hardware tts")
	}
	now := uint64(time.Now().UnixMilli())
	existing := &HardwareTTSSettings{}
	err := s.DB.WithContext(ctx).Where("device_id = ?", cfg.DeviceID).First(existing).Error
	if err == nil {
		cfg.ID = existing.ID
		cfg.CreatedAtMs = existing.CreatedAtMs
		cfg.UpdatedAtMs = now
		return s.DB.WithContext(ctx).Model(&HardwareTTSSettings{}).Where("id = ?", cfg.ID).Updates(cfg).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cfg.CreatedAtMs = now
		cfg.UpdatedAtMs = now
		return s.DB.WithContext(ctx).Create(cfg).Error
	}
	return err
}

// GetEffectiveDeviceTTS returns the effective provider/model/voice and source priority: device > app > default.
type EffectiveTTS struct {
	Provider string
	Model    string
	Voice    string
	Source   string // device|app|default
}

func (s *ApplicationService) GetEffectiveDeviceTTS(ctx context.Context, deviceID string, appID *uint64) (*EffectiveTTS, error) {
	// device level
	var h HardwareTTSSettings
	if err := s.DB.WithContext(ctx).Where("device_id = ? AND is_deleted = 0", deviceID).First(&h).Error; err == nil {
		return &EffectiveTTS{Provider: h.Provider, Model: h.Model, Voice: h.Voice, Source: "device"}, nil
	}
	// app level
	if appID != nil {
		var a AppTTSSettings
		if err := s.DB.WithContext(ctx).Where("app_id = ?", *appID).First(&a).Error; err == nil {
			return &EffectiveTTS{Provider: a.Provider, Model: a.Model, Voice: a.Voice, Source: "app"}, nil
		}
	}
	// default
	return &EffectiveTTS{Provider: "doubao", Model: "speech-1", Voice: "doubao-standard", Source: "default"}, nil
}

// GetVoiceSampleURL finds a sample_url for provider+voice (space scoped first, then global)
func (s *ApplicationService) GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error) {
	var v TTSVoice
	db := s.DB.WithContext(ctx).Model(&TTSVoice{}).Where("provider = ? AND voice_code = ?", provider, voice)
	if spaceID != nil {
		db = db.Where("space_id = ? OR space_id IS NULL", *spaceID)
	}
	if err := db.Order("space_id IS NULL").First(&v).Error; err != nil {
		return "", err
	}
	if v.SampleURL != nil {
		return *v.SampleURL, nil
	}
	return "", nil
}