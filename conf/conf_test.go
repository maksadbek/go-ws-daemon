package conf

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	mockConf := `
[ds]
    [ds.mysql]
		dsn = "root:zqwW4XYLzNwN3Dsa@tcp(54.72.185.137:3306)/test"
		limit = 20
	[ds.redis]
		port = 6379
		chan = "orders"
[srv]
	port = 1234
`

	r := strings.NewReader(mockConf)
	app, err := Read(r)
	if err != nil {
		t.Error("Read error")
	}
	want := 6379
	if got := app.DS.Redis.Port; got != want {
		t.Errorf("Datastore Redis Port %d, want %d", got, want)
	}

	want = 1234
	if got := app.SRV.Port; got != want {
		t.Errorf("Websocket Port %d, want %d", got, want)
	}

	want = 20
	if got := app.DS.Mysql.Limit; got != want {
		t.Errorf("Mysql limit %d, want %d", got, want)
	}
}
