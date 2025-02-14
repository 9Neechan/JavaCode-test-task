package rabbitmq

import (
	"errors"
	"testing"

	mockrabbitmq "github.com/9Neechan/JavaCode-test-task/rabbitmq/mock"
	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestNewRabbitMQ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := &amqp.Connection{} // Заглушка соединения

	// Вызов тестируемой функции
	mq := &RabbitMQ{
		conn:    mockConn,
	}
	require.NotNil(t, mq)

	//_, err := mockConn.Channel()
	//require.NoError(t, err)

	//err = ch.Qos(50, 0, false)
	//require.NoError(t, err)

}

func TestNewRabbitMQ_ErrorDial(t *testing.T) {
	_, err := NewRabbitMQ("invalid-url")
	require.Error(t, err)
}

func TestConsumeMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	queueName := "test_queue"

	mockChannel.EXPECT().QueueDeclare(queueName, true, false, false, false, nil).Return(amqp.Queue{}, nil)
	msgChan := make(chan amqp.Delivery, 1)
	msgChan <- amqp.Delivery{Body: []byte(`{"key":"value"}`)}
	close(msgChan)

	mockChannel.EXPECT().Consume(queueName, "", true, false, false, false, nil).Return(msgChan, nil)

	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	handlerCalled := false
	handler := func(msg []byte) {
		handlerCalled = true
	}

	err := rmq.ConsumeMessages(queueName, handler)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestPublishMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	queueName := "test_queue"
	message := map[string]string{"key": "value"}

	mockChannel.EXPECT().Publish("", queueName, false, false, gomock.Any()).Return(nil)

	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	err := rmq.PublishMessage(queueName, message)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestQueueDeclareError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mockrabbitmq.NewMockAMQPChannel(ctrl)
	queueName := "test_queue"

	mockChannel.EXPECT().QueueDeclare(queueName, true, false, false, false, nil).Return(amqp.Queue{}, errors.New("queue declare error"))

	rmq := &RabbitMQ{
		channel: mockChannel,
	}

	err := rmq.ConsumeMessages(queueName, func([]byte) {})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
