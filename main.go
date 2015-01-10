package main

import (
	"html/template"
	"io/ioutil"
	"log"

	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/route"

	"net/http"
	"strings"
)

func main() {
	d, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(string(d))

	config, err := conf.Read(reader)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	route.Initialize(config, t)

	http.HandleFunc("/", route.Index)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/orders", route.GetLastOrders)

	log.Fatal(http.ListenAndServe(config.SRV.IP+config.SRV.Port, nil))
}
