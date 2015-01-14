package datastore

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"strconv"
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

func statusToDesc(status int) string {
	switch status {
	case 0:
		return "В ожидании, новый заказ"
	case -2:
		return "Клиент оповешен"
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
	case 44:
		return "Передаем другому водителю"
	case 20:
		return "Водитель оповешен"

	}
	return "Статус не определен: " + strconv.Itoa(status)
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
