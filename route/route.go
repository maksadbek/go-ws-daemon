package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	_ "time"

	"github.com/Maksadbek/go-ws-daemon/conf"
	ds "github.com/Maksadbek/go-ws-daemon/datastore"
)

var orderLimit int
var t *template.Template

//Initialize the templates, ds
func Initialize(config conf.App, temp *template.Template) error {
	t = temp
	err := ds.Initialize(config.DS.Mysql.DSN, config.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = config.DS.Mysql.Limit
	return err
}

//GetLastOrders n orders and send in JSON
func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	fleetID := m["fleet"][0]

	orders, err := ds.GetAll(ds.Where{Field: "taxi_fleet_id", Crit: "=", Value: fleetID}, orderLimit)
	if err != nil {
		fmt.Println(err)
	}

	var ordersInJson []byte
	ordersInJson, err = json.Marshal(orders)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(ordersInJson))
}

//Index file
func Index(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("error")
		panic(err)
	}

	hash := m["hash"][0]
	webSiteCookies, err := r.Cookie("PHPSESSID")
	if err != nil {
		log.Printf("%+v\n", err)
	}

	if hash == webSiteCookies.Value {
		log.Println("success")
		t.ExecuteTemplate(w, "index", nil)
	} else {
		log.Println("failure")
		fmt.Fprintf(w, "login fail")
	}
}
