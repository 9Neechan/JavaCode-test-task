package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	store   db.Store
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

func (r *RabbitMQ) PublishWalletUpdate(update db.TransferTxParams) error {
	body, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) processWalletUpdates() {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to consume messages: %v", err)
	}

	for msg := range msgs {
		var update db.TransferTxParams
		if err := json.Unmarshal(msg.Body, &update); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		_, err := r.store.TransferTx(nil, update)
		if err != nil {
			log.Printf("failed to update wallet balance: %v", err)
		}
	}
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
