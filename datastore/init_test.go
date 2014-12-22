package datastore

import (
	"fmt"
	"testing"
)

func TestInitialize(t *testing.T) {
	DSN := "root:zqwW4XYLzNwN3Dsa@tcp(54.72.185.137:3306)/test"
	err := Initialize(DSN, 6379)
	if err != nil {
		t.Error(err)
	}
	want := 20
	row, err := GetLast(want)
	if err != nil {
		t.Error(err)
	}

	if len(row) != want {
		t.Errorf("want %d got %d", want, len(row))
	}
	fmt.Printf("%+v\n", row)
}
