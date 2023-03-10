[![Go Vet](https://github.com/patrikcze/go-graph/actions/workflows/workflow.yml/badge.svg)](https://github.com/patrikcze/go-graph/actions/workflows/workflow.yml)

# PoC 
Trying to mess around with GoLang, prepare simple APi which will collect data from ESP32 + BME280 Module, uploads data to MySQL Database and will draw a line chart.

# Updated thx to kX !
- updated makefile
- added workflow
- Updated workflow, using selfhosted runner
- fixed error handling in main.go
- Updated Makefile there were multiple rules in Makefile with the same target name 'build', and the second one is overriding the first one.
- Added comments / used CamelCase for Consts ...

**THIS ALL HAS BEEN DONE WITH HELP OF OpenAI!**

Let's see how far we can go. :-)
I'm messing with this all the time, so appologize for inconsistency.

# Update 27.2.2023

- Trying to review and test on Hardware again. 
- Fixing Problems with runner on Github

# Updates 12.1.2023

- Created docker-compose.yml to use `multi-container` deployment.
- MySQL Server is now running on official `mysql:latest` image. 
- Go API is being compiled in `golang:latest`
- `ubuntu:22.04` has been used for hosting `go-api`.
- `main.go` - Updated API To use `db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(db)/"+dbName+"?parseTime=true")`
- `main.go` - changed CONSTS to variables and reading data from ENV variables.

## MySQL Database

Database created : 
```
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
| temperature_db     |
+--------------------+
5 rows in set (0.00 sec)
```
Permissions added automatically within native MySQL image and entrypoint script. 

```
go-graph-db-1      | 2023-01-12 04:55:03+00:00 [Note] [Entrypoint]: Creating database temperature_db
go-graph-db-1      | 2023-01-12 04:55:03+00:00 [Note] [Entrypoint]: Creating user dbuser
go-graph-db-1      | 2023-01-12 04:55:03+00:00 [Note] [Entrypoint]: Giving user dbuser access to schema temperature_db
go-graph-db-1      | 
go-graph-db-1      | 2023-01-12 04:55:03+00:00 [Note] [Entrypoint]: /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/01.sql
```

Table created and filled with example data by `./init/01.sql` script:
```
mysql> show tables;
+--------------------------+
| Tables_in_temperature_db |
+--------------------------+
| data                     |
+--------------------------+
1 row in set (0.00 sec)
```

Example data: 
```
mysql> select * from data;
+----+---------------------+-------------+----------+----------+
| id | time                | temperature | humidity | pressure |
+----+---------------------+-------------+----------+----------+
|  1 | 2022-01-01 00:00:00 |          20 |       78 |     1023 |
|  2 | 2022-01-01 01:00:00 |          21 |       75 |     1023 |
|  3 | 2022-01-01 02:00:00 |          22 |       34 |     1025 |
|  4 | 2022-01-01 03:00:00 |          23 |       55 |     1022 |
|  5 | 2022-01-01 04:00:00 |          24 |       33 |     1021 |
|  6 | 2022-01-01 05:00:00 |          25 |       12 |     1023 |
|  7 | 2022-01-01 06:00:00 |          26 |       64 |     1024 |
|  8 | 2022-01-01 07:00:00 |          27 |       23 |     1019 |
|  9 | 2022-01-01 08:00:00 |          28 |       22 |     1022 |
| 10 | 2022-01-01 09:00:00 |          10 |       33 |     1021 |
| 11 | 2022-01-01 10:00:00 |          15 |       55 |     1022 |
| 12 | 2022-01-01 11:00:00 |          31 |       66 |     1024 |
| 13 | 2022-01-01 12:00:00 |          25 |       33 |     1025 |
| 14 | 2022-01-01 13:00:00 |          33 |       22 |     1023 |
| 15 | 2022-01-01 14:00:00 |          34 |       22 |     1021 |
| 16 | 2022-01-01 15:00:00 |          35 |       33 |     1025 |
| 17 | 2022-01-01 16:00:00 |          11 |       33 |     1026 |
| 18 | 2022-01-01 17:00:00 |           0 |       55 |     1021 |
| 19 | 2022-01-01 18:00:00 |         -10 |       50 |     1012 |
| 20 | 2022-01-01 19:00:00 |          -5 |       65 |     1010 |
| 21 | 2022-01-01 20:00:00 |        -3.2 |       95 |     1023 |
+----+---------------------+-------------+----------+----------+
21 rows in set (0.00 sec)
```

### Using docker-compose

**Build image:**

```shell
docker-compose build
```

**Startup:**
```shell
docker-compose up -d
```

**Shutdown:**
```shell
docker-compose down
```

### Description

Builds a docker image with Ubuntu + MySQL and uploads compiled Go API into `/app` directory. Compiler is using ZSCaler Root Certificate (this is required for me, you can remove it completely). 

```bash
make build
```

Will start docker container with MySQL Database and GO Lang API running. 
You can edit `create_table.sql` and add example data if you want. 

```bash
make run 
```

Go API is listening on port `:8080`
Port `:8080` is exposed to port `:80` then. 

- `http://localhost/` will render graph from data collected.
- `http://localhost/writedata/` will write data to MySQL Database

You can test the api runnin `curl` command: 

```bash
curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d "time=2019-01-01 01:01:01&temperature=25.5&humidity=50.0&pressure=1005.5&param1=value1&param2=value2" http://SERVER_ADDRESS:SERVER_PORT/writedata

```

## Requirements
1. **ESP32** Dev Board
2. **BME280** Temperature sensor from Bosch
3. **MicroUSB Cable** / **Battery** or **USB 5V/1A** or other Source of energy
4. **Virtual** or **Cloud Machine** running **linux**
5. GoLang, VSCode, Arduino IDE
6. **Docker**, Podman or K8S for running Container.


## Introduction
At the beginning, I would like to mention here, that I'm not Arduiono specialist or programmer. I'm regular IT guy. I just wanted to try something new and different way. 

For basic understanding of things I've selected easy setup of Temperature measuring device. Since I have less knowledge of modern `APIs` or Containers I wanted to do it this way. 

## Rendering the Graph
GoLang API won't do much. Basically it can write data with `writedata` `func` to MySQL Database and It can draw simple line chart with `renderGraph` function. Both are pretty simple functions. In `renderGraph` function I've used Go-Echarts package to visualize the data. Visualization is basic and it is not perfect. It is important to say that some limits were reached already. Go-Echarts do not have all the functionality of Apache ECharts implemented, therefor there are some small glitches. 


![graph](./img/chart.png)

Data being sent from ESP32 device directly via `HTTP` `POST` requests. 


## ESP32 & BME280

Here is an example of how you can connect the BME280 sensor to the ESP32:

1. Connect the BME280 sensor to the ESP32 using the I2C interface. You will need to connect the SDA (data) and SCL (clock) lines of the sensor to the corresponding SDA and SCL pins on the ESP32. You will also need to connect the VCC and GND pins of the sensor to the appropriate power supply and ground pins on the ESP32.
2. Install the necessary libraries on your ESP32 board. You will need the "Adafruit BME280 Library" and the "Adafruit Unified Sensor Library" to use the BME280 sensor with the ESP32. You can install these libraries by going to Sketch > Include Library > Manage Libraries in the Arduino IDE and searching for the "Adafruit BME280" and "Adafruit Unified Sensor" libraries.
3. Initialize the BME280 sensor in your code by calling the begin() function of the Adafruit_BME280 class. You will need to pass the I2C address of the BME

![board](./img/board_esp32_bme280.jpeg)

## ESP 32 Example code could look like this

### C Language Code for ESP32 Device

This code uses the libcurl library to send an HTTP POST request to your GoLang server at the /writedata endpoint with the form data specified in the CURLOPT_POSTFIELDS option. The server will then handle the request and write the data to the database.

Keep in mind that you will need to replace "your-server-ip" with the actual IP address of your server. You can also modify the form data as needed to send different values for the time, temperature, humidity, and pressure.

```c
#include <stdio.h>
#include <curl/curl.h>

int main(void)
{
    CURL *curl;
    CURLcode res;

    curl = curl_easy_init();
    if(curl) {
        curl_easy_setopt(curl, CURLOPT_URL, "http://your-server-ip:8080/writedata");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, "time=2022-01-01 12:00:00&temperature=25&humidity=50&pressure=1013");

        res = curl_easy_perform(curl);
        if(res != CURLE_OK) {
            fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));
        }

        curl_easy_cleanup(curl);
    }

    return 0;
}

```

### Direct MySQL Insert

```c
#include <stdio.h>
#include <stdlib.h>
#include <mysql/mysql.h>
#include <bme280.h>

#define HOST "localhost"
#define USER "dbuser"
#define PASSWORD "heslo"
#define DATABASE "temperature_db"

int main(int argc, char *argv[]) {
  // Initialize the BME280 sensor
  bme280_init();

  // Connect to the MySQL database
  MYSQL *conn = mysql_init(NULL);
  if (!mysql_real_connect(conn, HOST, USER, PASSWORD, DATABASE, 0, NULL, 0)) {
    fprintf(stderr, "%s\n", mysql_error(conn));
    return 1;
  }

  // Read data from the BME280 sensor
  float temperature = bme280_read_temperature();
  float humidity = bme280_read_humidity();
  float pressure = bme280_read_pressure();

  // Build the SQL query
  char query[256];
  sprintf(query, "INSERT INTO data (temperature, humidity, pressure) VALUES (%f, %f, %f)", temperature, humidity, pressure);

  // Execute the query
  if (mysql_query(conn, query)) {
    fprintf(stderr, "%s\n", mysql_error(conn));
    return 1;
  }

  // Close the connection
  mysql_close(conn);

  return 0;
}
```

## Actual ESP32 C Language code used
![C langguage Code](./ESP32/esp32_go_api.c?raw=true)

It uses a number of libraries, including `WiFi`, `WiFiUdp`, `time`, `WebServer`, `Wire`, `Adafruit_Sensor`, `Adafruit_BME280`, and `HTTPClient` to connect to a WiFi network, retrieve the current time from a time server, read sensor data from a `BME280` sensor (a sensor that measures temperature, humidity, and pressure), and post that data to an Go `API server`.

First, it initializes the `BME280` sensor and connects to the WiFi network using the credentials provided in the WIFI_SSID and WIFI_PASSWORD constants. Then it uses the WiFiUDP library to send a packet to a time server to retrieve the current time. The current time is then used to set the time on the device.

In the loop, it repeatedly `checks the WiFi connection` and if the connection is lost it tries to reconnect. Then it reads the `temperature`, `humidity`, and `pressure` data from the `BME280 sensor` and prints the data to the serial monitor along with the `current time`. Finally, the code then post this data to the Go `api server` with an `http client`.

Overall, this code is designed to retrieve sensor data from a BME280 sensor, get the current time, and then post that data to an API server via a wifi connection.


## Few suggestions

1. Make sure that your GoLang API is running and accessible from the ESP32. You can test this by trying to access the API from your web browser or by using a CURL command like the one I provided in my previous message.
2. Make sure that the WiFi connection on the ESP32 is stable and that it can reach the GoLang API server. You can check the status of the WiFi connection by looking at the serial output from the ESP32.
3. Make sure that the sendDataToAPI function is being called correctly. You can check this by adding some debug statements in the function to print out the values of the temperature, humidity, and pressure variables.
4. If you are still having trouble, you may want to check the error messages that are being printed by the GoLang API. This can help you identify any problems with the API itself or with the data that is being sent to the API.

