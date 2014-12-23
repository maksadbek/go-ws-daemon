package route

import (
	"fmt"
	"github.com/Maksadbek/go-ws-daemon/conf"
	"github.com/Maksadbek/go-ws-daemon/datastore"
	"html/template"
	"net/http"
)

var orderLimit

func Initialize(conf App) error {
	err := datastore.Initialize(conf.DS.Mysql.DSN, conf.DS.Redis.Port)
	if err != nil {
		return err
	}
	orderLimit = conf.DS.Mysql.Limit
	return err
}

func GetLastOrders(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v", datastore.GetLast(lastOrder))
}
