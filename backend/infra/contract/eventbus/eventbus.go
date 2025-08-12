/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eventbus

import "context"

//go:generate  mockgen -destination ../../../internal/mock/infra/contract/eventbus/eventbus_mock.go -package mock -source eventbus.go Factory
type Producer interface {
	Send(ctx context.Context, body []byte, opts ...SendOpt) error
	BatchSend(ctx context.Context, bodyArr [][]byte, opts ...SendOpt) error
}

var defaultSVC ConsumerService

func SetDefaultSVC(svc ConsumerService) {
	defaultSVC = svc
}

func GetDefaultSVC() ConsumerService {
	return defaultSVC
}

// ProducerFactory 提供统一的 Producer 构造，避免上层直接依赖具体实现
// nameServer/topic/group/retries 的含义与实现层保持一致
// 若未设置默认工厂，GetDefaultProducerFactory 将返回 nil
// 上层需在应用初始化阶段注入
//
// 该接口不破坏现有 ConsumerService 注入模式
type ConsumerService interface {
	RegisterConsumer(nameServer, topic, group string, consumerHandler ConsumerHandler, opts ...ConsumerOpt) error
}

type ConsumerHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

type Message struct {
	Topic string
	Group string
	Body  []byte
}

// ProducerFactory 定义
type ProducerFactory interface {
	NewProducer(nameServer, topic, group string, retries int) (Producer, error)
}

var defaultProducerFactory ProducerFactory

func SetDefaultProducerFactory(f ProducerFactory) {
	defaultProducerFactory = f
}

func GetDefaultProducerFactory() ProducerFactory {
	return defaultProducerFactory
}
