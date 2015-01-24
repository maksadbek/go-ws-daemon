package datastore

import "testing"

func TestGetAllOrderLogs(t *testing.T) {
	DSN := "root:zqwW4XYLzNwN3Dsa@tcp(54.72.185.137:3306)/test"
	err := Initialize(DSN, 6379)
	if err != nil {
		t.Error(err)
	}

	row, err := GetAllActiveOrders(10)
	want := 10
	if len(row) != want {
		t.Error("want %d got %d", want, len(row))
	}
}
