package route

import (
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/datastore"
	_ "html/template"
	"net/http"
)

var orderLimit int

func Initialize(config conf.App) error {
	err := datastore.Initialize(config.DS.Mysql.DSN, config.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = config.DS.Mysql.Limit
	orders, _ := datastore.GetLast(10)
	fmt.Printf("%+v\n", orders)
	return err
}

func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	//	lastOrders, _ := datastore.GetLast(orderLimit)
	fmt.Fprintf(w, "hello")
}
