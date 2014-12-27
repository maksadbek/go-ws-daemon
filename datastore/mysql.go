package datastore

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "time"
)

type orderLog struct {
	ID            int
	OrderID       int
	DriverConnID  int
	InsertDate    []byte
	ClickTime     []byte
	Status        int
	TaxiFleetID   int
	UnitID        int
	DrvAcceptTime []byte
	Active        int
}

type Fleet []orderLog

func sqlConnect(DSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	return db, err
}

func GetLast(last int) (Fleet, error) {
	rows, err := db.Query("SELECT id, order_id, conn_driver_id, date_insert, dtime_click, status, taxi_fleet_id, unit_id, drv_accepted_date_time, active from max_taxi_deamon_log LIMIT ?", last)
	fleet := make(Fleet, last)
	if err != nil {
		return fleet, err
	}

	defer rows.Close()
	var n = 0
	for rows.Next() {
		err = rows.Scan(&fleet[n].ID, &fleet[n].OrderID, &fleet[n].DriverConnID, &fleet[n].InsertDate,
			&fleet[n].ClickTime, &fleet[n].Status, &fleet[n].TaxiFleetID, &fleet[n].UnitID,
			&fleet[n].DrvAcceptTime, &fleet[n].Active)

		if err != nil {
			return fleet, err
		}
		n += 1
	}

	return fleet, err
}
