package utils

import (
	"fmt"

	"github.com/go-ini/ini"
)

type Face struct {
	Minions string
	Group   string
}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}
type Config struct {
	ConnectAll bool
	RunMode    string
	Port       string
	Face
	Database
}

var config *Config

func InitConfig() {
	cfg, err := ini.Load("config.ini")
	config = new(Config)
	err = cfg.MapTo(config)
	fmt.Printf("%+v\n", cfg)
	if err != nil {
		fmt.Println(err)
	}
}

func GetCfg() *Config {
	return config
}
