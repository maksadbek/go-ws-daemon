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

func GetActiveOrders(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI())
	webSiteCookies, err := r.Cookie("PHPSESSID")
	if err != nil {
		log.Println("failure: no cookie")
		fmt.Fprintf(w, "login fail")
		return
	}

	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("error")
		panic(err)
	}

	hash, ok := m["hash"]
	if ok && hash[0] == webSiteCookies.Value {
		log.Println("success")

		// get command,
		// activate, next, cancel
		cmd, cmdOk := m["cmd"]

		//get order id to impl. cmd
		orderID, orderOk := m["id"]

		if cmdOk && orderOk {
			switch cmd[0] {
			case "cancel":
				ds.CancelActOrder(orderID[0])
			case "next":
				ds.ToNextSt(orderID[0])
			case "activate":
				ds.ActivateOrder(orderID[0])
			}
		} else {
			t.ExecuteTemplate(w, "activeOrders", nil)
		}
	} else {
		log.Println("failure: hash do not match")
		fmt.Fprintf(w, "login fail")
	}
}

//GetLastOrders n orders and send in JSON
func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI())
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	// Check for existence of the fleet param
	if f, ok := m["fleet"]; ok {
		fleetID := f[0] // Get the param from array
		orders, err := ds.GetAllOrderLogs(
			ds.Where{
				Field: "taxi_fleet_id",
				Crit:  "=",
				Value: fleetID,
			},
			orderLimit,
		)

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
	} else {
		orders, err := ds.GetAllActiveOrders(50)
		//TODO: make a helper function
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		}
		var ordersJSON []byte
		ordersJSON, err = json.Marshal(orders)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(ordersJSON))
	}
}

//Index file
func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI())
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("error")
		panic(err)
	}

	hash := m["hash"][0]
	fleet := m["fleet"][0]

	webSiteCookies, err := r.Cookie("PHPSESSID")
	if err != nil {
		panic(err)
	}

	if hash == webSiteCookies.Value {
		log.Println("success")
		t.ExecuteTemplate(w, "index", fleet)
	} else {
		log.Println("failure")
		fmt.Fprintf(w, "login fail")
	}
}
