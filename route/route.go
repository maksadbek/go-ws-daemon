package route

import (
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/datastore"
	"html/template"
	"net/http"
)

var orderLimit int
var t *template.Template

func Initialize(config conf.App, temp *template.Template) error {
	t = temp
	err := datastore.Initialize(config.DS.Mysql.DSN, config.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = config.DS.Mysql.Limit
	return err
}

func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetLastOrders")
	orders, _ := datastore.GetLast(orderLimit)
	fmt.Printf("%+v", orders)
	t.ExecuteTemplate(w, "orders", orders)
}
