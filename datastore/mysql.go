package datastore

import (
	"database/sql"
	"strconv"
	"log"
	"github.com/go-sql-driver/mysql"
)

type orderLog struct {
	OrderID    int
	InsertDate string
	Status     string
	Name       string
	CarNum     string
	StCode     int
	ClientName string
	ClientPhone     string
}

type activeOrder struct {
	ID          int `db: "id"`
	Status      string
	OrderTime   string `db: "time_order"`
	Companies   string
	ClientPhone string `db: "client_phone_number`
	CarNum      string `db: "car_number"`
	DriverPhone string `db: "driver_phone"`
	StCode      int    `db: "status"`
	ClientName  string `db: "client_name"`
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
func GetAllActiveOrders(fleet int, last int) (Order, error) {
	var where string
	if fleet == 0 {
		where = " "
	} else {
		where = "AND o.companies = '" + strconv.Itoa(fleet) + "'"
	}

	query :=
		`
		SELECT
				o.id,
				o.STATUS,
				o.time_order,
				o.companies,
				c.Mobile AS client_phone_number,
				d.driver_phone,
				CONCAT(u.name, ' ', u.NUMBER) AS car_number,
				c.FirstName as client_name

		FROM
			max_taxi_incoming_orders o
			LEFT OUTER JOIN max_taxi_server_clients c ON c.ClientID = o.client_id
			LEFT OUTER JOIN max_drivers d ON d.id = o.driver_id
			LEFT OUTER JOIN max_units u ON u.id = d.unit_id
			LEFT OUTER JOIN max_users us ON us.id = o.user_id
		WHERE
		      o.driver_id <> 0 ` +
			where +
			`ORDER BY o.time_order DESC ` +
			` LIMIT 0, ` + strconv.Itoa(last)
	orders := make(Order, last)
	rows, err := db.Query(query)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	// n is iteration index for each
	var n = 0

	for rows.Next() {
		var tmpOrderTime mysql.NullTime
		var tmpClientPhone []byte
		var tmpCarNum []byte
		var tmpDriverPhone []byte
		var status int
		var tmpClientName []byte
		var tmpCompanies string

		err := rows.Scan(
			&orders[n].ID,
			&status,
			&tmpOrderTime,
			&tmpCompanies,
			&tmpClientPhone,
			&tmpCarNum,
			&tmpDriverPhone,
			&tmpClientName,
		)
		orders[n].ClientPhone = string(tmpClientPhone)
		orders[n].CarNum = string(tmpCarNum)
		orders[n].DriverPhone = string(tmpDriverPhone)
		orders[n].ClientName = string(tmpClientName)

		if err != nil {
			return orders, err
		}

		if tmpOrderTime.Valid {
			orders[n].OrderTime = tmpOrderTime.Time.Format("2 01 2006 at 15:04")
		} else {
			orders[n].OrderTime = ""
		}
		orders[n].Status = T(strconv.Itoa(status))
		orders[n].Companies = tmpCompanies
		orders[n].StCode = status

		n += 1
	}
	return orders, err
}

func GetAllOrderLogs(where Where, last int) (Fleet, error) {
	getALLquery := ` 
		SELECT order_id, 
				date_insert, 
				status, 
				carNum, 
				driverName,
		        SUBSTRING_INDEX(Name_Mobile, ':>', -1) Mobile,
		        SUBSTRING_INDEX(Name_Mobile, ':>',  1) FirstName
		FROM (
		    SELECT l.order_id, 
		    		l.date_insert, 
		    		l.status,
		    (SELECT u.number 
		     FROM max_units u 
		     WHERE u.id = l.unit_id
		     ) as carNum,
		    (SELECT d.driver_name 
		    FROM max_drivers d 
		    WHERE d.unit_id = l.unit_id
			LIMIT 1) as driverName,
		    (SELECT CONCAT(cl.FirstName, ':>', cl.Mobile)
		    FROM max_taxi_server_clients cl
		    WHERE cl.ClientID = 
		    	(SELECT io.client_id 
				FROM max_taxi_incoming_orders io
		        WHERE io.client_id = cl.ClientID AND io.id = l.order_id
		        LIMIT 1)
		    ) Name_Mobile
		    FROM max_taxi_deamon_log l
		    WHERE l.` + where.Field +
				` ` +
				where.Crit +
				` ` +
				where.Value +
				` ORDER BY l.id DESC 
				LIMIT 0, ? ) t`

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
		var tmpClientName []byte
		var tmpClientPhone []byte
		var tmpName []byte
		var tmpCarNum []byte

		err = rows.Scan(
			&fleet[n].OrderID,
			&insertDate,
			&status,
			&tmpCarNum,
			&tmpName,
			&tmpClientPhone,
			&tmpClientName,
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
		fleet[n].ClientPhone = string(tmpClientPhone)
		fleet[n].ClientName = string(tmpClientName)
		fleet[n].Name = string(tmpName)
		fleet[n].CarNum = string(tmpCarNum)
		log.Println(fleet[n].Status)
		log.Println(strconv.Itoa(status))
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

func ToNextSt(id int, status int) error {

	query := `UPDATE max_taxi_incoming_orders 
				SET status = ?, next_step = ? 
				WHERE id = ?`
	var newSt string
	var err error
	switch status {
	case 1:
		newSt = "5"

		_, err = db.Exec(query, newSt, newSt, strconv.Itoa(id))
		if err != nil {
			return err
		}
	case 5:
		newSt = "6"
		_, err = db.Exec(query, newSt, newSt, strconv.Itoa(id))
		if err != nil {
			return err
		}
	case 2:
	case 4:
	case 9:
		q := `UPDATE max_taxi_incoming_orders
			 SET next_step = '1', order_attached = 1,
			time_order = now(), status = 1
			 where id = ?`
		_, err = db.Exec(q, strconv.Itoa(id))
		if err != nil {
			return err
		}
	}
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
