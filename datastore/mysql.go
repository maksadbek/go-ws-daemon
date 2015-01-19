package datastore

import (
	"database/sql"
	"strconv"

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

type activeOrder struct {
	Id          int `db: "id"`
	ClientID    int `db: "client_id"`
	Status      string
	AddrFrom    string `db: "from_adres"`
	StCode      int    `db: "status"`
	Date        string `db: "date"`
	OrderTime   string `db: "time_order"`
	Companies   int    `db: "comanies"`
	TariffID    int    `db: "tariffID"`
	ClientPhone string `db: "client_phone_number`
	ClientName  string `db: "client_name"`
	CarNum      string `db: "car_number"`
	DriverPhone string `db: "driver_phone"`
	UserName    string `db: "user_name"`
}

type Where struct {
	Field string
	Crit  string
	Value string
}

type Fleet []orderLog
type Order []activeOrder

func sqlConnect(DSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	return db, err
}

func GetAllActiveOrders(fleetID int, last int) (Order, error) {
	fleet := strconv.Itoa(fleetID)
	query := `
	SELECT
		o.id, 
		o.status, 
		o.client_id,
		o.from_adres,
		o.date,
		o.time_order,
		o.companies, 
		o.tariffID,
		c.Mobile as client_phone_number,
		d.driver_phone, 
		us.login as user_name,
		CONCAT(u.name, ' ', u.number) as car_number,
		c.FirstName as client_name
	FROM
		max_taxi_incoming_orders o
		LEFT OUTER JOIN max_taxi_server_clients c ON c.ClientID = o.client_id
		LEFT OUTER JOIN max_drivers d ON d.id = o.driver_id
		LEFT OUTER JOIN max_units u ON u.id = d.unit_id
		LEFT OUTER JOIN max_users us ON us.id = o.user_id
	WHERE
		(
			(o.driver_id <> 0 
			AND o.driver_id
			IN (SELECT id 
				FROM max_drivers 
				WHERE fleet_id = ` + fleet + `)
			)
			OR (o.driver_id = 0 
				AND companies 
				LIKE '%` + fleet + ` %')
		)
		LIMIT 0,` + strconv.Itoa(last)
	orders := make(Order, last)

	rows, err := db.Query(query)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	// n is iteration index for each
	var n = 0

	for rows.Next() {
		var tmpDate mysql.NullTime
		var tmpOrderTime mysql.NullTime
		var status int
		err := rows.Scan(
			&orders[n].Id,
			&orders[n].ClientID,
			&status,
			&orders[n].AddrFrom,
			&tmpDate,
			&tmpOrderTime,
			&orders[n].Companies,
			&orders[n].TariffID,
			&orders[n].ClientPhone,
			&orders[n].ClientName,
			&orders[n].CarNum,
			&orders[n].DriverPhone,
			&orders[n].UserName,
		)

		if err != nil {
			return orders, err
		}

		if tmpDate.Valid {
			orders[n].Date = tmpDate.Time.Format("2 01 2006 at 15:04")
		} else if tmpOrderTime.Valid {
			orders[n].OrderTime = tmpOrderTime.Time.Format("2 01 2006 at 15:04")
		} else {
			orders[n].Date = ""
			orders[n].OrderTime = ""
		}

		orders[n].Status = statusToDesc(status)
		orders[n].StCode = status

		n += 1
	}
	return orders, err
}

func GetAllOrderLogs(where Where, last int) (Fleet, error) {
	getALLquery := ` 
			SELECT 
				l.order_id as order_id, 
				l.date_insert as data_insert, 
				l.status as status, 
				u.number as carNum,
				d.driver_name as driverName
				FROM max_taxi_deamon_log l, max_drivers d, max_units u
				WHERE u.id = d.unit_id and l.unit_id = u.id and l.` +
		where.Field +
		` ` +
		where.Crit +
		` ` +
		where.Value +
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
		if err != nil {
			return fleet, err
		}

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
