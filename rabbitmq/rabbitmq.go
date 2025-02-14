package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	//db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	mockrabbitmq "github.com/9Neechan/JavaCode-test-task/rabbitmq/mock"
	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
)

type AMQPChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Qos(prefetchCount, prefetchSize int, global bool) error
	Close() error
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel AMQPChannel // *amqp.Channel
	//queue   amqp.Queue
	//store   db.Store
}

func NewMockRabbitMQ(ctrl *gomock.Controller) *RabbitMQ {
	return &RabbitMQ{
		channel: mockrabbitmq.NewMockAMQPChannel(ctrl),
	}
}

func NewRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	//!!!!!!!!!!!!!
	err = ch.Qos(50, 0, false) // Обработчик берет до 50 сообщений сразу
	if err != nil {
		log.Fatalf("Ошибка установки QoS: %s", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

// ConsumeMessages слушает очередь и обрабатывает сообщения
func (r *RabbitMQ) ConsumeMessages(queueName string, handler func([]byte)) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ошибка объявления очереди: %w", err)
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("ошибка получения сообщений: %w", err)
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()

	log.Printf("✅ Начат прием сообщений из очереди: %s", queueName)
	return nil
}

// ✅ Реализуем метод PublishMessage
func (r *RabbitMQ) PublishMessage(queueName string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("✅ [RabbitMQ] Сообщение отправлено в очередь %s", queueName)
	return nil
}
