package main

import (
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	r "github.com/Maksadbek/go-ws-daemon/route"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getLastOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
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

	r.Initialize(config)
	r.Print("HELLO")
	http.Handle("/assets", http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets/"))))
	http.HandleFunc("/", r.GetLastOrders)
	http.ListenAndServe(":3000", nil)
}
