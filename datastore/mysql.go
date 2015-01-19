package datastore

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type orderLog struct {
	OrderID    int
	InsertDate string
	Status     string
	Name       string
	CarNum     string
	StCode     int
}

type Where struct {
	Field string
	Crit  string
	Value string
}

type Fleet []orderLog

func sqlConnect(DSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	return db, err
}

func GetAll(where Where, last int) (Fleet, error) {
	getALLquery := ` 
			SELECT 
				l.order_id as order_id, 
				l.date_insert as data_insert, 
				l.status as status, 
				u.number as carNum,
				d.driver_name as driverName
				FROM max_taxi_deamon_log l, max_drivers d, max_units u
				WHERE u.id = d.unit_id and l.unit_id = u.id and l.` + where.Field + ` ` + where.Crit + ` ` + where.Value +
		` ORDER BY l.id DESC 
				LIMIT 0, ? `
	rows, err := db.Query(getALLquery, last)
	fleet := make(Fleet, last)
	if err != nil {
		return fleet, err
	}

	defer rows.Close()
	var n = 0
	for rows.Next() {
		var status int
		var insertDate mysql.NullTime
		err = rows.Scan(
			&fleet[n].OrderID,
			&insertDate,
			&status,
			&fleet[n].CarNum,
			&fleet[n].Name,
		)

		if insertDate.Valid {
			fleet[n].InsertDate = insertDate.Time.Format("2 01 2006 at 15:04")
		} else {
			fleet[n].InsertDate = ""
		}

		fleet[n].Status = statusToDesc(status)
		fleet[n].StCode = status

		n += 1
	}

	return fleet, err
}
