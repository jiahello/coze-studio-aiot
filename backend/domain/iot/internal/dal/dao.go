package dal

import (
	"context"

	"gorm.io/gorm"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
)

type DeviceDAO struct { db *gorm.DB }

type VoiceDAO struct { db *gorm.DB }

type SettingsDAO struct { db *gorm.DB }

func NewDeviceDAO(db *gorm.DB) *DeviceDAO   { return &DeviceDAO{db: db} }
func NewVoiceDAO(db *gorm.DB) *VoiceDAO     { return &VoiceDAO{db: db} }
func NewSettingsDAO(db *gorm.DB) *SettingsDAO { return &SettingsDAO{db: db} }

// device
func (d *DeviceDAO) ListDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*domain.HardwareDevice, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 200 { pageSize = 20 }
	var (
		listPO  []HardwareDevice
		total int64
	)
	db := d.db.WithContext(ctx).Model(&HardwareDevice{}).Where("space_id = ?", spaceID)
	if keyword != "" {
		db = db.Where("device_id LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	if err := db.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&listPO).Error; err != nil { return nil, 0, err }
	res := make([]*domain.HardwareDevice, 0, len(listPO))
	for i := range listPO {
		po := &listPO[i]
		res = append(res, &domain.HardwareDevice{
			ID: po.ID, DeviceID: po.DeviceID, Name: po.Name, Description: po.Description,
			AppID: po.AppID, SpaceID: po.SpaceID, CreatedUserID: po.CreatedUserID, UpdatedUserID: po.UpdatedUserID,
			CreatedAtMs: po.CreatedAtMs, UpdatedAtMs: po.UpdatedAtMs, FirmwareVersion: po.FirmwareVersion,
			MacAddress: po.MacAddress, Status: po.Status, LastPingAtMs: po.LastPingAtMs, VerifyCode: po.VerifyCode,
			FromEndUserID: po.FromEndUserID, FromAccountID: po.FromAccountID,
		})
	}
	return res, total, nil
}

func (d *DeviceDAO) UpsertDevice(ctx context.Context, dev *domain.HardwareDevice) error {
	existing := &HardwareDevice{}
	err := d.db.WithContext(ctx).Where("space_id = ? AND device_id = ?", dev.SpaceID, dev.DeviceID).First(existing).Error
	if err == nil {
		upd := &HardwareDevice{ID: existing.ID, DeviceID: dev.DeviceID, Name: dev.Name, Description: dev.Description, AppID: dev.AppID,
			SpaceID: dev.SpaceID, CreatedUserID: dev.CreatedUserID, UpdatedUserID: dev.UpdatedUserID,
			CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: dev.UpdatedAtMs, FirmwareVersion: dev.FirmwareVersion,
			MacAddress: dev.MacAddress, Status: dev.Status, LastPingAtMs: dev.LastPingAtMs, VerifyCode: dev.VerifyCode,
			FromEndUserID: dev.FromEndUserID, FromAccountID: dev.FromAccountID,
		}
		return d.db.WithContext(ctx).Model(&HardwareDevice{}).Where("id = ?", existing.ID).Updates(upd).Error
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		po := &HardwareDevice{DeviceID: dev.DeviceID, Name: dev.Name, Description: dev.Description, AppID: dev.AppID,
			SpaceID: dev.SpaceID, CreatedUserID: dev.CreatedUserID, UpdatedUserID: dev.UpdatedUserID,
			CreatedAtMs: dev.CreatedAtMs, UpdatedAtMs: dev.UpdatedAtMs, FirmwareVersion: dev.FirmwareVersion,
			MacAddress: dev.MacAddress, Status: dev.Status, LastPingAtMs: dev.LastPingAtMs, VerifyCode: dev.VerifyCode,
			FromEndUserID: dev.FromEndUserID, FromAccountID: dev.FromAccountID,
		}
		return d.db.WithContext(ctx).Create(po).Error
	}
	return err
}

// voice
func (v *VoiceDAO) ListVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*domain.TTSVoice, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 200 { pageSize = 20 }
	var (
		listPO  []TTSVoice
		total int64
	)
	db := v.db.WithContext(ctx).Model(&TTSVoice{})
	if spaceID != nil { db = db.Where("space_id = ? OR space_id IS NULL", *spaceID) }
	if provider != "" { db = db.Where("provider = ?", provider) }
	if language != "" { db = db.Where("language = ?", language) }
	if gender != "" { db = db.Where("gender = ?", gender) }
	if err := db.Count(&total).Error; err != nil { return nil, 0, err }
	if err := db.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&listPO).Error; err != nil { return nil, 0, err }
	res := make([]*domain.TTSVoice, 0, len(listPO))
	for i := range listPO {
		po := &listPO[i]
		res = append(res, &domain.TTSVoice{
			ID: po.ID, Provider: po.Provider, Model: po.Model, VoiceType: po.VoiceType, Name: po.Name, VoiceCode: po.VoiceCode,
			Description: po.Description, Gender: po.Gender, Language: po.Language, Scenario: po.Scenario, SoundQuality: po.SoundQuality,
			SampleRate: po.SampleRate, TimestampSup: po.TimestampSup, ErhuaSupport: po.ErhuaSupport, SampleURL: po.SampleURL,
			CreatedAtMs: po.CreatedAtMs, UpdatedAtMs: po.UpdatedAtMs,
			LanguageName: po.LanguageName, VoiceID: po.VoiceID, ModelType: po.ModelType, PreviewText: po.PreviewText,
			EmotionSupport: po.EmotionSupport, Emotions: po.Emotions, FromAccountID: po.FromAccountID, APIType: po.APIType, SpaceID: po.SpaceID,
		})
	}
	return res, total, nil
}

func (v *VoiceDAO) GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error) {
	var m TTSVoice
	db := v.db.WithContext(ctx).Model(&TTSVoice{}).Where("provider = ? AND voice_code = ?", provider, voice)
	if spaceID != nil { db = db.Where("space_id = ? OR space_id IS NULL", *spaceID) }
	if err := db.Order("space_id IS NULL").First(&m).Error; err != nil { return "", err }
	if m.SampleURL != nil { return *m.SampleURL, nil }
	return "", nil
}

// settings
func (s *SettingsDAO) UpsertAppTTS(ctx context.Context, cfg *domain.AppTTSSettings) error {
	existing := &AppTTSSettings{}
	err := s.db.WithContext(ctx).Where("app_id = ?", cfg.AppID).First(existing).Error
	if err == nil {
		upd := &AppTTSSettings{ID: existing.ID, AppID: cfg.AppID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs}
		return s.db.WithContext(ctx).Model(&AppTTSSettings{}).Where("id = ?", existing.ID).Updates(upd).Error
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		po := &AppTTSSettings{AppID: cfg.AppID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: cfg.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs}
		return s.db.WithContext(ctx).Create(po).Error
	}
	return err
}

func (s *SettingsDAO) UpsertHardwareTTS(ctx context.Context, cfg *domain.HardwareTTSSettings) error {
	existing := &HardwareTTSSettings{}
	err := s.db.WithContext(ctx).Where("device_id = ?", cfg.DeviceID).First(existing).Error
	if err == nil {
		upd := &HardwareTTSSettings{ID: existing.ID, DeviceID: cfg.DeviceID, HardwareDeviceID: cfg.HardwareDeviceID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs, IsDeleted: cfg.IsDeleted}
		return s.db.WithContext(ctx).Model(&HardwareTTSSettings{}).Where("id = ?", existing.ID).Updates(upd).Error
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		po := &HardwareTTSSettings{DeviceID: cfg.DeviceID, HardwareDeviceID: cfg.HardwareDeviceID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: cfg.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs, IsDeleted: cfg.IsDeleted}
		return s.db.WithContext(ctx).Create(po).Error
	}
	return err
}

func (s *SettingsDAO) GetAppTTSByAppID(ctx context.Context, appID uint64) (*domain.AppTTSSettings, error) {
	var a AppTTSSettings
	if err := s.db.WithContext(ctx).Where("app_id = ?", appID).First(&a).Error; err != nil { return nil, err }
	return &domain.AppTTSSettings{ID: a.ID, AppID: a.AppID, Provider: a.Provider, Model: a.Model, Voice: a.Voice, VoiceRef: a.VoiceRef, CreatedUserID: a.CreatedUserID, UpdatedUserID: a.UpdatedUserID, CreatedAtMs: a.CreatedAtMs, UpdatedAtMs: a.UpdatedAtMs}, nil
}

func (s *SettingsDAO) GetHardwareTTSByDeviceID(ctx context.Context, deviceID string) (*domain.HardwareTTSSettings, error) {
	var h HardwareTTSSettings
	if err := s.db.WithContext(ctx).Where("device_id = ? AND is_deleted = 0", deviceID).First(&h).Error; err != nil { return nil, err }
	return &domain.HardwareTTSSettings{ID: h.ID, DeviceID: h.DeviceID, HardwareDeviceID: h.HardwareDeviceID, Provider: h.Provider, Model: h.Model, Voice: h.Voice, VoiceRef: h.VoiceRef, CreatedUserID: h.CreatedUserID, UpdatedUserID: h.UpdatedUserID, CreatedAtMs: h.CreatedAtMs, UpdatedAtMs: h.UpdatedAtMs, IsDeleted: h.IsDeleted}, nil
}