package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connStr    string
	conn       *amqp.Connection
	channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string

	mu sync.Mutex
}

// NewRabbitMQ 初始化 RabbitMQ 实例
func NewRabbitMQ(connStr, queueName, exchange, routingKey string) (*RabbitMQ, error) {
	r := &RabbitMQ{
		connStr:    connStr,
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingKey,
	}

	if err := r.connect(); err != nil {
		return nil, err
	}

	go r.watchConnection() // 后台自动重连
	return r, nil
}

// connect 建立连接和通道
func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.connStr)
	if err != nil {
		return fmt.Errorf("failed to connect RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// 队列声明
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	// 交换机声明
	if r.Exchange != "" {
		if err = r.channel.ExchangeDeclare(
			r.Exchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("exchange declare failed: %w", err)
		}
	}

	return nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

// watchConnection 自动重连
func (r *RabbitMQ) watchConnection() {
	for {
		if r.conn == nil || r.conn.IsClosed() {
			log.Println("RabbitMQ connection closed, reconnecting...")
			for {
				if err := r.connect(); err == nil {
					log.Println("RabbitMQ reconnected")
					break
				}
				log.Println("Reconnect failed, retry in 5s")
				time.Sleep(5 * time.Second)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
