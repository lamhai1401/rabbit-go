package client

import "context"

type RabbitClient struct {
	ctxDone   context.Context    // to signal reconnect
	ctxCancel context.CancelFunc // to cancel ctx
}
