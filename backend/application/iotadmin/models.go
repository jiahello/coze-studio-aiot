package iotadmin

import "time"

// HardwareDevice maps to table `hardware_device`.
type HardwareDevice struct {
	ID              uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeviceID        string    `gorm:"column:device_id;size:64;not null" json:"device_id"`
	Name            string    `gorm:"column:name;size:255;not null" json:"name"`
	Description     string    `gorm:"column:description" json:"description"`
	AppID           *uint64   `gorm:"column:app_id" json:"app_id,omitempty"`
	SpaceID         uint64    `gorm:"column:space_id;not null" json:"space_id"`
	CreatedUserID   uint64    `gorm:"column:created_user_id;not null" json:"created_user_id"`
	UpdatedUserID   uint64    `gorm:"column:updated_user_id;not null" json:"updated_user_id"`
	CreatedAtMs     uint64    `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAtMs     uint64    `gorm:"column:updated_at;not null" json:"updated_at"`
	FirmwareVersion *string   `gorm:"column:firmware_version;size:64" json:"firmware_version,omitempty"`
	MacAddress      *string   `gorm:"column:mac_address;size:64" json:"mac_address,omitempty"`
	Status          string    `gorm:"column:status;size:255;not null" json:"status"`
	LastPingAtMs    *uint64   `gorm:"column:last_ping_at" json:"last_ping_at,omitempty"`
	VerifyCode      *string   `gorm:"column:verify_code;size:8" json:"verify_code,omitempty"`
	FromEndUserID   *string   `gorm:"column:from_end_user_id" json:"from_end_user_id,omitempty"`
	FromAccountID   *string   `gorm:"column:from_account_id" json:"from_account_id,omitempty"`
}

func (HardwareDevice) TableName() string { return "hardware_device" }

// TTSVoice maps to table `tts_voice`.
type TTSVoice struct {
	ID             uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Provider       string     `gorm:"column:provider;size:255;not null" json:"provider"`
	Model          *string    `gorm:"column:model;size:255" json:"model,omitempty"`
	VoiceType      string     `gorm:"column:voice_type;size:32;not null" json:"voice_type"`
	Name           string     `gorm:"column:name;size:255;not null" json:"name"`
	VoiceCode      string     `gorm:"column:voice_code;size:255;not null" json:"voice_code"`
	Description    *string    `gorm:"column:description" json:"description,omitempty"`
	Gender         *string    `gorm:"column:gender;size:32" json:"gender,omitempty"`
	Language       *string    `gorm:"column:language;size:64" json:"language,omitempty"`
	Scenario       *string    `gorm:"column:scenario;size:255" json:"scenario,omitempty"`
	SoundQuality   *string    `gorm:"column:sound_quality;size:64" json:"sound_quality,omitempty"`
	SampleRate     *string    `gorm:"column:sample_rate;size:64" json:"sample_rate,omitempty"`
	TimestampSup   *string    `gorm:"column:timestamp_support;size:16" json:"timestamp_support,omitempty"`
	ErhuaSupport   *string    `gorm:"column:erhua_support;size:16" json:"erhua_support,omitempty"`
	SampleURL      *string    `gorm:"column:sample_url;size:512" json:"sample_url,omitempty"`
	CreatedAtMs    uint64     `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAtMs    uint64     `gorm:"column:updated_at;not null" json:"updated_at"`
	SourceFileID   *string    `gorm:"column:source_file_id" json:"source_file_id,omitempty"`
	SampleFileID   *string    `gorm:"column:sample_file_id" json:"sample_file_id,omitempty"`
	IsDeleted      bool       `gorm:"column:is_deleted;not null" json:"is_deleted"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	DeletedByID    *uint64    `gorm:"column:deleted_by_id" json:"deleted_by_id,omitempty"`
	LanguageName   *string    `gorm:"column:language_name;size:255" json:"language_name,omitempty"`
	VoiceID        *string    `gorm:"column:voice_id;size:255" json:"voice_id,omitempty"`
	ModelType      *string    `gorm:"column:model_type;size:64" json:"model_type,omitempty"`
	PreviewText    *string    `gorm:"column:preview_text" json:"preview_text,omitempty"`
	EmotionSupport *bool      `gorm:"column:emotion_support" json:"emotion_support,omitempty"`
	Emotions       *string    `gorm:"column:emotions" json:"emotions,omitempty"`
	FromAccountID  *string    `gorm:"column:from_account_id" json:"from_account_id,omitempty"`
	APIType        *string    `gorm:"column:api_type;size:64" json:"api_type,omitempty"`
	SpaceID        *uint64    `gorm:"column:space_id" json:"space_id,omitempty"`
}

func (TTSVoice) TableName() string { return "tts_voice" }

// AppTTSSettings maps to table `app_tts_settings`.
type AppTTSSettings struct {
	ID            uint64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AppID         uint64  `gorm:"column:app_id;not null" json:"app_id"`
	Provider      string  `gorm:"column:provider;size:255;not null" json:"provider"`
	Model         string  `gorm:"column:model;size:255;not null" json:"model"`
	Voice         string  `gorm:"column:voice;size:255;not null" json:"voice"`
	VoiceRef      *uint64 `gorm:"column:voice_ref" json:"voice_ref,omitempty"`
	CreatedUserID uint64  `gorm:"column:created_user_id;not null" json:"created_user_id"`
	UpdatedUserID uint64  `gorm:"column:updated_user_id;not null" json:"updated_user_id"`
	CreatedAtMs   uint64  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAtMs   uint64  `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (AppTTSSettings) TableName() string { return "app_tts_settings" }

// HardwareTTSSettings maps to table `hardware_tts_settings`.
type HardwareTTSSettings struct {
	ID              uint64   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeviceID        string   `gorm:"column:device_id;size:255;not null" json:"device_id"`
	HardwareDeviceID *uint64 `gorm:"column:hardware_device_id" json:"hardware_device_id,omitempty"`
	Provider        string   `gorm:"column:provider;size:255;not null" json:"provider"`
	Model           string   `gorm:"column:model;size:255;not null" json:"model"`
	Voice           string   `gorm:"column:voice;size:255;not null" json:"voice"`
	VoiceRef        *uint64  `gorm:"column:voice_ref" json:"voice_ref,omitempty"`
	CreatedUserID   uint64   `gorm:"column:created_user_id;not null" json:"created_user_id"`
	UpdatedUserID   uint64   `gorm:"column:updated_user_id;not null" json:"updated_user_id"`
	CreatedAtMs     uint64   `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAtMs     uint64   `gorm:"column:updated_at;not null" json:"updated_at"`
	IsDeleted       bool     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	DeletedByID     *uint64  `gorm:"column:deleted_by_id" json:"deleted_by_id,omitempty"`
}

func (HardwareTTSSettings) TableName() string { return "hardware_tts_settings" }