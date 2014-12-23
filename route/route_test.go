package route

import (
	"github.com/Maksadbek/go-ws-daemon/conf"
	"strings"
	"testing"
)

func TestInitialize(t *testing.T) {
	mockConf := `
[ds]
    [ds.mysql]
    dsn = "root:zqwW4XYLzNwN3Dsa@tcp(54.72.185.137:3306)/test"
    limit = 20

    [ds.redis]
    port = 6379
    chan = "orders"
    
[ws]
    port = 1234

`

	reader := strings.NewReader(mockConf)
	config, err := conf.Read(reader)
	if err != nil {
		t.Error(err)
	}

	err = Initialize(config)
	if err != nil {
		t.Error(err)
	}
}
