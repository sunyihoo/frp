package transport

import (
	"context"
	"github.com/sunyihoo/frp/pkg/msg"
	"sync"
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
	return &
}

type transportImpl struct {
	sendCh chan msg.Message

	// 第一个键是消息类型，第二个键是通道键。
	// Dispatch 会按消息类型和通道键将消息调度到相关通道。
	//
	registry map[string]map[string]chan msg.Message
	mu sync.RWMutex
}