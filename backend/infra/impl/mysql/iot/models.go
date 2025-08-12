package iot

import "time"

// 注意：以下为持久化模型（PO），仅供本目录内部使用

 type HardwareDevicePO struct {
	ID              uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID        string    `gorm:"column:device_id;size:64;not null"`
	Name            string    `gorm:"column:name;size:255;not null"`
	Description     string    `gorm:"column:description"`
	AppID           *uint64   `gorm:"column:app_id"`
	SpaceID         uint64    `gorm:"column:space_id;not null"`
	CreatedUserID   uint64    `gorm:"column:created_user_id;not null"`
	UpdatedUserID   uint64    `gorm:"column:updated_user_id;not null"`
	CreatedAtMs     uint64    `gorm:"column:created_at;not null"`
	UpdatedAtMs     uint64    `gorm:"column:updated_at;not null"`
	FirmwareVersion *string   `gorm:"column:firmware_version;size:64"`
	MacAddress      *string   `gorm:"column:mac_address;size:64"`
	Status          string    `gorm:"column:status;size:255;not null"`
	LastPingAtMs    *uint64   `gorm:"column:last_ping_at"`
	VerifyCode      *string   `gorm:"column:verify_code;size:8"`
	FromEndUserID   *string   `gorm:"column:from_end_user_id"`
	FromAccountID   *string   `gorm:"column:from_account_id"`
 }

 func (HardwareDevicePO) TableName() string { return "hardware_device" }

 type TTSVoicePO struct {
	ID             uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Provider       string     `gorm:"column:provider;size:255;not null"`
	Model          *string    `gorm:"column:model;size:255"`
	VoiceType      string     `gorm:"column:voice_type;size:32;not null"`
	Name           string     `gorm:"column:name;size:255;not null"`
	VoiceCode      string     `gorm:"column:voice_code;size:255;not null"`
	Description    *string    `gorm:"column:description"`
	Gender         *string    `gorm:"column:gender;size:32"`
	Language       *string    `gorm:"column:language;size:64"`
	Scenario       *string    `gorm:"column:scenario;size:255"`
	SoundQuality   *string    `gorm:"column:sound_quality;size:64"`
	SampleRate     *string    `gorm:"column:sample_rate;size:64"`
	TimestampSup   *string    `gorm:"column:timestamp_support;size:16"`
	ErhuaSupport   *string    `gorm:"column:erhua_support;size:16"`
	SampleURL      *string    `gorm:"column:sample_url;size:512"`
	CreatedAtMs    uint64     `gorm:"column:created_at;not null"`
	UpdatedAtMs    uint64     `gorm:"column:updated_at;not null"`
	SourceFileID   *string    `gorm:"column:source_file_id"`
	SampleFileID   *string    `gorm:"column:sample_file_id"`
	IsDeleted      bool       `gorm:"column:is_deleted;not null"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
	DeletedByID    *uint64    `gorm:"column:deleted_by_id"`
	LanguageName   *string    `gorm:"column:language_name;size:255"`
	VoiceID        *string    `gorm:"column:voice_id;size:255"`
	ModelType      *string    `gorm:"column:model_type;size:64"`
	PreviewText    *string    `gorm:"column:preview_text"`
	EmotionSupport *bool      `gorm:"column:emotion_support"`
	Emotions       *string    `gorm:"column:emotions"`
	FromAccountID  *string    `gorm:"column:from_account_id"`
	APIType        *string    `gorm:"column:api_type;size:64"`
	SpaceID        *uint64    `gorm:"column:space_id"`
 }

 func (TTSVoicePO) TableName() string { return "tts_voice" }

 type AppTTSSettingsPO struct {
	ID            uint64  `gorm:"column:id;primaryKey;autoIncrement"`
	AppID         uint64  `gorm:"column:app_id;not null"`
	Provider      string  `gorm:"column:provider;size:255;not null"`
	Model         string  `gorm:"column:model;size:255;not null"`
	Voice         string  `gorm:"column:voice;size:255;not null"`
	VoiceRef      *uint64 `gorm:"column:voice_ref"`
	CreatedUserID uint64  `gorm:"column:created_user_id;not null"`
	UpdatedUserID uint64  `gorm:"column:updated_user_id;not null"`
	CreatedAtMs   uint64  `gorm:"column:created_at;not null"`
	UpdatedAtMs   uint64  `gorm:"column:updated_at;not null"`
 }

 func (AppTTSSettingsPO) TableName() string { return "app_tts_settings" }

 type HardwareTTSSettingsPO struct {
	ID               uint64   `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID         string   `gorm:"column:device_id;size:255;not null"`
	HardwareDeviceID *uint64  `gorm:"column:hardware_device_id"`
	Provider         string   `gorm:"column:provider;size:255;not null"`
	Model            string   `gorm:"column:model;size:255;not null"`
	Voice            string   `gorm:"column:voice;size:255;not null"`
	VoiceRef         *uint64  `gorm:"column:voice_ref"`
	CreatedUserID    uint64   `gorm:"column:created_user_id;not null"`
	UpdatedUserID    uint64   `gorm:"column:updated_user_id;not null"`
	CreatedAtMs      uint64   `gorm:"column:created_at;not null"`
	UpdatedAtMs      uint64   `gorm:"column:updated_at;not null"`
	IsDeleted        bool     `gorm:"column:is_deleted;not null"`
 }

 func (HardwareTTSSettingsPO) TableName() string { return "hardware_tts_settings" }