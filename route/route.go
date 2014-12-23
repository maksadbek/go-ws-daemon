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
	return err
}

func GetLastOrders(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	lastOrders, _ := datastore.GetLast(orderLimit)
	fmt.Fprintf(w, "%+v", lastOrders)
=======
	fmt.Fprintf(w, "%+v", datastore.GetLast(orderLimit))
>>>>>>> bfe3682e62f1473954d564b21829b949be5cc8dc
}
