package main

import (
	"fmt"
	"godis/config"
	"godis/lib/logger"
	"godis/tcp"
)

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address:       fmt.Sprintf("%s:%d", config.RedisCfg.IP, config.RedisCfg.Port),
			MaxConnectNum: int(config.RedisCfg.MaxConn),
		},
		tcp.NewHandler())
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
