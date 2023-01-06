#include <WiFi.h>
#include <WebServer.h>
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>
#include <HTTPClient.h>

// BME280 init
#define BME280_ADDRESS 0x76

// WIFI Settings
#define WIFI_SSID "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_PASSWORD"
#define BME280_ADDRESS 0x76

// Possible DB Server configuration
//#define DB_SERVER "mysql"
//#define DB_PORT 3306

// API Server Address and Port
#define SERVER_ADDRESS "YOUR_API_SERVER"
#define SERVER_PORT 80

// ALTITUDE
#define ALTITUDE 218.4 // Altitude for BRNO in Meters

Adafruit_BME280 bme;
 
// Declare variables for storing temperature, humidity, and pressure data
float temperature, humidity, pressure;
 
void setup() {
  // Initialize BME280 sensor
  if (!bme.begin(BME280_ADDRESS)) {
    Serial.println("Could not find a valid BME280 sensor, check wiring!");
    while (1);
  }
 
  // Initialize ESP32 WiFi
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
 
  // Print the IP address of the ESP32
  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}
 
void loop() {
  // Read temperature, humidity, and pressure data from the BME280 sensor
  temperature = bme.readTemperature();
  humidity = bme.readHumidity();
  pressure = bme.readPressure() / 100.0F; // hPa

  // Calculate the absolute air pressure at a given altitude
  float seaLevelPressure = bme.seaLevelForAltitude(ALTITUDE, pressure); // hPa
 
  // Print the data to the serial monitor
  Serial.print("Temperature: ");
  Serial.print(temperature);
  Serial.println(" *C");
  Serial.print("Humidity: ");
  Serial.print(humidity);
  Serial.println(" %");
  Serial.print("Pressure: ");
  Serial.print(pressure);
  Serial.println(" hPa");
  Serial.print("Sea level pressure (at ");
  Serial.print(ALTITUDE);
  Serial.print("m altitude): ");
  Serial.print(seaLevelPressure);
  Serial.println(" hPa");
 
  // Send the data to the GoLang API
  sendDataToAPI(temperature, humidity, seaLevelPressure);
 
  // Wait for a certain amount of time before collecting new data
  delay(60000); // Collect new data every 60 seconds
}

String getTimeString() {
  time_t now = time(nullptr);
  struct tm *timeinfo;
  timeinfo = localtime(&now);
  char buffer[80];
  strftime(buffer, 80, "%Y-%m-%d %H:%M:%S", timeinfo);
  return String(buffer);
}

// Function for sending data to the GoLang API
void sendDataToAPI(float temperature, float humidity, float pressure) {
  // Create an HTTP client
  HTTPClient http;

  // Set the API endpoint URL
  String apiEndpoint = "http://" + String(SERVER_ADDRESS) + ":" + String(SERVER_PORT) + "/writedata";
  http.begin(apiEndpoint);

  // Set the request body
  String body = "time=" + getTimeString() + "&temperature=" + String(temperature) + "&humidity=" + String(humidity) + "&pressure=" + String(pressure);

  // Set the request headers
  http.addHeader("Content-Type", "application/x-www-form-urlencoded");
  http.addHeader("Content-Length", String(body.length()));

  // Make the POST request
  int httpCode = http.POST(body);

  // Check the response status code
  if (httpCode > 0) {
    // HTTP request was successful
    String response = http.getString();
    Serial.println(httpCode);
    Serial.println(response);
  } else {
    // HTTP request failed
    Serial.println("Error: " + http.errorToString(httpCode));
  }

  // Close the HTTP client
  http.end();
}
