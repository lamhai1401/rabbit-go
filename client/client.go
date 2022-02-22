package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/lamhai1401/gologs/logs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	queueName string
	cfg       *tls.Config
	conn      *amqp.Connection
	url       string
	mutex     sync.RWMutex
	ctxDone   context.Context    // to signal reconnect
	ctxCancel context.CancelFunc // to cancel ctx
	callBack  func([]byte)
	isClosed  bool // check peer close or not
}

func NewRabbitClient(url string, queueName string, cfg *tls.Config, callBack func([]byte)) *RabbitClient {
	// ctx, cancel := context.WithCancel(context.Background())
	return &RabbitClient{
		queueName: queueName,
		cfg:       cfg,
		url:       url,
		// ctxDone:   ctx,
		// ctxCancel: cancel,
		callBack: callBack,
	}
}

// dial dial to rabbit
func (b *RabbitClient) dial(cfg *tls.Config) (*amqp.Connection, error) {
	// TODO add context to check timeout connection
	var conn *amqp.Connection
	var err error

	errChann := make(chan error)
	resultChan := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*time.Second))

	go func() {
		if cfg != nil {
			conn, err = amqp.DialTLS(b.url, cfg)
		} else {
			conn, err = amqp.Dial(b.url)
		}

		if err != nil {
			errChann <- err
		} else {
			resultChan <- 1
		}
		cancel()
	}()

	select {
	case err := <-errChann:
		return nil, err
	case <-resultChan:
		return conn, err
	case <-ctx.Done():
		return nil, fmt.Errorf(ctx.Err().Error())
	}
}

func (b *RabbitClient) Close() {
	if !b.checkClose() {
		b.setClose(true)
		cancel := b.getCtxCancel()
		if cancel != nil {
			cancel()
		}
	}
}

// Connect linter
func (b *RabbitClient) connect() (*amqp.Connection, error) {
	limit := 100
	count := 0

	var conn *amqp.Connection
	var err error

	for {
		if count == limit {
			return nil, fmt.Errorf("cannot connect to %s after 100 times try", b.url)
		}
		count++
		logs.Warn(fmt.Sprintf("Try to connect (%s) at %d times", b.url, count))
		conn, err = b.dial(b.cfg)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		b.setConn(conn)
		break
	}
	logs.Info("Connecting to ", b.url)
	return conn, nil
}

func (b *RabbitClient) reconnect() error {
	// close old connection
	if conn := b.getConn(); conn != nil {
		b.setConn(nil)
		conn.Close()
	}

	// signal ctx
	if cancel := b.getCtxCancel(); cancel != nil {
		b.setCtx(context.TODO())
		b.setCtxCancel(nil)
		cancel()
	}

	// set new context
	ctx, cancel := context.WithCancel(context.Background())
	b.setCtx(ctx)
	b.setCtxCancel(cancel)

	_, err := b.connect()
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitClient) Start() error {
	if b.url == "" {
		return fmt.Errorf("AMQP url is nil")
	}

	err := b.reconnect()
	if err != nil {
		return err
	}

	err = b.initConsume(b.getConn())
	if err != nil {
		return err
	}

	return nil
}

func (b *RabbitClient) initConsume(conn *amqp.Connection) error {
	if conn == nil {
		return fmt.Errorf("AMQP connection is nil")
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	// consumer
	msgs, err := ch.Consume(
		b.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if !b.checkClose() {
				b.Start()
			}
		}()
		ctx := b.getCtx()
		cancel := b.getCtxCancel()
		for {
			select {
			case msg, open := <-msgs:
				if !open {
					cancel()
					logs.Warn(fmt.Sprintf("%s amqp channel was closed", b.queueName))
					return
				}
				if b.callBack != nil {
					b.callBack(msg.Body)
				}
			case <-ctx.Done():
				logs.Warn("amqp channel was closed by call func Closes")
				return
			}
		}
	}()

	// body := "hello listener"
	// ch.Publish("", b.queueName, false, false, amqp.Publishing{
	// 	DeliveryMode: amqp.Persistent,
	// 	ContentType:  "text/plain",
	// 	Body:         []byte(body),
	// })

	return nil
}
