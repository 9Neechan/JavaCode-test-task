package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	//db "github.com/9Neechan/JavaCode-test-task/db/sqlc"

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
		return nil, fmt.Errorf("Ошибка установки QoS: %s", err)
		//log.Fatalf("Ошибка установки QoS: %s", err)
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
	// Проверка на nil-обработчик
	if handler == nil {
		return fmt.Errorf("обработчик не может быть nil")
	}

	// Проверка на пустое имя очереди
	if queueName == "" {
		return fmt.Errorf("ошибка объявления очереди: имя очереди не может быть пустым")
	}

	// Объявляем очередь
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("ошибка объявления очереди: %w", err)
	}

	// Получаем сообщения
	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
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
