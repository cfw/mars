package amqp

import (
	log "github.com/sirupsen/logrus"
	amqplib "github.com/streadway/amqp"
	"sync"
	"time"
)

type Amqp struct {
	conn      *amqplib.Connection
	consumeCh *amqplib.Channel
	produceCh *amqplib.Channel
	mutex     sync.Mutex
}

func (a *Amqp) Consume(queue string, autoAck bool, callback Callback) {
	for {
		a.mutex.Lock()
		err := a.setupConsumeChannel()
		if err != nil {
			panic(err)
		}
		a.mutex.Unlock()
		messages, _ := a.consumeCh.Consume(
			queue,
			"",
			autoAck,
			false,
			false,
			false,
			nil,
		)
		for d := range messages {
			go callback(d)
		}
	}

}

func (a *Amqp) Publish(exchange, routingKey string, publishing amqplib.Publishing) {
	_ = a.setupProduce()
	var err error
	err = a.produceCh.Publish(exchange, routingKey, false, false, publishing)
	if err != nil {
		log.Error(err)
	}
}
func (a *Amqp) Close() {
	log.Info("Stopping RabbitMQ")
	if a.produceCh != nil {
		_ = a.produceCh.Close()
	}
	if a.consumeCh != nil {
		_ = a.consumeCh.Close()
	}
	if a.conn != nil || a.conn.IsClosed() {
		_ = a.conn.Close()
	}
	log.Info("Stopped RabbitMQ")
}

func (a *Amqp) setupProduce() error {
	var err error
	for {
		if a.conn.IsClosed() {
			continue
		}
		if a.produceCh == nil {
			a.produceCh, err = a.conn.Channel()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (a *Amqp) setupConsumeChannel() error {
	var err error
	for {
		if a.conn.IsClosed() {
			continue
		}
		if a.consumeCh == nil {
			a.consumeCh, err = a.conn.Channel()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (a *Amqp) connect(url string) error {
	var err error
	if a.conn == nil || a.conn.IsClosed() {
		a.consumeCh = nil
		a.conn, err = amqplib.Dial(url)
		log.Info("Connected to RabbitMQ")
		if err != nil {
			return err
		}
		go func() {
			for {
				reason, ok := <-a.conn.NotifyClose(make(chan *amqplib.Error))
				if !ok {
					log.Info("Connection closed")
					break
				}
				log.Info("Connection closed, reason: ", reason)
				if a.consumeCh != nil {
					_ = a.consumeCh.Close()
					a.consumeCh = nil
				}
				a.consumeCh = nil
				if a.produceCh != nil {
					_ = a.produceCh.Close()
					a.produceCh = nil
				}

				for {
					time.Sleep(time.Second)
					if conn, err := amqplib.Dial(url); err == nil {
						a.conn = conn
						log.Info("Reconnect success")
						break
					}
					log.Info("Reconnect failed, err: ", err)
				}
			}
		}()
	}
	return err
}

type Callback func(delivery amqplib.Delivery)

func NewAmqp(c *Config) *Amqp {
	instance := &Amqp{}
	if err := instance.connect(c.Url()); err != nil {
		panic(err)
	}
	return instance
}
