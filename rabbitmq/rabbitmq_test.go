package rabbitmq

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	mockrabbitmq "github.com/9Neechan/JavaCode-test-task/rabbitmq/mock"
	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

const testAMQPURL = "amqp://guest:guest@localhost:5672/"
const testQueueName = "test_queue"

// * Не забудьте make rabbitmq

func TestNewRabbitMQ_Success(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)
	require.NotNil(t, rmq)
	require.NotNil(t, rmq.conn)
	require.NotNil(t, rmq.channel)

	// Закрываем соединение после теста
	rmq.Close()
}

func TestNewRabbitMQ_InvalidURL(t *testing.T) {
	_, err := NewRabbitMQ("amqp://invalid:invalid@localhost:5672/")
	require.Error(t, err, "должна быть ошибка при неправильном URL")
}

func TestNewRabbitMQ_QosError(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)

	// Закрываем канал перед установкой QoS, чтобы вызвать ошибку
	rmq.channel.Close()

	err = rmq.channel.Qos(50, 0, false)
	require.Error(t, err, "должна быть ошибка при установке QoS на закрытом канале")

	rmq.Close()
}

func TestClose(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)

	// Проверяем, что соединение открыто
	require.NotNil(t, rmq.conn)
	require.NotNil(t, rmq.channel)

	// Закрываем соединение
	rmq.Close()

	// Проверяем, что закрытие не вызывает панику (в реальном коде можно проверять закрыто ли соединение)
	require.NotNil(t, rmq)
}

func TestConsumeMessages_NilHandler(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)
	defer rmq.Close()

	err = rmq.ConsumeMessages(testQueueName, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "обработчик не может быть nil")
}

func TestConsumeMessages_EmptyQueueName(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)
	defer rmq.Close()

	err = rmq.ConsumeMessages("", func([]byte) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "имя очереди не может быть пустым")
}

func TestConsumeMessages_Success(t *testing.T) {
	rmq, err := NewRabbitMQ(testAMQPURL)
	require.NoError(t, err)
	defer rmq.Close()

	messageReceived := make(chan []byte, 1)

	handler := func(body []byte) {
		messageReceived <- body
	}

	err = rmq.ConsumeMessages(testQueueName, handler)
	require.NoError(t, err)

	// Публикуем тестовое сообщение
	testMessage := "Hello, RabbitMQ!"
	err = rmq.channel.Publish(
		"",            // exchange
		testQueueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(testMessage),
		},
	)
	require.NoError(t, err)

	select {
	case received := <-messageReceived:
		require.Equal(t, testMessage, string(received))
	case <-time.After(2 * time.Second):
		t.Fatal("Не получено сообщение за 2 секунды")
	}
}

func TestConsumeMessages_ConsumeError_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)

	// Ожидаем вызов QueueDeclare, который должен пройти успешно
	mockChannel.EXPECT().
		QueueDeclare(gomock.Any(), true, false, false, false, gomock.Nil()).
		Return(amqp.Queue{Name: "test_queue"}, nil)

	// Ожидаем вызов Consume, который должен вернуть ошибку
	mockChannel.EXPECT().
		Consume(gomock.Any(), gomock.Any(), true, false, false, false, gomock.Nil()).
		Return(nil, errors.New("ошибка получения сообщений"))

	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	err := rmq.ConsumeMessages("test_queue", func([]byte) {})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ошибка получения сообщений")
}

func TestPublishMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	message := map[string]string{"key": "value"}
	body, _ := json.Marshal(message)

	// Ожидаем вызов Publish без ошибок
	mockChannel.EXPECT().
		Publish("", "test_queue", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}).
		Return(nil)

	err := rmq.PublishMessage("test_queue", message)
	require.NoError(t, err)
}

func TestPublishMessage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	message := map[string]string{"key": "value"}
	body, _ := json.Marshal(message)

	// Ожидаем вызов Publish, который вернёт ошибку
	mockChannel.EXPECT().
		Publish("", "test_queue", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}).
		Return(errors.New("ошибка публикации"))

	err := rmq.PublishMessage("test_queue", message)
	require.Error(t, err)
	require.Contains(t, err.Error(), "ошибка публикации")
}

func TestPublishMessage_ErrorOnMarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	// Создаем объект, который вызовет ошибку при сериализации
	message := make(chan int) // Каналы нельзя сериализовать в JSON

	err := rmq.PublishMessage("test_queue", message)
	require.Error(t, err)
	require.Contains(t, err.Error(), "json: unsupported type")
}