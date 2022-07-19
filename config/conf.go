package config

import (
	"github.com/BurntSushi/toml"
	"path/filepath"
	"runtime"
)

var (
	nowPath  string
	Config   *config
	RedisCfg *redisCfg
)

func init() {
	_, nowPath, _, _ = runtime.Caller(0)
	cfgFilePath := filepath.Join(filepath.Dir(nowPath), "config.toml")
	if _, err := toml.DecodeFile(cfgFilePath, &Config); err != nil {
		panic(err)
	}
	RedisCfg = &Config.Redis
}

type config struct {
	Redis redisCfg
}

type redisCfg struct {
	IP      string `toml:"ip"`
	Port    int32  `toml:"port"`
	MaxConn int32  `toml:"max_conn"`
}
