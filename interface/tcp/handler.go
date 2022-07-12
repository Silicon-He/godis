package tcp

import (
	"context"
	"net"
)

// this package only implement handle a connection

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
