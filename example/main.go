package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// Data represents a temperature and humidity reading from the ESP32 device.
type Data struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Time        string  `json:"time"`
}

// API represents the API key and the database connection information.
type API struct {
	Key    string
	Db     *sql.DB
	HTTP   *http.Client
	Router *mux.Router
}

// NewAPI initializes the API.
func NewAPI(key string) (*API, error) {
	api := &API{
		Key:    key,
		Router: mux.NewRouter(),
	}

	// Connect to the database.
	db, err := sql.Open("mysql", "root:password@/database_name?parseTime=true")
	if err != nil {
		return nil, err
	}
	api.Db = db

	// Add the routes to the router.
	api.Router.HandleFunc("/data", api.handleData).Methods("POST")
	api.Router.HandleFunc("/chart", api.handleChart).Methods("GET")

	return api, nil
}

// handleData handles data POST requests.
func (api *API) handleData(w http.ResponseWriter, r *http.Request) {
	// Verify the API key.
	if r.Header.Get("x-api-key") != api.Key {
		http.Error(w, "invalid API key", http.StatusForbidden)
		return
	}

	// Decode the request body into a Data struct.
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the data into the database.
	result, err := api.Db.Exec("INSERT INTO data (temperature, humidity, time) VALUES (?, ?, ?)", data.Temperature, data.Humidity, data.Time)
	if err != nil {
		http.Error(w, "error inserting data into database", http.StatusInternalServerError)
		return
	}

	// Return the number of rows affected.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "error getting number of rows affected", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d rows affected\n", rowsAffected)
}
