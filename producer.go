package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	tag        string
	done       chan error
}

func NewProducer(idx int, conn *amqp.Connection, exchange, exchangeType, key, ctag string, reliable bool) (*Producer, error) {
	p := &Producer{
		connection: conn,
		channel:    nil,
		tag:        ctag,
		done:       make(chan error),
	}

	var err error

	log.Printf("Getting Channel ")
	p.channel, err = p.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: ", err)
	}

	log.Printf("Declaring Exchange (%s)", exchange)
	if err := p.channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	// if reliable {
	// 	if err := p.channel.Confirm(false); err != nil {
	// 		return nil, fmt.Errorf("Channel could not be put into confirm mode: ", err)
	// 	}

	// 	ack, nack := p.channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	// 	// defer confirmOne(ack, nack)
	// }

	return p, nil
}

func (p *Producer) Publish(exchange, routingKey, body string) error {
	log.Printf("Publishing %s (%dB)", body, len(body))

	if err := p.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: ", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	log.Printf("waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}
