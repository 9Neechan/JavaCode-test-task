// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/9Neechan/JavaCode-test-task/rabbitmq (interfaces: AMQPChannel)

// Package mockrabbitmq is a generated GoMock package.
package mockrabbitmq

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	amqp "github.com/streadway/amqp"
)

// MockAMQPChannel is a mock of AMQPChannel interface.
type MockAMQPChannel struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPChannelMockRecorder
}

// MockAMQPChannelMockRecorder is the mock recorder for MockAMQPChannel.
type MockAMQPChannelMockRecorder struct {
	mock *MockAMQPChannel
}

// NewMockAMQPChannel creates a new mock instance.
func NewMockAMQPChannel(ctrl *gomock.Controller) *MockAMQPChannel {
	mock := &MockAMQPChannel{ctrl: ctrl}
	mock.recorder = &MockAMQPChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAMQPChannel) EXPECT() *MockAMQPChannelMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockAMQPChannel) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAMQPChannelMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAMQPChannel)(nil).Close))
}

// Consume mocks base method.
func (m *MockAMQPChannel) Consume(arg0, arg1 string, arg2, arg3, arg4, arg5 bool, arg6 amqp.Table) (<-chan amqp.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(<-chan amqp.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Consume indicates an expected call of Consume.
func (mr *MockAMQPChannelMockRecorder) Consume(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockAMQPChannel)(nil).Consume), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// Publish mocks base method.
func (m *MockAMQPChannel) Publish(arg0, arg1 string, arg2, arg3 bool, arg4 amqp.Publishing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockAMQPChannelMockRecorder) Publish(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockAMQPChannel)(nil).Publish), arg0, arg1, arg2, arg3, arg4)
}

// Qos mocks base method.
func (m *MockAMQPChannel) Qos(arg0, arg1 int, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Qos", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Qos indicates an expected call of Qos.
func (mr *MockAMQPChannelMockRecorder) Qos(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Qos", reflect.TypeOf((*MockAMQPChannel)(nil).Qos), arg0, arg1, arg2)
}

// QueueDeclare mocks base method.
func (m *MockAMQPChannel) QueueDeclare(arg0 string, arg1, arg2, arg3, arg4 bool, arg5 amqp.Table) (amqp.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueDeclare", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(amqp.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueDeclare indicates an expected call of QueueDeclare.
func (mr *MockAMQPChannelMockRecorder) QueueDeclare(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueDeclare", reflect.TypeOf((*MockAMQPChannel)(nil).QueueDeclare), arg0, arg1, arg2, arg3, arg4, arg5)
}
