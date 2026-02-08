package rabbitmq

import "errors"

var (
	// ErrMessageNotConfirmed 消息未被确认
	ErrMessageNotConfirmed = errors.New("message was not confirmed by broker")
	// ErrConfirmTimeout 确认超时
	ErrConfirmTimeout = errors.New("message confirmation timeout")
	// ErrConsumerStopped 消费者已停止
	ErrConsumerStopped = errors.New("consumer stopped")
)
