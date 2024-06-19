// Copyright 2023 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"context"
	"reflect"
	"sync"

	"github.com/sunyihoo/frp/pkg/msg"

	"github.com/fatedier/golib/errors"
)

type MessageTransporter interface {
	Send(msg.Message) error
	// Recv(ctx context.Context, laneKey string, msgType string) (Message, error)
	// Do 将首先发送 msg，然后使用相同的 laneKey 和指定的 msgType 接收 msg。
	Do(ctx context.Context, req msg.Message, laneKey, recvMsgType string) (msg.Message, error)
	// Dispatch 会按消息类型和 laneKey 将消息调度到 Do 函数中注册的相关 channel
	Dispatch(m msg.Message, laneKey string) bool
	// DispatchWithType 这个函数和 Dispatch 函数类似，但它只处理特定类型的消息
	DispatchWithType(m msg.Message, msgType, laneKey string) bool
}

func NewMessageTransporter(sendCh chan msg.Message) MessageTransporter {
	return &transportImpl{
		sendCh:   sendCh,
		registry: make(map[string]map[string]chan msg.Message),
	}
}

type transportImpl struct {
	sendCh chan msg.Message

	// 第一个键是消息类型，第二个键是通道键。
	// Dispatch 会按消息类型和通道键将消息调度到相关通道。
	//
	registry map[string]map[string]chan msg.Message
	mu       sync.RWMutex
}

func (impl *transportImpl) Send(m msg.Message) error {
	return errors.PanicToError(func() {
		impl.sendCh <- m
	})
}

func (impl *transportImpl) Do(ctx context.Context, req msg.Message, laneKey, recvMsgType string) (msg.Message, error) {
	ch := make(chan msg.Message, 1)
	defer close(ch)
	unregisterFn := impl.registerMsgChan(ch, laneKey, recvMsgType)
	defer unregisterFn()

	if err := impl.Send(req); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-ch:
		return resp, nil
	}
}

func (impl *transportImpl) DispatchWithType(m msg.Message, msgType string, laneKey string) bool {
	var ch chan msg.Message
	impl.mu.RLock()
	byLaneKey, ok := impl.registry[msgType]
	if ok {
		ch = byLaneKey[msgType]
	}
	impl.mu.RUnlock()

	if ch == nil {
		return false
	}

	if err := errors.PanicToError(func() {
		ch <- m
	}); err != nil {
		return false
	}
	return true
}

func (impl *transportImpl) Dispatch(m msg.Message, laneKey string) bool {
	msgType := reflect.TypeOf(m).Elem().Name()
	return impl.DispatchWithType(m, msgType, laneKey)
}

func (impl *transportImpl) registerMsgChan(recvCh chan msg.Message, laneKey string, msgType string) (unregister func()) {
	impl.mu.Lock()
	byLaneKey, ok := impl.registry[laneKey]
	if !ok {
		byLaneKey = make(map[string]chan msg.Message)
		impl.registry[msgType] = byLaneKey
	}
	byLaneKey[laneKey] = recvCh
	impl.mu.Unlock()

	unregister = func() {
		impl.mu.Lock()
		delete(byLaneKey, laneKey)
		impl.mu.Unlock()
	}
	return
}
