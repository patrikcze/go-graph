# PoC 
Trying to mess around with GoLang, prepare simple APi which will collect data from ESP32 + BME280 Module, uploads data to MySQL Database and will draw a line chart.

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
## Few suggestions

1. Make sure that your GoLang API is running and accessible from the ESP32. You can test this by trying to access the API from your web browser or by using a CURL command like the one I provided in my previous message.
2. Make sure that the WiFi connection on the ESP32 is stable and that it can reach the GoLang API server. You can check the status of the WiFi connection by looking at the serial output from the ESP32.
3. Make sure that the sendDataToAPI function is being called correctly. You can check this by adding some debug statements in the function to print out the values of the temperature, humidity, and pressure variables.
4. If you are still having trouble, you may want to check the error messages that are being printed by the GoLang API. This can help you identify any problems with the API itself or with the data that is being sent to the API.
