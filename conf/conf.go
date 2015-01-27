package conf

import (
	"github.com/BurntSushi/toml"
	"io"
)

type Datastore struct {
	Redis struct {
		Port int
		Chan string
	}

	Mysql struct {
		DSN        string
		OrderLimit int
		LogLimit   int
	}
}

type Server struct {
	IP   string
	Port string
}
type I18n struct {
	Path string
}

type App struct {
	DS    Datastore
	SRV   Server
	I18n  I18n
	Order Order
}

type Order struct {
	Attrs []string
}

func Read(r io.Reader) (config App, err error) {
	_, err = toml.DecodeReader(r, &config)
	return config, err
}
