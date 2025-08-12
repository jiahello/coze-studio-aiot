# IoT 模块 DDD 对齐说明

## 分层与目录
- 领域层：`backend/domain/iot/`
  - 实体：`entity.go`（无 ORM 标签）
  - 仓储接口：`repository/repository.go`
  - 领域服务：`service/service.go`
  - DAL（仅内部）：`internal/dal/`（GORM PO 与 DAO 实现）
- 应用层：`backend/application/iotadmin/`
  - 组合仓储与领域服务，编排用例、补齐时间戳
- API 层：`backend/api/handler/coze/iot_service.go`
  - 入参校验与 DTO 映射，不承载业务
- 基础设施：通过契约注入（事件总线在 `infra/contract/eventbus`）
- 跨域能力：`crossdomain/contract/iot` + `crossdomain/impl/iot`

## 关键规则
- 仓储只做数据访问；“生效 TTS 规则（device > app > default）”位于领域服务
- API/应用层禁止越层访问 DAL；所有访问通过仓储接口
- 事件总线/外部系统通过契约获取，避免应用层依赖实现

## 跨域调用与初始化
- 初始化（在 `backend/application/iotadmin/init.go` 完成）：
```go
import (
    crossIOT "github.com/coze-dev/coze-studio/backend/crossdomain/contract/iot"
    implIOT "github.com/coze-dev/coze-studio/backend/crossdomain/impl/iot"
    domainSvc "github.com/coze-dev/coze-studio/backend/domain/iot/service"
)

// 装配仓储与领域服务
svc := domainSvc.NewService(deviceRepo, voiceRepo, settingsRepo)
// 注册跨域默认服务
crossIOT.SetDefaultSVC(implIOT.NewAdapter(svc))
```
- 业务侧调用（如 `backend/application/iot/init.go`）：
```go
import crossIOT "github.com/coze-dev/coze-studio/backend/crossdomain/contract/iot"

svc := crossIOT.GetDefaultSVC()
if svc != nil {
    eff, _ := svc.GetEffectiveTTS(ctx, deviceID, appIDPtr)
    // 使用 eff.Provider / eff.Model / eff.Voice
}
```

## 扩展方式
- 新增规则：在 `service/service.go` 中扩展方法，必要时补充仓储查询接口
- 新增存储字段：同步更新 `internal/dal/model.go` 与实体字段，映射在 DAO 中补齐
- 跨域调用：在 `crossdomain/contract/iot` 增口，在 `crossdomain/impl/iot` 适配领域服务

## 校验
- `go build ./...` 通过
- 单测：`backend/domain/iot/service/service_test.go` 覆盖 TTS 生效优先级逻辑