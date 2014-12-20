package datastore

import (
	"database/sql"
	"gihub.com/Maksadbek/go-ws-daemon/conf"
	"github.com/garyburd/redigo"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func Initialize() {
	//get configs
	var config conf.App
	config = conf.Read()
	//connect to sql db
	db, err := sql.Open("mysql", config.Mysql.DSN)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
