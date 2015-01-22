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
	ID          int `db: "id"`
	ClientID    int `db: "client_id"`
	Status      string
	AddrFrom    string `db: "from_adres"`
	StCode      int    `db: "status"`
	Date        string `db: "date"`
	OrderTime   string `db: "time_order"`
	Companies   string `db: "comanies"`
	TariffID    string `db: "tariffID"`
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

// GetAllActiveOrders can be used to get active orders
// without filtration for fleets
func GetAllActiveOrders(last int) (Order, error) {
	query :=
		`
		SELECT
			o.id,
			o.STATUS,
			o.client_id,
			o.from_adres,
			o.DATE,
			o.time_order,
			o.companies,
			o.tariffID,
			c.Mobile AS client_phone_number,
			d.driver_phone,
			us.login AS user_name,
			CONCAT(u.name, ' ', u.NUMBER) AS car_number,
			c.FirstName AS client_name
		FROM
			max_taxi_incoming_orders o
			LEFT OUTER JOIN max_taxi_server_clients c ON c.ClientID = o.client_id
			LEFT OUTER JOIN max_drivers d ON d.id = o.driver_id
			LEFT OUTER JOIN max_units u ON u.id = d.unit_id
			LEFT OUTER JOIN max_users us ON us.id = o.user_id
		WHERE
		      o.driver_id <> 0
		  LIMIT 0, ` + strconv.Itoa(last)

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
		var tmpTariffID []byte
		var tmpClientPhone []byte
		var tmpClientName []byte
		var tmpCarNum []byte
		var tmpDriverPhone []byte
		var tmpUserName []byte
		var status int

		err := rows.Scan(
			&orders[n].ID,
			&status,
			&orders[n].ClientID,
			&orders[n].AddrFrom,
			&tmpDate,
			&tmpOrderTime,
			&orders[n].Companies,
			&tmpTariffID,
			&tmpClientPhone,
			&tmpClientName,
			&tmpCarNum,
			&tmpDriverPhone,
			&tmpUserName,
		)
		orders[n].TariffID = string(tmpTariffID)
		orders[n].ClientPhone = string(tmpClientPhone)
		orders[n].ClientName = string(tmpClientName)
		orders[n].CarNum = string(tmpCarNum)
		orders[n].DriverPhone = string(tmpDriverPhone)
		orders[n].UserName = string(tmpUserName)

		if err != nil {
			return orders, err
		}

		if tmpDate.Valid {
			orders[n].Date = tmpDate.Time.Format("2 01 2006 at 15:04")
		} else {
			orders[n].Date = ""
		}

		if tmpOrderTime.Valid {
			orders[n].OrderTime = tmpOrderTime.Time.Format("2 01 2006 at 15:04")
		} else {
			orders[n].OrderTime = ""
		}

		orders[n].Status = T(strconv.Itoa(status))
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
	//n is for data index
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

		fleet[n].Status = T(strconv.Itoa(status))
		fleet[n].StCode = status

		n += 1
	}

	return fleet, err
}

func CancelActOrder(id int) error {
	query := `UPDATE max_taxi_incoming_orders
				SET status = 9 
				WHERE id = ?`
	_, err := db.Exec(query, strconv.Itoa(id))
	return err
}

func ToNextSt(id int) error {
	query := `UPDATE max_taxi_incoming_orders 
				SET status = ?, next_step = ? 
				WHERE id = ?`
	_, err := db.Exec(query, strconv.Itoa(id))
	return err
}

func ActivateOrder(id int) error {
	queryDaemonLog := `UPDATE max_taxi_deamon_log SET active = 0
	         WHERE order_id = ?`
	_, err := db.Exec(queryDaemonLog, strconv.Itoa(id))
	if err != nil {
		return err
	}
	queryIncOrders := `UPDATE max_taxi_incoming_orders SET 
			status = 0, 
			driver_id = 0, 
			date= now(), 
			time_order = now()
			WHERE id = ?`
	_, err = db.Exec(queryIncOrders, strconv.Itoa(id))
	return err
}
