package conf

import (
	"log"
	"testing"
)

func TestRead(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Error(err)
	}
	log.Println(config.Datastore["mysql_test"].DSN)
}
