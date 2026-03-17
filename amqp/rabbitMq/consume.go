package rabbitmq

import "errors"

// Consume 消费消息
func (r *RabbitMQ) Consume(autoAck bool, handler func([]byte) error, workerNum int) error {
	if r.channel == nil {
		return errors.New("channel not initialized")
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		autoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerNum; i++ {
		go func() {
			for d := range msgs {
				// 幂等处理：handler 返回 nil 才 ack
				if err := handler(d.Body); err == nil {
					d.Ack(false)
				} else {
					d.Nack(false, true) // requeue
				}
			}
		}()
	}
	return nil
}
