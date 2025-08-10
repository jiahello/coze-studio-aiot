package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	admin "github.com/coze-dev/coze-studio/backend/application/iotadmin"
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

// ListHardwareDevices
// @router /api/iot/devices/list [POST]
func ListHardwareDevices(ctx context.Context, c *app.RequestContext) {
	var req listDevicesReq
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	list, total, err := admin.SVC.ListHardwareDevices(ctx, req.SpaceID, req.Page, req.PageSize, req.Keyword)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"list": list, "total": total})
}

// UpsertHardwareDevice
// @router /api/iot/devices/upsert [POST]
func UpsertHardwareDevice(ctx context.Context, c *app.RequestContext) {
	var req upsertDeviceReq
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	dev := &admin.HardwareDevice{
		ID: req.ID,
		SpaceID: req.SpaceID,
		DeviceID: req.DeviceID,
		Name: req.Name,
		AppID: req.AppID,
		Status: req.Status,
		Description: req.Desc,
	}
	if err := admin.SVC.UpsertHardwareDevice(ctx, dev); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

type listVoicesReq struct {
	SpaceID *uint64 `json:"space_id"`
	Provider string `json:"provider"`
	Page    int    `json:"page"`
	PageSize int   `json:"page_size"`
}

// ListTTSVoices
// @router /api/tts/voices/list [POST]
func ListTTSVoices(ctx context.Context, c *app.RequestContext) {
	var req listVoicesReq
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	list, total, err := admin.SVC.ListTTSVoices(ctx, req.SpaceID, req.Provider, req.Page, req.PageSize)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"list": list, "total": total})
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
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	cfg := &admin.AppTTSSettings{AppID: req.AppID, Provider: req.Provider, Model: req.Model, Voice: req.Voice, VoiceRef: req.VoiceRef}
	if err := admin.SVC.UpsertAppTTS(ctx, cfg); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
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
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}
	cfg := &admin.HardwareTTSSettings{DeviceID: req.DeviceID, Provider: req.Provider, Model: req.Model, Voice: req.Voice, VoiceRef: req.VoiceRef}
	if err := admin.SVC.UpsertHardwareTTS(ctx, cfg); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}