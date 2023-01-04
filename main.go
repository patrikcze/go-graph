package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

const (
	DB_USER     = "dbuser"
	DB_PASSWORD = "heslo"
	DB_NAME     = "temperature_db"
)

func main() {
	/*r := mux.NewRouter()
	r.HandleFunc("/graph", chartHandler)
	http.ListenAndServe(":8080", r)
	*/
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8080", nil)
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// Reset Items
	temperatures := make([]opts.LineData, 0)
	humidities := make([]opts.LineData, 0)
	preasures := make([]opts.LineData, 0)
	var times = []time.Time{}
	// Connect to the database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME+"?parseTime=true")
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the data from the database
	rows, err := db.Query("SELECT time, temperature, humidity, preasure FROM data")
	if err != nil {
		http.Error(w, "Error querying data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the data and add points to the series
	for rows.Next() {
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
		preasures = append(preasures, opts.LineData{Value: pres})

	}

	// Create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Graf teplot",
			Subtitle: "Čárový graf teplot vykreslený.",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:    true,
			Bottom:  "20",
			Padding: 5,
		}),
	)

	// Put data into instance
	line.SetXAxis(times).
		AddSeries("Teploty", temperatures).
		AddSeries("Vlhkosti", humidities).
		AddSeries("Tlaky", preasures).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)

}
