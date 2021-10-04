package client

import "github.com/streadway/amqp"

type ExchangeOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueOptions struct {
	Name        string
	Durable     bool
	AutoDelete  bool
	Exclusive   bool
	NoWait      bool
	BindOptions *QueueBindOptions
	Args        amqp.Table
}

type QueueBindOptions struct {
	ExchangeName string
	RoutingKey   string
	NoWait       bool
	Args         amqp.Table
}

type ConsumerOptions struct {
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type PublisherOptions struct {
	ExchangeName string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
	Persistent   bool
	Payload      *MessagePayload
}
