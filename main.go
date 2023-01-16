// Package main is the entry point of the application.
// It contains the main function that runs the application.
// / default API will render chart
// writedata will write data into MySQL Database
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

/*
const (
	DB_USER     = "dbuser"
	DB_PASSWORD = "heslo"
	DB_NAME     = "temperature_db"
)
*/

// var dbServer string
var dbUser string
var dbPassword string
var dbName string

// Config holds the configuration options for the chart.
// config.json structure
// used to configure look and feel of the chart
type Config struct {
	Initialization opts.Initialization `json:"initialization"`
	Title          opts.Title          `json:"title"`
	Tooltip        opts.Tooltip        `json:"tooltip"`
	Legend         opts.Legend         `json:"legend"`
	DataZoom       opts.DataZoom       `json:"dataZoom"`
	Toolbox        opts.Toolbox        `json:"toolbox"`
	XAxis          opts.XAxis          `json:"xAxis"`
}

func main() {

	// Read and parse the config.json
	config := Config{}
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	// Read variables from ENV
	// dbUser is default MySQL User name
	dbUser = os.Getenv("MYSQL_USER")
	// dbPassword default MySQL Password
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	// dbName is Default MySQL Database name
	dbName = os.Getenv("MYSQL_DATABASE")

	/*r := mux.NewRouter()
	r.HandleFunc("/graph", chartHandler)
	http.ListenAndServe(":8080", r)
	*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderGraph(w, r, config)
	})
	http.HandleFunc("/writedata", writeData)
	// Check error
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// writeData handles HTTP requests to write data to the MySQL database.
// It expects the following form values in the request body:
// - "time": a string in the format "2006-01-02 15:04:05"
// - "temperature": a float value representing the temperature
// - "humidity": a float value representing the humidity
// - "pressure": a float value representing the pressure
// If the form values are not in the expected format or there is a problem connecting to the database,
// it will return an appropriate error status code and message.
// Update (12.1.2023 - Using TCP connection for MySQL using "db" name of container.)
func writeData(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the form values
	timeStr := r.FormValue("time")
	tempStr := r.FormValue("temperature")
	humStr := r.FormValue("humidity")
	presStr := r.FormValue("pressure")

	// Convert the form values to the appropriate types
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Printf("Error parsing time value : %v", err)
		http.Error(w, "Error parsing time value: "+err.Error(), http.StatusBadRequest)
		return
	}
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		log.Printf("Error parsing temperature value : %v", err)
		http.Error(w, "Error parsing temperature value: "+err.Error(), http.StatusBadRequest)
		return
	}
	hum, err := strconv.ParseFloat(humStr, 64)
	if err != nil {
		log.Printf("Error parsing humidity value : %v", err)
		http.Error(w, "Error parsing humidity value: "+err.Error(), http.StatusBadRequest)
		return
	}
	pres, err := strconv.ParseFloat(presStr, 64)
	if err != nil {
		log.Printf("Error parsing preassure value : %v", err)
		http.Error(w, "Error parsing pressure value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(db)/"+dbName+"?parseTime=true")
	if err != nil {
		log.Printf("There was problem with connection to databsae : %v", err)
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute the INSERT statement
	_, err = db.Exec("INSERT INTO data (time, temperature, humidity, pressure) VALUES (?, ?, ?, ?)", t, temp, hum, pres)
	if err != nil {
		log.Printf("Error writing data to database : %v", err)
		http.Error(w, "Error writing data to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fixed error handling here (hopefully)
	n, err := w.Write([]byte("Data written to database successfully"))
	if n != 200 {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}
	// Return a success response
	w.WriteHeader(http.StatusOK)
	log.Println("Data wuccessfully written to database.")
}

// renderGraph handles HTTP requests to render a chart of the data from the MySQL database.
// It queries the data from the "data" table in the database and pass it to the go-echarts package
// to generate the chart.
// It takes no parameters and returns a rendered chart as an HTTP response.
// Render the chart (12.1.2023 - tpc connection to MySQL using "db" name of container.)
func renderGraph(w http.ResponseWriter, _ *http.Request, config Config) {
	// Reset Items
	temperatures := make([]opts.LineData, 0)
	humidities := make([]opts.LineData, 0)
	pressures := make([]opts.LineData, 0)
	var times []time.Time
	// Connect to the database
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(db)/"+dbName+"?parseTime=true")
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the data from the database
	rows, err := db.Query("SELECT time, temperature, humidity, pressure FROM data")
	if err != nil {
		http.Error(w, "Error querying data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the data and add points to the series
	for rows.Next() {
		// Create slice with times
		var t time.Time
		var temp float64
		var hum float64
		var pres float64
		err := rows.Scan(&t, &temp, &hum, &pres)
		if err != nil {
			log.Fatal(err)
		}
		// Convert the time value to a string
		tString := t.Format("2006-01-02 15:04:05")

		// Parse the time value from the database
		t, err = time.Parse("2006-01-02 15:04:05", tString)
		if err != nil {
			log.Fatal(err)
		}
		// Append the time and temperature values to the chart data
		times = append(times, t)
		temperatures = append(temperatures, opts.LineData{Value: temp})
		humidities = append(humidities, opts.LineData{Value: hum})
		pressures = append(pressures, opts.LineData{Value: pres})

	}
	/*
		fn := `function(value){
			        let label;
			        if (value.getMinutes() < 10){
			          label = value.getHours() + ":0" +value.getMinutes();
			        }
			        else {
			          label = value.getHours() + ":" +value.getMinutes();
			        }
			        return label;
			      }`
	*/

	// Create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		// Initial option of Chart
		charts.WithInitializationOpts(config.Initialization),
		// Name of the chart and subtitle
		charts.WithTitleOpts(config.Title),
		// Shows tool tip on click
		charts.WithTooltipOpts(config.Tooltip),
		// Will try to render Legend (you can click on each series)
		charts.WithLegendOpts(config.Legend),
		// This will setup DataZoom Slider in the chart
		charts.WithDataZoomOpts(config.DataZoom),
		// This will add Toolbox to top right corner and allows to export to PNG or Show dataset.
		charts.WithToolboxOpts(config.Toolbox),
		// Some basic options for X Axis
		charts.WithXAxisOpts(config.XAxis),
	)

	// Puts data into instance and setup some further options to each serie.
	line.SetXAxis(times).
		AddSeries("Teploty (℃)", temperatures, charts.WithLabelOpts(opts.Label{Show: true})).
		AddSeries("Vlhkosti (%)", humidities, charts.WithLabelOpts(opts.Label{Show: true})).
		AddSeries("Tlaky (hPa)", pressures, charts.WithLabelOpts(opts.Label{Show: true})).
		SetSeriesOptions(charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: true,
			}),
		)
	// render if no issue
	e := line.Render(w)
	if e != nil {
		http.Error(w, "Error rendering the chart : "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error in rendering the chart : %v", e)
		return
	}

}
