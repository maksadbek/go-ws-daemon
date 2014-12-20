package conf

import (
	"github.com/BurntSushi/toml"
	"log"
)

type mysqlTest struct {
	DSN string
}

type redisTest struct {
	IP      string
	Channel string
}

type webSocket struct {
	Port string
}

type App struct {
	Datastore map[string]mysqlTest
	Redis     map[string]redisTest
	WebSocket webSocket
}

func Read() (config App, err error) {
	if _, err := toml.DecodeFile("../config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config, err
}
