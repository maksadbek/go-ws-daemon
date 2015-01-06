package route

import (
	"encoding/json"
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/datastore"
	"html/template"
	"net/http"
)

var orderLimit int
var t *template.Template

//Initialize the templates, datastore
func Initialize(config conf.App, temp *template.Template) error {
	t = temp
	err := datastore.Initialize(config.DS.Mysql.DSN, config.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = config.DS.Mysql.Limit
	return err
}

//GetLastOrders n orders and send in JSON
func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := datastore.GetLast(orderLimit)
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
	t.ExecuteTemplate(w, "index", nil)
}
