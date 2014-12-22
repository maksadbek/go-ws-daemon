package datastore

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var redisPool *redis.Pool

func Initialize(DSN string, redisPort int) (err error) {
	//connect to sql db
	db, err = sqlConnect(DSN)
	if err != nil {
		return err
	}
	return err
}
