package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/", renderGraph)
	http.HandleFunc("/writedata", writeData)

	//http.ListenAndServe(":8080", nil)
	// Check error
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME+"?parseTime=true")
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

	// Return a success response
	w.WriteHeader(http.StatusOK)
	log.Println("Data wuccessfully written to database.")
	w.Write([]byte("Data written to database successfully"))
	if err != nil {
		http.Error(w, "Error writting data: "+err.Error(), http.StatusBadRequest)
		return
	}
}

// Render the chart
func renderGraph(w http.ResponseWriter, _ *http.Request) {
	// Reset Items
	temperatures := make([]opts.LineData, 0)
	humidities := make([]opts.LineData, 0)
	pressures := make([]opts.LineData, 0)
	var times []time.Time
	// Connect to the database
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME+"?parseTime=true")
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
		charts.WithInitializationOpts(opts.Initialization{
			Theme:     types.ThemeWesteros,
			PageTitle: "Grafík",
			Height:    "768px",
			Width:     "1024px",
		}),
		// Name of the chart and subtitle
		charts.WithTitleOpts(opts.Title{
			Title:    "Graf Teplot (ESP32 & BME280)",
			Subtitle: "Pokusí se vykreslit data z databáze. Teploty, Vlhkosti a tlaky.",
		}),
		// Shows tool tip on click
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "click",
		}),
		// Will try to render Legend (you can click on each series)
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Bottom: "50%",
			Align:  "right",
			Left:   "90%",
			Right:  "10%",
			Top:    "50%",
			Orient: "vertical",
		}),
		// This will setup DataZoom Slider in the chart
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   100,
		}),
		// This will add Toolbox to top right corner and allows to export to PNG or Show dataset.
		charts.WithToolboxOpts(opts.Toolbox{
			Show: true,
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Name:  "Heyrovskeho5",
					Title: "Uložit",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
				},
				DataZoom: &opts.ToolBoxFeatureDataZoom{
					Show: true,
				},
				Restore: &opts.ToolBoxFeatureRestore{
					Show:  true,
					Title: "Refresh",
				},
			},
			Top: "",
		}),
		// Some basic options for X Axis
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Datum a čas",
			Show: true,
			//Max:  "dataMax",
			//Type: "category",
			//Data: times,
			//Width: "50%",
			//Type: "time",
			// Options for X Axis Labels (GRID is not supported in Go-ECharts package!!!)
			AxisLabel: &opts.AxisLabel{
				Show: true,
				//Interval:  "10",
				Inside: false,
				//Rotate: 90,
				Margin: 0,
				//Formatter: opts.FuncOpts(fn),
				//Formatter: "{HH}:{mm}",
				Align: "",
				//VerticalAlign: "right",
				//LineHeight: "250",
			}},
		),
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
	line.Render(w)
	if err != nil {
		http.Error(w, "Error rendering the chart : "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error in rendering the chart : %v", err)
	}
}
