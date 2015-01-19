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
                sid: 12224
            user_id: 149
          client_id: 583
          driver_id: 0
             status: NULL
     push_life_time: 20
driver_arrival_time: 0
        dateRequest: 00:00:00
         from_adres: 41.294003 69.246033
    sub_region_from: 
   coord_from_adres:        ���㡥D@'/2�OQ@
    poi_from_adress: Шиномонтаж - 41м, Детская Школа №18 музыки и исскуств (Bolalar musiqa va san`at maktabi) - 163м, Детский сад №68 - 244м, Школа №73 - 246м, Опорный пункт милиции №433 - 318м, Детский сад №73 - 333м, Отряд пожарной безопасности - 347м, Семейная поликлиника №57 - 347м, Клиника Sergo Dental - 352м, Проектный институт ТошкентбошпланЛИТИ - 356
           to_adres: 41.294003 69.246033
      sub_region_to: 
     coord_to_adres:        ���㡥D@'/2�OQ@
               date: 2015-01-13 19:57:15
         time_order: 2015-01-13 20:00:11
               reqs: 2
          companies: 202
          orderFrom: 1
           distance: 2000
        description: oooo
    dr_arrived_time: 0000-00-00 00:00:00
    client_sat_time: 0000-00-00 00:00:00
order_finished_time: 0000-00-00 00:00:00
        status_time: 2015-01-13 19:57:15
   order_del_status: 0
      order_details: 
          next_step: 0
     order_attached: 0
           tariffID: 2
     fleet_tarif_id: 3
 coord_from_adres_t: POINT(41.294003 69.246033)
   coord_to_adres_t: POINT(41.294003 69.246033)
client_phone_number: 998000000000
        client_name: Без номера
         car_number: NULL
            unit_id: NULL
           vaqtRazn: 0
         vaqtRazn_7: 0
        BUDvaqtRazn: 0
       driver_phone: NULL
          user_name: newmax
```
