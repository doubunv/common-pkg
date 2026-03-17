package rabbitmq

import (
	"errors"
	"github.com/streadway/amqp"
	"time"
)

// Publish 单条消息
func (r *RabbitMQ) Publish(msg []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.channel == nil {
		return errors.New("channel not initialized")
	}

	return r.channel.Publish(
		r.Exchange,
		r.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
			Timestamp:   time.Now(),
		})
}

// PublishBatch 批量发送
func (r *RabbitMQ) PublishBatch(msgs [][]byte) error {
	for _, msg := range msgs {
		if err := r.Publish(msg); err != nil {
			return err
		}
	}
	return nil
}
