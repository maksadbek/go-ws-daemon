package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/route"

	"net/http"
	"strings"
)

func main() {
	confPath := flag.String("conf", "config.toml", "configuration file")
	viewsPath := flag.String("views", "views", "views folder")
	flag.Parse()
	d, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	reader := strings.NewReader(string(d))

	config, err := conf.Read(reader)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles(
		*viewsPath+"/index.html",
		*viewsPath+"/header.html",
		*viewsPath+"/active.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	route.Initialize(config, t)

	http.HandleFunc("/", route.Index)
	http.Handle("/favicon.ico", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/orders", route.GetLastOrders)

	log.Fatal(http.ListenAndServe(config.SRV.IP+config.SRV.Port, nil))
}
