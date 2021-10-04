package client

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Client interface {
	ExchangeDeclare(opts *ExchangeOptions) error
	QueueDeclare(opts *QueueOptions) error
	Listen(opts *ConsumerOptions) error
	HandleFunc(messageType interface{}, handle func(data []byte) error)
	Send(opts *PublisherOptions) error
	Close()
}

type _client struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	exchanges map[string]*ExchangeOptions
	queues    map[string]*amqp.Queue
	handlers  map[string]func(data []byte) error
}

func NewClient(cfg Config) (Client, error) {
	var c _client
	var err error

	c.conn, err = amqp.Dial(cfg.URL())
	if err != nil {
		return nil, err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	c.exchanges = make(map[string]*ExchangeOptions)
	c.queues = make(map[string]*amqp.Queue)
	c.handlers = make(map[string]func(data []byte) error)

	return &c, nil
}

func (c *_client) Close() {
	HandleClose(c.channel)
	HandleClose(c.conn)
}

func (c *_client) ExchangeDeclare(opts *ExchangeOptions) error {
	err := c.channel.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}
	c.exchanges[opts.Name] = opts
	return nil
}

func (c *_client) QueueDeclare(opts *QueueOptions) error {
	queue, err := c.channel.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}

	if opts.BindOptions != nil {

		bindOpts := opts.BindOptions

		err = c.channel.QueueBind(
			queue.Name,
			bindOpts.RoutingKey,
			bindOpts.ExchangeName,
			bindOpts.NoWait,
			bindOpts.Args,
		)
		if err != nil {
			return err
		}
	}

	c.queues[opts.Name] = &queue

	return nil
}

func (c *_client) HandleFunc(messageType interface{}, handler func(data []byte) error) {
	c.handlers[fmt.Sprintf("%T", messageType)] = handler
}

func (c *_client) Listen(opts *ConsumerOptions) error {

	messages, err := c.channel.Consume(
		opts.QueueName,
		opts.Consumer,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}

	for message := range messages {

		var payload MessagePayload
		err := json.Unmarshal(message.Body, &payload)
		if err != nil {
			return err
		}

		if !opts.AutoAck {
			err := message.Ack(false)
			if err != nil {
				return err
			}
		}

		handler, ok := c.handlers[payload.Type]
		if !ok {
			log.Printf("received message without registered handler: MessageType=%s", payload.Type)
			continue
		}

		err = handler(payload.Data)
		if err != nil {
			log.Printf("error on handler message: %T\n", payload.Data)
		}
	}

	return nil
}

func (c *_client) Send(opts *PublisherOptions) error {

	payload, err := json.Marshal(opts.Payload)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}

	if opts.Persistent {
		message.DeliveryMode = amqp.Persistent
	}

	return c.channel.Publish(
		opts.ExchangeName,
		opts.RoutingKey,
		opts.Mandatory,
		opts.Immediate,
		message,
	)
}
