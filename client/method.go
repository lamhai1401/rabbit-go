package client

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (b *RabbitClient) getConn() *amqp.Connection {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.conn
}

func (b *RabbitClient) setConn(conn *amqp.Connection) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.conn = conn
}

func (b *RabbitClient) getCtx() context.Context {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.ctxDone
}

func (b *RabbitClient) setCtx(ctx context.Context) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.ctxDone = ctx
}

func (b *RabbitClient) getCtxCancel() context.CancelFunc {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.ctxCancel
}

func (b *RabbitClient) setCtxCancel(ctx context.CancelFunc) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.ctxCancel = ctx
}

func (b *RabbitClient) checkClose() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.isClosed
}

func (b *RabbitClient) setClose(state bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.isClosed = state
}
