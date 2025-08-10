package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/coze-dev/coze-studio/backend/api/model/conversation/run"
	convApp "github.com/coze-dev/coze-studio/backend/application/conversation"
	convMsg "github.com/coze-dev/coze-studio/backend/api/model/crossdomain/message"
	agentrunEntity "github.com/coze-dev/coze-studio/backend/domain/conversation/agentrun/entity"

	"github.com/coze-dev/coze-studio/backend/infra/impl/eventbus"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	contract "github.com/coze-dev/coze-studio/backend/infra/contract/eventbus"
	msg "github.com/coze-dev/coze-studio/backend/infra/contract/iot"
)

// Init wires IoT/Voice NSQ topics for minimal LLM <-> TTS loop.
func Init(ctx context.Context) (*Service, error) {
	nameServer := os.Getenv(consts.MQServer)
	if nameServer == "" {
		// Fallback to nsqd service name in docker network for local debug
		nameServer = "nsqd:4150"
	}

	// Register consumers
	svc := &Service{}
	if err := eventbus.DefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicDeviceInbound, consts.RMQConsumeGroupIoT, svc); err != nil {
		return nil, fmt.Errorf("register device inbound consumer failed: %w", err)
	}
	if err := eventbus.DefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicLLMResults, consts.RMQConsumeGroupLLM, svc); err != nil {
		return nil, fmt.Errorf("register llm results consumer failed: %w", err)
	}
	if err := eventbus.DefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicLLMTasks, consts.RMQConsumeGroupLLM, svc); err != nil {
		return nil, fmt.Errorf("register llm tasks consumer failed: %w", err)
	}
	if err := eventbus.DefaultSVC().RegisterConsumer(nameServer, consts.RMQTopicTTSResults, consts.RMQConsumeGroupTTS, svc); err != nil {
		return nil, fmt.Errorf("register tts results consumer failed: %w", err)
	}

	// Create producers
	var err error
	svc.deviceOutboundP, err = eventbus.NewProducer(nameServer, consts.RMQTopicDeviceOutbound, consts.RMQConsumeGroupIoT, 1)
	if err != nil {
		return nil, fmt.Errorf("init device outbound producer failed: %w", err)
	}
	svc.llmTasksP, err = eventbus.NewProducer(nameServer, consts.RMQTopicLLMTasks, consts.RMQConsumeGroupLLM, 1)
	if err != nil {
		return nil, fmt.Errorf("init llm tasks producer failed: %w", err)
	}
	svc.ttsTasksP, err = eventbus.NewProducer(nameServer, consts.RMQTopicTTSTasks, consts.RMQConsumeGroupTTS, 1)
	if err != nil {
		return nil, fmt.Errorf("init tts tasks producer failed: %w", err)
	}

	logs.Infof("IoT/Voice bus wired on %s", nameServer)
	return svc, nil
}

type Service struct {
	deviceOutboundP contract.Producer
	llmTasksP       contract.Producer
	ttsTasksP       contract.Producer
}

// HandleMessage implements ConsumerHandler for all three topics we subscribed.
func (s *Service) HandleMessage(ctx context.Context, m *contract.Message) error {
	var env msg.Envelope
	if err := json.Unmarshal(m.Body, &env); err != nil {
		logs.Errorf("[iot] invalid envelope: %v", err)
		return nil // skip bad message
	}
	switch env.Type {
	case "llm.request":
		// Forward to LLM tasks (Coze agent side)
		return s.forward(ctx, s.llmTasksP, &env)
	case "llm.task":
		// Consume llm.tasks: call agent and produce llm.results
		return s.handleLLMTask(ctx, &env)
	case "llm.result":
		// Forward to TTS tasks when final text arrives
		if r, ok := env.Payload.(map[string]any); ok {
			if final, _ := r["final"].(bool); final {
				// convert to tts.request
				text, _ := r["text"].(string)
				env.Type = "tts.request"
				env.Payload = msg.TTSRequest{Text: text}
				return s.forward(ctx, s.ttsTasksP, &env)
			}
		}
		return nil
	case "tts.result":
		// Deliver to device via device.outbound
		return s.forward(ctx, s.deviceOutboundP, &env)
	case "device.event", "asr.text":
		// Upstream into LLM
		env.Type = "llm.request"
		return s.forward(ctx, s.llmTasksP, &env)
	default:
		logs.Warnf("[iot] unknown type: %s", env.Type)
	}
	return nil
}

func (s *Service) forward(ctx context.Context, p contract.Producer, env *msg.Envelope) error {
	b, _ := json.Marshal(env)
	return p.Send(ctx, b)
}

func (s *Service) handleLLMTask(ctx context.Context, env *msg.Envelope) error {
	// map payload to ChatV3Request minimally
	req := &run.ChatV3Request{}
	// minimal fields: BotID(app/agent), Content, ContentType
	// For demo, parse from payload assuming {text: string, bot_id: number}
	m, _ := env.Payload.(map[string]any)
	text, _ := m["text"].(string)
	botID := int64(0)
	if v, ok := m["bot_id"].(float64); ok {
		botID = int64(v)
	}
	req.BotID = botID
	req.User = env.DeviceID
	msgItem := &run.EnterMessage{Role: "user", Content: text, ContentType: run.ContentTypeText}
	req.AdditionalMessages = []*run.EnterMessage{msgItem}
	// ensure conversation id kept across messages could be added later

	// build AgentRunMeta and call domain directly
	arm := &agentrunEntity.AgentRunMeta{
		ConversationID:   0,
		ConnectorID:      consts.APIConnectorID,
		SpaceID:          0,
		Scene:            0,
		SectionID:        0,
		Name:             "",
		UserID:           env.DeviceID,
		AgentID:          req.BotID,
		ContentType:      convMsg.ContentTypeText,
		Content:          []*convMsg.InputMetaData{{Type: convMsg.InputTypeText, Text: text}},
		PreRetrieveTools: nil,
		IsDraft:          false,
		CustomerConfig:   nil,
		DisplayContent:   text,
		CustomVariables:  nil,
		Version:          "",
		Ext:              nil,
	}
	if _, err := convApp.ConversationSVC.AgentRunDomainSVC.AgentRun(ctx, arm); err != nil {
		logs.Errorf("[iot] AgentRun failed: %v", err)
		return nil
	}
	// TODO: read stream and capture final text; temporarily return first chunk text
	finalText := text

	res := &msg.Envelope{
		MessageID: env.MessageID,
		Type:      "llm.result",
		DeviceID:  env.DeviceID,
		SpaceID:   env.SpaceID,
		AppID:     env.AppID,
		TS:        time.Now().UnixMilli(),
		Payload:   msg.LLMTxtResult{Text: finalText, Final: true},
	}
	return s.forward(ctx, s.llmTasksP, res)
}

// withSyntheticAPIAuth creates minimal context that satisfies openapi agent run requirements.
func withSyntheticAPIAuth(ctx context.Context) context.Context {
	// TODO: set API auth/session if required by ctxutil in deeper layers
	return ctx
}

// NOTE: streaming parse to extract final text will be implemented later