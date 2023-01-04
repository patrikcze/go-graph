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
	var times = []time.Time{}
	// Connect to the database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME+"?parseTime=true")
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the data from the database
	rows, err := db.Query("SELECT time, temperature FROM data")
	if err != nil {
		http.Error(w, "Error querying data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the data and add points to the series
	for rows.Next() {
		var t time.Time
		var temp float64
		err := rows.Scan(&t, &temp)
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

	}

	// Create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Graf teplot",
			Subtitle: "Čárový graf teplot vykreslený.",
		}))

	// Put data into instance
	line.SetXAxis(times).
		AddSeries("Teploty", temperatures).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)

}

/*
func chartHandler(c http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME+"?parseTime=true")
	if err != nil {
		http.Error(c, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the data from the database
	rows, err := db.Query("SELECT time, temperature FROM data")
	if err != nil {
		http.Error(c, "Error querying data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a line series for the data
	temperatureSeries := chart.TimeSeries{
		Name: "Temperature",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
		},
	}

	// Iterate through the data and add points to the series
	for rows.Next() {
		var t time.Time
		var temp float64
		err := rows.Scan(&t, &temp)
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
		temperatureSeries.XValues = append(temperatureSeries.XValues, t)
		temperatureSeries.YValues = append(temperatureSeries.YValues, temp)

	}
	// Log what has been added to the series
	log.Printf("Following time values will be added : %v", temperatureSeries.XValues)
	log.Printf("With following temperature values : %v", temperatureSeries.YValues)

	// Create the chart with the series
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				if t, ok := v.(time.Time); ok {
					return t.Format("2006-01-02 15:04:05")
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: 100,
			},
		},
		Series: []chart.Series{
			temperatureSeries,
		},
	}

	// Render the chart as a PNG image
	buf := bytes.NewBuffer([]byte{})
	// Setup size of chart
	graph.Width = 512
	graph.Height = 512
	err = graph.Render(chart.PNG, buf)
	if err != nil {
		log.Printf("Error rendering chart: %v", err)
		http.Error(c, "Error rendering chart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the chart to the HTTP response
	c.Header().Set("Content-Type", "image/png")
	c.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	_, err = c.Write(buf.Bytes())
	if err != nil {
		log.Printf("Error writing chart to response: %v", err)
		http.Error(c, "Error writing chart to response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
*/
