package api

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type RmqMessage struct {
	Event string `json:"event"`
}

// RmqConnection - main object with connection data
type RmqConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	stop chan struct{}
	wg   *sync.WaitGroup
}

// Close connection to rabbit
func (c *RmqConnection) Close() {
	close(c.stop)
	c.wg.Wait()
	c.ch.Close()
	c.conn.Close()
}

// RabbitMQConnect - connect and create data channel
func RabbitMQConnect(addr string) (*RmqConnection, error) {
	var c RmqConnection
	var err error
	if c.conn, err = amqp.Dial(addr); err != nil {
		return nil, err
	}
	if c.ch, err = c.conn.Channel(); err != nil {
		return nil, err

	}
	c.stop = make(chan struct{})
	c.wg = &sync.WaitGroup{}
	return &c, nil
}

// DeclareQueue - create queue in rabbit with default parameters
func (c *RmqConnection) DeclareQueue(name string) error {
	_, err := c.ch.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// SendMessage - marshaling object into json and send into queue
func (c *RmqConnection) SendMessage(name string, message []byte) error {
	err := c.ch.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         message,
		})
	return err
}

// Subscribe to message queue in rabbit
func (c *RmqConnection) Subscribe(name string) (<-chan []byte, error) {
	rmqchan, err := c.ch.Consume(
		name,  // queue name
		"",    // consumer
		false, // auto ask
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil)   // arguments
	if err != nil {
		return nil, err
	}
	c.wg.Add(1)
	datach := make(chan []byte, 1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.stop:
				close(datach)
				return
			case msg := <-rmqchan:
				datach <- msg.Body
				if err := msg.Ack(false); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
	return datach, nil
}

// IsNoQueueError - check error for queue is not declared in rabbitmq
func IsNoQueueError(err error) bool {
	if amqperr, ok := err.(*amqp.Error); ok && amqperr.Code == amqp.NotFound {
		return true
	}
	return false
}
