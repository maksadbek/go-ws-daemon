package datastore

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"strconv"
)

type orderLog struct {
	ID           int
	OrderID      int
	DriverConnID int
	InsertDate   string
	ClickTime    string
	Status       string
	TaxiFleetID  int
	UnitID       int
	Active       int
}

type Fleet []orderLog

func sqlConnect(DSN string) (*sql.DB, error) {
	db, err := sql.Open("mysql", DSN)
	return db, err
}

func statusToDesc(status int) string {
	switch status {
	case 0:
		return "В ожидании, новый заказ"
	case 1:
		return "Принять Водителем"
	case 2:
		return "Водитель не найден"
	case 3:
		return "В процессе обработки"
	case 4:
		return "Такси отменил"
	case 5:
		return "Машина приехала"
	case 6:
		return "Клиент на машине"
	case 7:
		return "Конец"
	case 8:
		return "Клиент отменил"
	case 9:
		return "Сервер отменил"
	case 10:
		return "Доехал"
	case 11:
		return "Водитель отменил: Технический неполадка"
	case 12:
		return "Водитель отменил: Не успею"
	case 13:
		return "Водитель отменил: Другая причина отказа"
	case 14:
		return "Водитель отменил: Клиент не выходил"
	case 15:
		return "Водитель отменил: Другая причина отказа"
	case 16:
		return "Водитель отменил: Клиент отказался"
	case 17:
		return "Водитель отменил: Технический неполадка"
	case 18:
		return "Водитель отменил: Другая причина отказа"
	}
	return "Статус не определен: " + strconv.Itoa(status)
}

func GetLast(last int) (Fleet, error) {
	getALLquery := ` 
			SELECT 
				id, 
				order_id, 
				conn_driver_id, 
				date_insert, 
				dtime_click, 
				status, 
				taxi_fleet_id, 
				unit_id, 
				active 
				FROM max_taxi_deamon_log 
				ORDER BY id DESC 
				LIMIT 0, ?
			`
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
		var clickTime mysql.NullTime
		err = rows.Scan(
			&fleet[n].ID,
			&fleet[n].OrderID,
			&fleet[n].DriverConnID,
			&insertDate,
			&clickTime,
			&status,
			&fleet[n].TaxiFleetID,
			&fleet[n].UnitID,
			&fleet[n].Active,
		)

		if insertDate.Valid {
			fleet[n].InsertDate = insertDate.Time.Format("2 Jan 2006 at 15:04")
		} else {
			fleet[n].InsertDate = ""
		}

		if clickTime.Valid {
			fleet[n].ClickTime = clickTime.Time.Format("2 Jan 2006 at 15:04")
		} else {
			fleet[n].ClickTime = ""
		}

		fleet[n].Status = statusToDesc(status)

		n += 1
	}

	return fleet, err
}
