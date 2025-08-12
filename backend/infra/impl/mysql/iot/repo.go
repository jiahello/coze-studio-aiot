package iot

import (
	"context"

	"gorm.io/gorm"

	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
)

type RepositoryBundle struct {
	Devices  domain.DeviceRepository
	Voices   domain.VoiceRepository
	Settings domain.TTSSettingsRepository
}

type deviceRepo struct { db *gorm.DB }

type voiceRepo struct { db *gorm.DB }

type settingsRepo struct { db *gorm.DB }

func NewRepositories(db *gorm.DB) *RepositoryBundle {
	return &RepositoryBundle{
		Devices:  &deviceRepo{db: db},
		Voices:   &voiceRepo{db: db},
		Settings: &settingsRepo{db: db},
	}
}

// DeviceRepository
func (r *deviceRepo) ListDevices(ctx context.Context, spaceID uint64, page, pageSize int, keyword string) ([]*domain.HardwareDevice, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 200 { pageSize = 20 }
	var (
		listPO  []HardwareDevicePO
		total int64
	)
	db := r.db.WithContext(ctx).Model(&HardwareDevicePO{}).Where("space_id = ?", spaceID)
	if keyword != "" {
		db = db.Where("device_id LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).Find(&listPO).Error; err != nil {
		return nil, 0, err
	}
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

func (r *deviceRepo) UpsertDevice(ctx context.Context, d *domain.HardwareDevice) error {
	now := d.UpdatedAtMs
	if now == 0 {
		// 由业务层传入时间；如未传入则数据库触发器/默认值，或此处保持 0
	}
	if d.ID == 0 {
		// 查询是否存在
		existing := &HardwareDevicePO{}
		err := r.db.WithContext(ctx).Where("space_id = ? AND device_id = ?", d.SpaceID, d.DeviceID).First(existing).Error
		if err == nil {
			// 更新
			upd := &HardwareDevicePO{
				ID: existing.ID,
				DeviceID: d.DeviceID, Name: d.Name, Description: d.Description, AppID: d.AppID,
				SpaceID: d.SpaceID, CreatedUserID: d.CreatedUserID, UpdatedUserID: d.UpdatedUserID,
				CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: d.UpdatedAtMs,
				FirmwareVersion: d.FirmwareVersion, MacAddress: d.MacAddress, Status: d.Status,
				LastPingAtMs: d.LastPingAtMs, VerifyCode: d.VerifyCode,
				FromEndUserID: d.FromEndUserID, FromAccountID: d.FromAccountID,
			}
			return r.db.WithContext(ctx).Model(&HardwareDevicePO{}).Where("id = ?", existing.ID).Updates(upd).Error
		}
		if err != nil && err == gorm.ErrRecordNotFound {
			po := &HardwareDevicePO{
				DeviceID: d.DeviceID, Name: d.Name, Description: d.Description, AppID: d.AppID,
				SpaceID: d.SpaceID, CreatedUserID: d.CreatedUserID, UpdatedUserID: d.UpdatedUserID,
				CreatedAtMs: d.CreatedAtMs, UpdatedAtMs: d.UpdatedAtMs,
				FirmwareVersion: d.FirmwareVersion, MacAddress: d.MacAddress, Status: d.Status,
				LastPingAtMs: d.LastPingAtMs, VerifyCode: d.VerifyCode,
				FromEndUserID: d.FromEndUserID, FromAccountID: d.FromAccountID,
			}
			return r.db.WithContext(ctx).Create(po).Error
		}
		return err
	}
	// 更新
	upd := &HardwareDevicePO{
		ID: d.ID,
		DeviceID: d.DeviceID, Name: d.Name, Description: d.Description, AppID: d.AppID,
		SpaceID: d.SpaceID, CreatedUserID: d.CreatedUserID, UpdatedUserID: d.UpdatedUserID,
		CreatedAtMs: d.CreatedAtMs, UpdatedAtMs: d.UpdatedAtMs,
		FirmwareVersion: d.FirmwareVersion, MacAddress: d.MacAddress, Status: d.Status,
		LastPingAtMs: d.LastPingAtMs, VerifyCode: d.VerifyCode,
		FromEndUserID: d.FromEndUserID, FromAccountID: d.FromAccountID,
	}
	return r.db.WithContext(ctx).Model(&HardwareDevicePO{}).Where("id = ?", d.ID).Updates(upd).Error
}

// VoiceRepository
func (r *voiceRepo) ListVoices(ctx context.Context, spaceID *uint64, provider, language, gender string, page, pageSize int) ([]*domain.TTSVoice, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 200 { pageSize = 20 }
	var (
		listPO  []TTSVoicePO
		total int64
	)
	db := r.db.WithContext(ctx).Model(&TTSVoicePO{})
	if spaceID != nil {
		db = db.Where("space_id = ? OR space_id IS NULL", *spaceID)
	}
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
			CreatedAtMs: po.CreatedAtMs, UpdatedAtMs: po.UpdatedAtMs, SourceFileID: po.SourceFileID, SampleFileID: po.SampleFileID,
			IsDeleted: po.IsDeleted, LanguageName: po.LanguageName, VoiceID: po.VoiceID, ModelType: po.ModelType, PreviewText: po.PreviewText,
			EmotionSupport: po.EmotionSupport, Emotions: po.Emotions, FromAccountID: po.FromAccountID, APIType: po.APIType, SpaceID: po.SpaceID,
		})
	}
	return res, total, nil
}

func (r *voiceRepo) GetVoiceSampleURL(ctx context.Context, provider, voice string, spaceID *uint64) (string, error) {
	var v TTSVoicePO
	db := r.db.WithContext(ctx).Model(&TTSVoicePO{}).Where("provider = ? AND voice_code = ?", provider, voice)
	if spaceID != nil {
		db = db.Where("space_id = ? OR space_id IS NULL", *spaceID)
	}
	if err := db.Order("space_id IS NULL").First(&v).Error; err != nil { return "", err }
	if v.SampleURL != nil { return *v.SampleURL, nil }
	return "", nil
}

// SettingsRepository
func (r *settingsRepo) UpsertAppTTS(ctx context.Context, cfg *domain.AppTTSSettings) error {
	existing := &AppTTSSettingsPO{}
	err := r.db.WithContext(ctx).Where("app_id = ?", cfg.AppID).First(existing).Error
	if err == nil {
		upd := &AppTTSSettingsPO{ID: existing.ID, AppID: cfg.AppID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs}
		return r.db.WithContext(ctx).Model(&AppTTSSettingsPO{}).Where("id = ?", existing.ID).Updates(upd).Error
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		po := &AppTTSSettingsPO{AppID: cfg.AppID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: cfg.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs}
		return r.db.WithContext(ctx).Create(po).Error
	}
	return err
}

func (r *settingsRepo) UpsertHardwareTTS(ctx context.Context, cfg *domain.HardwareTTSSettings) error {
	existing := &HardwareTTSSettingsPO{}
	err := r.db.WithContext(ctx).Where("device_id = ?", cfg.DeviceID).First(existing).Error
	if err == nil {
		upd := &HardwareTTSSettingsPO{ID: existing.ID, DeviceID: cfg.DeviceID, HardwareDeviceID: cfg.HardwareDeviceID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: existing.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs, IsDeleted: cfg.IsDeleted}
		return r.db.WithContext(ctx).Model(&HardwareTTSSettingsPO{}).Where("id = ?", existing.ID).Updates(upd).Error
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		po := &HardwareTTSSettingsPO{DeviceID: cfg.DeviceID, HardwareDeviceID: cfg.HardwareDeviceID, Provider: cfg.Provider, Model: cfg.Model, Voice: cfg.Voice, VoiceRef: cfg.VoiceRef, CreatedUserID: cfg.CreatedUserID, UpdatedUserID: cfg.UpdatedUserID, CreatedAtMs: cfg.CreatedAtMs, UpdatedAtMs: cfg.UpdatedAtMs, IsDeleted: cfg.IsDeleted}
		return r.db.WithContext(ctx).Create(po).Error
	}
	return err
}

func (r *settingsRepo) GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*domain.EffectiveTTS, error) {
	// device level
	var h HardwareTTSSettingsPO
	if err := r.db.WithContext(ctx).Where("device_id = ? AND is_deleted = 0", deviceID).First(&h).Error; err == nil {
		return &domain.EffectiveTTS{Provider: h.Provider, Model: h.Model, Voice: h.Voice, Source: "device"}, nil
	}
	// app level
	if appID != nil {
		var a AppTTSSettingsPO
		if err := r.db.WithContext(ctx).Where("app_id = ?", *appID).First(&a).Error; err == nil {
			return &domain.EffectiveTTS{Provider: a.Provider, Model: a.Model, Voice: a.Voice, Source: "app"}, nil
		}
	}
	return &domain.EffectiveTTS{Provider: "doubao", Model: "speech-1", Voice: "doubao-standard", Source: "default"}, nil
}