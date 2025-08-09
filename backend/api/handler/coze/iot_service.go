package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	admin "github.com/coze-dev/coze-studio/backend/application/iotadmin"
)

// ListHardwareDevices
// @router /api/iot/devices/list [POST]
func ListHardwareDevices(ctx context.Context, c *app.RequestContext) {
	// TODO: replace with real filters and paging
	_ = ctx
	c.JSON(consts.StatusOK, map[string]any{
		"list": []any{},
		"total": 0,
	})
}

// UpsertHardwareDevice
// @router /api/iot/devices/upsert [POST]
func UpsertHardwareDevice(ctx context.Context, c *app.RequestContext) {
	// TODO: bind request and persist via admin.SVC.DB
	if admin.SVC == nil || admin.SVC.DB == nil {
		httputil.InternalError(ctx, c, nil)
		return
	}
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

// ListTTSVoices
// @router /api/tts/voices/list [POST]
func ListTTSVoices(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{
		"list": []any{},
		"total": 0,
	})
}

// UpsertAppTTS
// @router /api/apps/tts/set [POST]
func UpsertAppTTS(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}

// UpsertHardwareTTS
// @router /api/iot/devices/tts/set [POST]
func UpsertHardwareTTS(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, map[string]any{"ok": true})
}