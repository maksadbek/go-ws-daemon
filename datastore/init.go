package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize(DSN string, redisPort int) (err error) {
	//connect to sql db
	db, err = sqlConnect(DSN)
	if err != nil {
		return err
	}
	return err
}
