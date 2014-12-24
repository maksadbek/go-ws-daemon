package main

import (
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/route"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	d, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	reader, err := strings.NewReader(string(d))
	if err != nil {
		log.Fatal(err)
	}

	config, err := conf.Read(reader)
	if err != nil {
		log.Fatal(err)
	}

	err = route.Initialize(config.DS.Mysql.DSN, 6379)
	if err != nil {
		log.Fatal(err)
	}
}
