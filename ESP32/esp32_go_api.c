#include <dummy.h>

#include <WiFi.h>
#include <WiFiUdp.h>
#include <time.h>

#include <WebServer.h>
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>
#include <HTTPClient.h>


// WIFI Settings
#define WIFI_SSID "YOUR-WIFI-SSID"
#define WIFI_PASSWORD "YOUR-PASSWORD"
#define BME280_ADDRESS 0x76
#define LED_PIN 2


// Setup Time Server 
// Time server IP address and port
IPAddress timeServerIP(216, 239, 35, 0);
const int timeServerPort = 123;
#define NTP_PACKET_SIZE 48
byte packetBuffer[NTP_PACKET_SIZE];

// Declare variables for storing current time
time_t currentTime;
struct tm *timeinfo;
char buffer[80];

// Possible DB Server configuration
//#define DB_SERVER "mysql"
//#define DB_PORT 3306

// API Server Address and Port
#define SERVER_ADDRESS "192.168.100.100"
#define SERVER_PORT 80

// ALTITUDE
#define ALTITUDE 218.4 // Altitude for BRNO Heyrovskeho in Meters

Adafruit_BME280 bme;
WiFiUDP udp;
 
// Declare variables for storing temperature, humidity, and pressure data
float temperature, humidity, pressure;
 
void setup() {
  // Serial Port
  Serial.begin(9600);

  // Init led
  pinMode(LED_PIN, OUTPUT);
  delay(100);

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

  // Connect to the time server
  udp.begin(timeServerPort);
  sendNTPpacket(timeServerIP);
  delay(2000);

  // Get the current time from the time server
  // time_t currentTime;
  if (getCurrentTime(currentTime)) {
    // Set the current time
    setTime(currentTime);

    // Print the current time to the serial monitor
    struct tm *timeinfo;
    timeinfo = localtime(&currentTime);
    char buffer[80];
    strftime(buffer, 80, "%A, %B %d %Y %H:%M:%S", timeinfo);
    Serial.println(buffer);
  }

}

void loop() {

   if (WiFi.status() != WL_CONNECTED) {
    // WiFi connection is lost, try to reconnect
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  }

  // Read temperature, humidity, and pressure data from the BME280 sensor
  temperature = bme.readTemperature();
  humidity = bme.readHumidity();
  pressure = bme.readPressure() / 100.0F; // hPa

  // Calculate the absolute air pressure at a given altitude
  float seaLevelPressure = bme.seaLevelForAltitude(ALTITUDE, pressure); // hPa
 
  // Print the data to the serial monitor
  Serial.println("----------------------------------");
  Serial.print("Current time is : ");
  Serial.println(getTimeString());
  Serial.println("----------------------------------");
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


// FORMAT DATETIME FOR API CALL
// This will return the time in the correct format, "YYYY-MM-DDTHH:MM:SSZ".
// Fixed function to return currentTime.
String getTimeString() {
  time_t now = time(nullptr);
  struct tm *timeinfo;
  timeinfo = localtime(&now);
  char buffer[80];
  strftime(buffer, 80, "%Y-%m-%d %H:%M:%S", timeinfo);
  return String(buffer);
}


// Sets actuall time CET Timezone by default
void setTime(time_t currentTime)
{
  struct tm timeinfo;
  gmtime_r(&currentTime, &timeinfo);
  configTime(0, 0, "pool.ntp.org", "time.nist.gov");
  setenv("TZ", "CET-1CEST,M3.5.0/2,M10.5.0/3", 1);
  tzset();
  timeinfo.tm_isdst = -1;
  currentTime = mktime(&timeinfo);
  Serial.printf("Setting time using settimeofday() to: %s", asctime(&timeinfo));
  timeval tv = { currentTime, 0 };
  settimeofday(&tv, nullptr);
}

// Send an NTP request to the time server
void sendNTPpacket(IPAddress& address) {
  // Set all bytes in the buffer to 0
  memset(packetBuffer, 0, NTP_PACKET_SIZE);
 
  // Initialize values needed to form NTP request
  packetBuffer[0] = 0b11100011;   // LI, Version, Mode
  packetBuffer[1] = 0;     // Stratum, or type of clock
  packetBuffer[2] = 6;     // Polling Interval
  packetBuffer[3] = 0xEC;  // Peer Clock Precision
 
  // 8 bytes of zero for Root Delay & Root Dispersion
  packetBuffer[12]  = 49;
  packetBuffer[13]  = 0x4E;
  packetBuffer[14]  = 49;
  packetBuffer[15]  = 52;
 
  // Send a packet requesting a timestamp:
  udp.beginPacket(address, 123); // NTP requests are to port 123
  udp.write(packetBuffer, NTP_PACKET_SIZE);
  udp.endPacket();
}

// Parse the current time from the NTP response
bool getCurrentTime(time_t& currentTime) {
  int cb = udp.parsePacket();
  if (!cb) {
    Serial.println("No NTP response :-(");
    return false;
  }
 
  // Read the NTP response packet
  udp.read(packetBuffer, NTP_PACKET_SIZE);
 
  // The first 48 bytes of the packet contain the NTP header data
  // The header consists of 10 words (4 bytes each)
  // The first word is a bitfield indicating the version and mode of the packet
  // The next two words are the timestamp in seconds
  // The last 7 words are unused in this implementation
 
  // Extract the timestamp from the packet
  unsigned long highWord = word(packetBuffer[40], packetBuffer[41]);
  unsigned long lowWord = word(packetBuffer[42], packetBuffer[43]);
  unsigned long secsSince1900 = highWord << 16 | lowWord;
 
  // The timestamp value is the number of seconds since January 1, 1900
  // Subtract 70 years (in seconds) to convert the timestamp to the number of seconds since January 1, 1970
  const unsigned long seventyYears = 2208988800UL;
  currentTime = secsSince1900 - seventyYears;
 
  // Print the current time
  Serial.print("Current time: ");
  Serial.println(ctime(&currentTime));
 
  return true;
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
