package conf

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {

	mockConf := `
[ds]
    [ds.mysql]
		dsn = "root:toor@tcp(11.22.33.44:3306)/test"
	[ds.redis]
		port = 6379
		chan = "orders"
[ws]
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
	if got := app.WS.Port; got != want {
		t.Errorf("Websocket Port %d, want %d", got, want)
	}
}
