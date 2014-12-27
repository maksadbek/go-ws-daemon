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
