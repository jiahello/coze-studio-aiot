package coze

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	admin "github.com/coze-dev/coze-studio/backend/application/iotadmin"
	domain "github.com/coze-dev/coze-studio/backend/domain/iot"
)

type listDevicesReq struct {
	SpaceID uint64 `json:"space_id"`
	Page    int    `json:"page"`
	PageSize int   `json:"page_size"`
	Keyword string `json:"keyword"`
}

type upsertDeviceReq struct {
	ID        uint64 `json:"id"`
	SpaceID   uint64 `json:"space_id"`
	DeviceID  string `json:"device_id"`
	Name      string `json:"name"`
	AppID     *uint64 `json:"app_id"`
	Status    string `json:"status"`
	Desc      string `json:"description"`
}

type deviceDTO struct {
	ID              uint64  `json:"id"`
	DeviceID        string  `json:"device_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	AppID           *uint64 `json:"app_id,omitempty"`
	SpaceID         uint64  `json:"space_id"`
	CreatedUserID   uint64  `json:"created_user_id"`
	UpdatedUserID   uint64  `json:"updated_user_id"`
	CreatedAtMs     uint64  `json:"created_at"`
	UpdatedAtMs     uint64  `json:"updated_at"`
	FirmwareVersion *string `json:"firmware_version,omitempty"`
	MacAddress      *string `json:"mac_address,omitempty"`
	Status          string  `json:"status"`
	LastPingAtMs    *uint64 `json:"last_ping_at,omitempty"`
	VerifyCode      *string `json:"verify_code,omitempty"`
	FromEndUserID   *string `json:"from_end_user_id,omitempty"`
	FromAccountID   *string `json:"from_account_id,omitempty"`
}

func mapDeviceToDTO(d *domain.HardwareDevice) *deviceDTO {
	if d == nil { return nil }
	return &deviceDTO{
		ID: d.ID, DeviceID: d.DeviceID, Name: d.Name, Description: d.Description, AppID: d.AppID,
		SpaceID: d.SpaceID, CreatedUserID: d.CreatedUserID, UpdatedUserID: d.UpdatedUserID,
		CreatedAtMs: d.CreatedAtMs, UpdatedAtMs: d.UpdatedAtMs, FirmwareVersion: d.FirmwareVersion,
		MacAddress: d.MacAddress, Status: d.Status, LastPingAtMs: d.LastPingAtMs, VerifyCode: d.VerifyCode,
		FromEndUserID: d.FromEndUserID, FromAccountID: d.FromAccountID,
	}
}

// ListHardwareDevices
// @router /api/iot/devices/list [POST]
func ListHardwareDevices(ctx context.Context, c *app.RequestContext) {
	var req listDevicesReq
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	if req.SpaceID == 0 { httputil.BadRequest(c, "space_id required"); return }
	if req.Page <= 0 { req.Page = 1 }
	if req.PageSize <= 0 || req.PageSize > 200 { req.PageSize = 20 }
	list, total, err := admin.SVC.ListHardwareDevices(ctx, req.SpaceID, req.Page, req.PageSize, req.Keyword)
	if err != nil { internalServerErrorResponse(ctx, c, err); return }
	// 映射为 DTO（snake_case）
	out := make([]*deviceDTO, 0, len(list))
	for _, d := range list { out = append(out, mapDeviceToDTO(d)) }
	c.JSON(consts.StatusOK, map[string]any{"list": out, "total": total})
}

// UpsertHardwareDevice
// @router /api/iot/devices/upsert [POST]
func UpsertHardwareDevice(ctx context.Context, c *app.RequestContext) {
	var req upsertDeviceReq
	if err := c.BindAndValidate(&req); err != nil { httputil.BadRequest(c, err.Error()); return }
	if req.SpaceID == 0 || req.DeviceID == "" || req.Name == "" {
		httputil.BadRequest(c, "space_id, device_id, name are required"); return
	}
	dev := &domain.HardwareDevice{
		ID: req.ID, SpaceID: req.SpaceID, DeviceID: req.DeviceID, Name: req.Name, AppID: req.AppID, Status: req.Status, Description: req.Desc,
	}
	if err := admin.SVC.UpsertHardwareDevice(ctx, dev); err != nil { internalServerErrorResponse(ctx, c, err); return }
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

type listVoicesReq struct {
	SpaceID *uint64 `json:"space_id"`
	Provider string `json:"provider"`
	Language string `json:"language"`
	Gender   string `json:"gender"`
	Page    int    `json:"page"`
	PageSize int   `json:"page_size"`
}

type voiceDTO struct {
	ID       uint64  `json:"id"`
	Provider string  `json:"provider"`
	Model    *string `json:"model,omitempty"`
	VoiceType string `json:"voice_type"`
	Name     string  `json:"name"`
	VoiceCode string `json:"voice_code"`
	Language *string `json:"language,omitempty"`
	Gender   *string `json:"gender,omitempty"`
	SampleURL *string `json:"sample_url,omitempty"`
}

// ListTTSVoices
// @router /api/tts/voices/list [POST]
func ListTTSVoices(ctx context.Context, c *app.RequestContext) {
	var req listVoicesReq
	if err := c.BindAndValidate(&req); err != nil { httputil.BadRequest(c, err.Error()); return }
	if req.Page <= 0 { req.Page = 1 }
	if req.PageSize <= 0 || req.PageSize > 200 { req.PageSize = 20 }
	list, total, err := admin.SVC.ListTTSVoices(ctx, req.SpaceID, req.Provider, req.Language, req.Gender, req.Page, req.PageSize)
	if err != nil { internalServerErrorResponse(ctx, c, err); return }
	out := make([]*voiceDTO, 0, len(list))
	for _, v := range list {
		out = append(out, &voiceDTO{ID: v.ID, Provider: v.Provider, Model: v.Model, VoiceType: v.VoiceType, Name: v.Name, VoiceCode: v.VoiceCode, Language: v.Language, Gender: v.Gender, SampleURL: v.SampleURL})
	}
	c.JSON(consts.StatusOK, map[string]any{"list": out, "total": total})
}

type upsertAppTTSReq struct {
	AppID uint64 `json:"app_id"`
	Provider string `json:"provider"`
	Model string `json:"model"`
	Voice string `json:"voice"`
	VoiceRef *uint64 `json:"voice_ref"`
}

// UpsertAppTTS
// @router /api/apps/tts/set [POST]
func UpsertAppTTS(ctx context.Context, c *app.RequestContext) {
	var req upsertAppTTSReq
	if err := c.BindAndValidate(&req); err != nil { httputil.BadRequest(c, err.Error()); return }
	cfg := &domain.AppTTSSettings{AppID: req.AppID, Provider: req.Provider, Model: req.Model, Voice: req.Voice, VoiceRef: req.VoiceRef}
	if err := admin.SVC.UpsertAppTTS(ctx, cfg); err != nil { internalServerErrorResponse(ctx, c, err); return }
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

type upsertHardwareTTSReq struct {
	DeviceID string `json:"device_id"`
	Provider string `json:"provider"`
	Model string `json:"model"`
	Voice string `json:"voice"`
	VoiceRef *uint64 `json:"voice_ref"`
}

// UpsertHardwareTTS
// @router /api/iot/devices/tts/set [POST]
func UpsertHardwareTTS(ctx context.Context, c *app.RequestContext) {
	var req upsertHardwareTTSReq
	if err := c.BindAndValidate(&req); err != nil { httputil.BadRequest(c, err.Error()); return }
	cfg := &domain.HardwareTTSSettings{DeviceID: req.DeviceID, Provider: req.Provider, Model: req.Model, Voice: req.Voice, VoiceRef: req.VoiceRef}
	if err := admin.SVC.UpsertHardwareTTS(ctx, cfg); err != nil { internalServerErrorResponse(ctx, c, err); return }
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

// GetEffectiveDeviceTTS
// @router /api/iot/devices/tts/get [GET]
func GetEffectiveDeviceTTS(ctx context.Context, c *app.RequestContext) {
	deviceID := string(c.QueryArgs().Peek("device_id"))
	appIDStr := string(c.QueryArgs().Peek("app_id"))
	var appID *uint64
	if appIDStr != "" {
		var v uint64
		_, _ = fmt.Sscan(appIDStr, &v)
		appID = &v
	}
	if deviceID == "" { httputil.BadRequest(c, "device_id required"); return }
	res, err := admin.SVC.GetEffectiveDeviceTTS(ctx, deviceID, appID)
	if err != nil { internalServerErrorResponse(ctx, c, err); return }
	c.JSON(consts.StatusOK, res)
}

// TTSPreview
// @router /api/tts/preview [POST]
func TTSPreview(ctx context.Context, c *app.RequestContext) {
	type req struct{ Provider, Model, Voice, Text string; SpaceID *uint64 }
	var r req
	if err := c.BindAndValidate(&r); err != nil { httputil.BadRequest(c, err.Error()); return }
	// try get sample url first
	url, _ := admin.SVC.GetVoiceSampleURL(ctx, r.Provider, r.Voice, r.SpaceID)
	if url == "" {
		// fallback: return a signed placeholder or instruct device to synth via NSQ (not implemented here)
	}
	c.JSON(consts.StatusOK, map[string]any{"sample_url": url})
}