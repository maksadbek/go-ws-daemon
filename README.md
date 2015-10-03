go-ws-daemon
============
This is a Go WebSocket server that updates data on UI when the data changes.
The UI part is done using ReactJS.

Table structure
```
+------------------------+-------------+------+-----+---------------------+----------------+
| Field                  | Type        | Null | Key | Default             |Extra           |
+------------------------+-------------+------+-----+---------------------+----------------+
| id                     | int(11)     | NO   | PRI | NULL                |auto_increment  |
| order_id               | int(11)     | NO   | MUL | NULL                |                |
| conn_driver_id         | int(11)     | NO   |     | NULL                |                |
| date_insert            | timestamp   | NO   |     | CURRENT_TIMESTAMP   |                |
| dtime_click            | timestamp   | YES  |     | NULL                |                |
| status                 | smallint(2) | NO   | MUL | NULL                |                |
| taxi_fleet_id          | int(11)     | NO   |     | NULL                |                |
| unit_id                | int(11)     | NO   |     | NULL                |                |
| drv_accepted_date_time | datetime    | NO   |     | 0000-00-00 00:00:00 |                |
| active                 | tinyint(1)  | NO   |     | 1                   |                |
+------------------------+-------------+------+-----+---------------------+----------------+
```

The sample JSON response
```JSON
{
  "id": 1252247,
  "client_id": 583,
  "status": "ok",
  "from_address": "41.294003 69.246033",
  "date": "2015-01-13 19:57:15",
  "time_order": "2015-01-13 20:00:11",
  "companies": 202,
  "tariffID": 2,
  "client_phone_number": "998000000000",
  "client_name": "Без номера",
  "car_number": "10 AA 12",
  "driver_phone": "998 90 1234",
  "user_name": "minimal"
}
```

To retrieve this data, the following query is runned.

```SQL
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
    ((o.driver_id <> 0 AND o.driver_id in (SELECT id FROM max_drivers WHERE fleet_id = 202))
    OR
    (o.driver_id = 0 AND companies like '%202%'))
    LIMIT 0,10
```

##API##
Get all orders filtered by fleet id
```
  curl -XGET localhost:9000/active-orders?fleet=<FLEET_NUMBER>&hash=<HASH_FROM_PHPSESSID>
```

Get all orders
```Curl
curl -XGET localhost:9000/active-orders?hash=<HASH_FROM_PHPSESSID>
```
