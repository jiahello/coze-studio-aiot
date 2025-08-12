package iot

import "context"

// EffectiveTTS 为跨域可见的值对象
 type EffectiveTTS struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Voice    string `json:"voice"`
	Source   string `json:"source"` // device|app|default
}

// Service 定义 IoT 跨域能力入口
 type Service interface {
	GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*EffectiveTTS, error)
}

var defaultSVC Service

func SetDefaultSVC(svc Service) { defaultSVC = svc }
func GetDefaultSVC() Service    { return defaultSVC }