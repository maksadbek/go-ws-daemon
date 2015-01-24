package route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	_ "time"

	"github.com/Maksadbek/go-ws-daemon/conf"
	ds "github.com/Maksadbek/go-ws-daemon/datastore"
)

var orderLimit int //SQL limit for orders
var logLimit int   //SQL limit for daemon logs

var t *template.Template

//Initialize the templates, ds
func Initialize(config conf.App, temp *template.Template) error {
	t = temp
	err := ds.Initialize(config.DS.Mysql.DSN, config.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = config.DS.Mysql.OrderLimit
	logLimit = config.DS.Mysql.LogLimit
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
	log.Println(m)

	h, ok := m["hash"]
	if ok && h[0] == webSiteCookies.Value {
		hash := h[0]
		log.Println("success")

		// get command,
		// activate, next, cancel
		cmd, cmdOk := m["cmd"]

		//get order id to impl. cmd
		orderID, orderOk := m["id"]

		if cmdOk && orderOk {
			id, err := strconv.Atoi(orderID[0])
			if err != nil {
				panic(err)
			}
			switch cmd[0] {
			case "cancel":
				log.Println(cmd[0], orderID[0])
				log.Println("cancel")
				err = ds.CancelActOrder(id)
				if err != nil {
					panic(err)
				}
			case "next":
				log.Println("next")
				st, stOk := m["status"]
				if stOk {
					log.Println("st OK")

					//cast status from str to int
					status, err := strconv.Atoi(st[0])
					if err != nil {
						panic(err)
					}

					//pass st and order id
					err = ds.ToNextSt(id, status)
					if err != nil {
						panic(err)
					}
				}
			case "activate":
				log.Println("activate")
				err = ds.ActivateOrder(id)
				if err != nil {
					panic(err)
				}
			}
		} else {
			t.ExecuteTemplate(w, "activeOrders", hash)
		}
	} else {
		log.Println("failure: hash do not match")
		fmt.Fprintf(w, "login fail")
	}
}

//GetLastOrders n orders and send in JSON
func GetOrders(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	sendErr(w, err)

	f, fleetOk := m["fleet"]

	// Check for existence of the type param
	if t, ok := m["type"]; ok {
		reqType := t[0]
		switch reqType {
		case "logs":
			var order ds.Fleet
			var err error
			//if fleet is given, filter by fleet
			if fleetOk {
				fleetID := f[0] // Get the param from array
				if fleetID == "" {
					order, err = ds.GetAllOrderLogs(
						ds.Where{
							Field: "taxi_fleet_id",
							Crit:  "",
							Value: "IS NOT NULL",
						},
						logLimit,
					)
				} else {
					order, err = ds.GetAllOrderLogs(
						ds.Where{
							Field: "taxi_fleet_id",
							Crit:  "=",
							Value: fleetID,
						},
						logLimit,
					)
				}
			}
			if sendErr(w, err) {
				return
			}
			orderJSON, err := json.Marshal(order)
			if sendErr(w, err) {
				return
			}
			sendOrders(w, orderJSON)
		case "orders":
			var orderJSON []byte
			var err error
			if fleetOk { //if fleet is given, filter by fleet
				//convert fleetID from string to int
				fleet, err := strconv.Atoi(f[0])
				if sendErr(w, err) {
					return
				}
				order, err := ds.GetAllActiveOrders(fleet, orderLimit)
				orderJSON, err = json.Marshal(order)
			} else { //if fleet is not given, do not filter
				order, err := ds.GetAllActiveOrders(0, orderLimit)
				if sendErr(w, err) {
					return
				}
				orderJSON, err = json.Marshal(order)
			}
			if sendErr(w, err) {
				return
			}
			sendOrders(w, orderJSON)
		}
	}
}

func sendOrders(w http.ResponseWriter, orders []byte) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(orders))
}

func sendErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		fmt.Printf("%+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return true
	}
	return false
}

func GetOrderLogs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI())
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("error")
		panic(err)
	}

	hash := m["hash"][0]
	fleet, fleetOk := m["fleet"]

	webSiteCookies, err := r.Cookie("PHPSESSID")
	if err != nil {
		panic(err)
	}

	if hash == webSiteCookies.Value {
		log.Println("success")
		if fleetOk {
			t.ExecuteTemplate(w, "index", fleet)
		} else {
			t.ExecuteTemplate(w, "index", nil)
		}
	} else {
		log.Println("failure")
		fmt.Fprintf(w, "login fail")
	}
}
