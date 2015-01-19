package datastore

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize(DSN string, redisPort int) (err error) {
	//connect to sql db
	db, err = sqlConnect(DSN)
	if err != nil {
		return err
	}
	return err
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
