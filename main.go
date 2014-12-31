package main

import (
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/route"
	"html/template"
	"io/ioutil"
	"log"
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

	t, err := template.ParseFiles("views/index.html", "views/orders.html")
	if err != nil {
		log.Fatal(err)
	}

	route.Initialize(config, t)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", route.Index)
	http.HandleFunc("/orders", route.GetLastOrders)

	log.Fatal(http.ListenAndServe(":4000", nil))
}
