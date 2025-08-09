package iot

// Envelope is the unified message wrapper for IoT <-> Coze communication via MQ.
// It intentionally uses simple json-friendly types and string IDs to be language-agnostic.
// Large binary payloads (audio) should be referenced via URI, not embedded.
//
// Versioning notes:
// - type: semantic category, e.g. "asr.text", "llm.request", "llm.result", "tts.request", "tts.result", "device.event"
// - streaming: if true, seq should be present and monotonic per message_id
// - payload: json object depending on type
//
// Required fields for routing and observability: message_id, device_id, space_id, app_id.
// Optional: user_id, session_id for detailed attribution.
//
// Time fields are milliseconds since epoch to align with existing schema.
//
// Warning: This contract is shared with an external service; keep it stable and backward-compatible.
type Envelope struct {
	MessageID string      `json:"message_id"`
	Type      string      `json:"type"`
	DeviceID  string      `json:"device_id"`
	SpaceID   string      `json:"space_id"`
	AppID     string      `json:"app_id"`
	UserID    string      `json:"user_id,omitempty"`
	SessionID string      `json:"session_id,omitempty"`
	TS        int64       `json:"ts"`
	Streaming bool        `json:"streaming,omitempty"`
	Seq       int64       `json:"seq,omitempty"`
	Payload   interface{} `json:"payload"`
}

// LLMTxtRequest is sent when IoT service asks Coze to generate agent reply.
// The Coze side should route to the bound app/agent and return LLMTxtResult on completion or stream chunks.
type LLMTxtRequest struct {
	Text       string            `json:"text"`
	Locale     string            `json:"locale,omitempty"`
	TraceCtx   map[string]string `json:"trace_ctx,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// LLMTxtResult is sent from Coze back to IoT service with the generated text.
type LLMTxtResult struct {
	Text       string            `json:"text"`
	Final      bool              `json:"final"`
	TraceCtx   map[string]string `json:"trace_ctx,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// TTSRequest requests the IoT service to synthesize audio based on current settings (device/app default).
// If Provider/Model/Voice are empty, service should resolve from device/app settings.
type TTSRequest struct {
	Text     string `json:"text"`
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
	Voice    string `json:"voice,omitempty"`
}

// TTSResult conveys synthesized audio by URL to avoid large payloads in MQ.
// The IoT service is expected to deliver audio to the device via its active transport.
type TTSResult struct {
	AudioURL  string            `json:"audio_url"`
	MimeType  string            `json:"mime_type,omitempty"`
	SampleRate int              `json:"sample_rate,omitempty"`
	DurationMs int64            `json:"duration_ms,omitempty"`
	TraceCtx   map[string]string `json:"trace_ctx,omitempty"`
}