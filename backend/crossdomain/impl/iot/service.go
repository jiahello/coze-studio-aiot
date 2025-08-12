package iot

import (
	"context"

	domainSvc "github.com/coze-dev/coze-studio/backend/domain/iot/service"
	contract "github.com/coze-dev/coze-studio/backend/crossdomain/contract/iot"
)

type Adapter struct {
	domain *domainSvc.Service
}

func NewAdapter(domain *domainSvc.Service) *Adapter { return &Adapter{domain: domain} }

func (a *Adapter) GetEffectiveTTS(ctx context.Context, deviceID string, appID *uint64) (*contract.EffectiveTTS, error) {
	res, err := a.domain.GetEffectiveTTS(ctx, deviceID, appID)
	if err != nil { return nil, err }
	return &contract.EffectiveTTS{Provider: res.Provider, Model: res.Model, Voice: res.Voice, Source: res.Source}, nil
}