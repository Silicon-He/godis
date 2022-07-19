package tcp

// implement tcp server

import (
	"context"
	"godis/interface/tcp"
	"godis/lib/logger"
	"godis/lib/sync/wait"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Address       string
	MaxConnectNum int
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) (err error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	closeChan := make(chan struct{})
	signChan := make(chan os.Signal)
	signal.Notify(signChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			closeChan <- struct{}{}
		}
	}()
	defer close(closeChan)
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// when closed
	go func() {
		<-closeChan
		logger.Info("shutting down listener...")
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var wg *wait.Wait
	// clean connect
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link...")
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}
