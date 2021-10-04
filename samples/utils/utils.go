package utils

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/ungame/go-rmq/client"
	"github.com/ungame/go-rmq/samples/models"
)

const (
	DefaultExchangeName = "Hello"
	DefaultQueueName    = "World"
)

var defaultMessageType = fmt.Sprintf("%T", models.User{})

func GetDefaultExchangeOptions() *client.ExchangeOptions {
	return &client.ExchangeOptions{
		Name:       DefaultExchangeName,
		Kind:       amqp.ExchangeDirect,
		Durable:    false,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

func GetDefaultQueueOptions() *client.QueueOptions {
	return &client.QueueOptions{
		Name:       DefaultQueueName,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		BindOptions: &client.QueueBindOptions{
			ExchangeName: DefaultExchangeName,
			RoutingKey:   "",
			NoWait:       false,
			Args:         nil,
		},
		Args: nil,
	}
}

func GetDefaultConsumerOptions() *client.ConsumerOptions {
	return &client.ConsumerOptions{
		QueueName: DefaultQueueName,
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
}

func GetDefaultPublisherOptions(data interface{}) *client.PublisherOptions {

	d, _ := json.Marshal(data)

	return &client.PublisherOptions{
		ExchangeName: DefaultExchangeName,
		RoutingKey:   "",
		Mandatory:    false,
		Immediate:    false,
		Persistent:   false,
		Payload: &client.MessagePayload{
			Type: defaultMessageType,
			Data: d,
		},
	}
}
