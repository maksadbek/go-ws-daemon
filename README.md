go-ws-daemon
============

##max_taxi_deamon_log table structure##
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

```
                 id: 1252247
          client_id: 583
             status: NULL
         from_adres: 41.294003 69.246033
               date: 2015-01-13 19:57:15
         time_order: 2015-01-13 20:00:11
          companies: 202
           tariffID: 2
client_phone_number: 998000000000
        client_name: Без номера
         car_number: NULL
       driver_phone: NULL
          user_name: newmax
```

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
