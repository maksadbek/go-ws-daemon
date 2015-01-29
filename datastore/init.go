package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nicksnyder/go-i18n/i18n"
)

var db *sql.DB
var T i18n.TranslateFunc

func Initialize(DSN string, redisPort int, i18nPath string) (err error) {
	//connect to sql db
	db, err = sqlConnect(DSN)
	if err != nil {
		return err
	}
	i18n.MustLoadTranslationFile(i18nPath)
	T, err = i18n.Tfunc("ru-ru")
	return err
}
