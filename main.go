package main

import (
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	d, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(string(d))

	config, err := conf.Read(reader)
	fmt.Println(config)
	if err != nil {
		log.Fatal(err)
	}
}
