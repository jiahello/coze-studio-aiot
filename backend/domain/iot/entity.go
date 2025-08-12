package iot

// HardwareDevice 是领域实体，表示硬件设备的业务信息
// 不包含任何 ORM 或存储层细节
// 时间字段统一以毫秒时间戳表达，符合现有上下游约定
// 注意：API 层需要自行做 DTO 映射以满足 snake_case 输出
 type HardwareDevice struct {
	ID              uint64
	DeviceID        string
	Name            string
	Description     string
	AppID           *uint64
	SpaceID         uint64
	CreatedUserID   uint64
	UpdatedUserID   uint64
	CreatedAtMs     uint64
	UpdatedAtMs     uint64
	FirmwareVersion *string
	MacAddress      *string
	Status          string
	LastPingAtMs    *uint64
	VerifyCode      *string
	FromEndUserID   *string
	FromAccountID   *string
}

// TTSVoice 是领域实体，表示 TTS 音色
 type TTSVoice struct {
	ID             uint64
	Provider       string
	Model          *string
	VoiceType      string
	Name           string
	VoiceCode      string
	Description    *string
	Gender         *string
	Language       *string
	Scenario       *string
	SoundQuality   *string
	SampleRate     *string
	TimestampSup   *string
	ErhuaSupport   *string
	SampleURL      *string
	CreatedAtMs    uint64
	UpdatedAtMs    uint64
	SourceFileID   *string
	SampleFileID   *string
	IsDeleted      bool
	LanguageName   *string
	VoiceID        *string
	ModelType      *string
	PreviewText    *string
	EmotionSupport *bool
	Emotions       *string
	FromAccountID  *string
	APIType        *string
	SpaceID        *uint64
}

// AppTTSSettings 是领域实体，表示应用级 TTS 设置（按 app 唯一）
 type AppTTSSettings struct {
	ID            uint64
	AppID         uint64
	Provider      string
	Model         string
	Voice         string
	VoiceRef      *uint64
	CreatedUserID uint64
	UpdatedUserID uint64
	CreatedAtMs   uint64
	UpdatedAtMs   uint64
}

// HardwareTTSSettings 是领域实体，表示设备级 TTS 设置（按 device 唯一）
 type HardwareTTSSettings struct {
	ID               uint64
	DeviceID         string
	HardwareDeviceID *uint64
	Provider         string
	Model            string
	Voice            string
	VoiceRef         *uint64
	CreatedUserID    uint64
	UpdatedUserID    uint64
	CreatedAtMs      uint64
	UpdatedAtMs      uint64
	IsDeleted        bool
}

// EffectiveTTS 是值对象，表示最终生效的 TTS 设置与来源
 type EffectiveTTS struct {
	Provider string
	Model    string
	Voice    string
	Source   string // device|app|default
}