package datastore

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nicksnyder/go-i18n/i18n"
)

var db *sql.DB
var T i18n.TranslateFunc

func Initialize(DSN string, redisPort int) (err error) {
	//connect to sql db
	db, err = sqlConnect(DSN)
	if err != nil {
		return err
	}
	i18n.MustLoadTranslationFile("../ru-ru.all.json")
	T, err = i18n.Tfunc("ru-RU")
	log.Println("i18n --> ", T("1"))
	return err
}
