package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	csvFilePath  = "ohlcv.csv"
	htmlFilePath = "ohlcv.html"

	openLabel  = "open"
	closeLabel = "close"
	lowLabel   = "low"
	highLabel  = "high"
)

func main() {
	records, loadErr := loadCsv(csvFilePath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	xAxe, yAxe := prepareLineData(records)

	line := plotChart(xAxe, yAxe)

	pageErr := createHtml(htmlFilePath, line)
	if pageErr != nil {
		log.Fatal(pageErr)
	}
}

func createHtml(filePath string, charts ...components.Charter) error {
	page := components.NewPage()
	page.AddCharts(charts...)

	file, createErr := os.Create(filePath)
	if createErr != nil {
		return createErr
	}
	return page.Render(io.MultiWriter(file))
}

func plotChart(xAxe []string, yAxe map[string][]opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | BTC-USDT",
			Subtitle: "OHLCV of 2022-01-01",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "go-echarts line example",
			Theme:     "dark",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:         true,
			SelectedMode: "multiple",
		}),
		// AXIS
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
	)

	line.SetXAxis(xAxe)

	line.AddSeries(openLabel+"_smooth", yAxe[openLabel]).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(closeLabel+"_smooth", yAxe[closeLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(lowLabel, yAxe[lowLabel])
	line.AddSeries(highLabel, yAxe[highLabel])

	return line
}

func prepareLineData(records [][]string) ([]string, map[string][]opts.LineData) {
	start := 0
	if strings.Contains(records[0][0], "CLOSED_AT") {
		start = 1
	}
	x := make([]string, 0)
	y := make(map[string][]opts.LineData, 0)
	for _, record := range records[start:] {
		// CLOSED_AT,OPENED_AT,OPEN,HIGH,LOW,CLOSE,VOLUME,COMPONENT,BUCKET
		openVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		closeVal, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}
		lowVal, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}
		highVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			continue
		}
		x = append(x, record[1])
		y[openLabel] = append(y[openLabel], opts.LineData{Value: openVal})
		y[closeLabel] = append(y[closeLabel], opts.LineData{Value: closeVal})
		y[lowLabel] = append(y[lowLabel], opts.LineData{Value: lowVal})
		y[highLabel] = append(y[highLabel], opts.LineData{Value: highVal})
	}

	return x, y
}

func loadCsv(filePath string) ([][]string, error) {
	file, openErr := os.Open(filePath)
	if openErr != nil {
		return nil, openErr
	}
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, readErr := reader.ReadAll()
	if readErr != nil {
		return nil, readErr
	}
	return records, nil
}
