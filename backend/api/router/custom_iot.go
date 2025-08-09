package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	coze "github.com/coze-dev/coze-studio/backend/api/handler/coze"
)

func init() {
	registerers = append(registerers, registerIoTRoutes)
}

func registerIoTRoutes(r *server.Hertz) {
	api := r.Group("/api")
	iot := api.Group("/iot")
	devices := iot.Group("/devices")
	devices.POST("/list", coze.ListHardwareDevices)
	devices.POST("/upsert", coze.UpsertHardwareDevice)
	devices.POST("/tts/set", coze.UpsertHardwareTTS)

	tts := api.Group("/tts")
	tts.POST("/voices/list", coze.ListTTSVoices)

	apps := api.Group("/apps")
	apps.POST("/tts/set", coze.UpsertAppTTS)
}