package bleclient

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

var bleAdaptor *gobot.Adaptor

func Connect(address string) {
	// Create a new BLE Adaptor
	bleAdaptor = ble.NewClientAdaptor(address)

	// Connect to the ESP32
	bleAdaptor.Connect()
}

func ReadTemperature(uuid string) ([]byte, error) {
	// Read the value of the temperature characteristic
	return bleAdaptor.ReadCharacteristic(uuid)
}

func ReadPressure(uuid string) ([]byte, error) {
	// Read the value of the pressure characteristic
	return bleAdaptor.ReadCharacteristic(uuid)
}

func ReadHumidity(uuid string) ([]byte, error) {
	// Read the value of the humidity characteristic
	return bleAdaptor.ReadCharacteristic(uuid)
}

func Disconnect() {
	// Disconnect from the ESP32
	bleAdaptor.Disconnect()
}
