package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type orderLog struct {
	Id      int
	OrderID int
}

type Fleet []orderLog

func sqlConnect(DSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	return db, err
}

func GetLast(last int) (Fleet, error) {
	rows, err := db.Query("SELECT id, order_id from max_taxi_deamon_log LIMIT ?", last)
	fleet := make(Fleet, last)
	if err != nil {
		return fleet, err
	}

	defer rows.Close()
	var n = 0
	for rows.Next() {
		if err := rows.Scan(&fleet[n].Id, &fleet[n].OrderID); err != nil {
			return fleet, err
		}
		n += 1
	}
	return fleet, err
}
