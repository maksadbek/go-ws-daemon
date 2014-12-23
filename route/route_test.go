package route

import (
	"github.com/Maksadbek/go-ws-daemon/conf"
	"log"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {
	mockConf := `
[ds]
    [datastore.mysql]
    dsn = "root:zqwW4XYLzNwN3Dsa@tcp(54.72.185.137:3306)/test"
    limit = 20

    [datastore.redis]
    port = 6379
    chan = "orders"
    
[ws]
    port = 1234

`

	reader := strings.NewReader(mockConf)
	config, err := conf.Read(reader)
<<<<<<< HEAD
	log.Println(config.DS.Mysql.DSN)
=======
	log.Println(config)
>>>>>>> bfe3682e62f1473954d564b21829b949be5cc8dc
	if err != nil {
		t.Error(err)
	}

	err = Initialize(config)
	if err != nil {
		t.Error(err)
	}
}
