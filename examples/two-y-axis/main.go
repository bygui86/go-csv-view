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
	ohlcvFilePath  = "ohlcv.csv"
	tradesFilePath = "trades.csv"
	htmlFilePath   = "two-y-axis.html"

	openLabel  = "open"
	closeLabel = "close"
	lowLabel   = "low"
	highLabel  = "high"
	sizeLabel  = "size"
)

func main() {
	ohlcvRecords, ohlcvErr := loadCsv(ohlcvFilePath)
	if ohlcvErr != nil {
		log.Fatal(ohlcvErr)
	}

	ohlcXaxe, ohlcYaxe := prepareOhlcData(ohlcvRecords)

	tradesRecords, tradesErr := loadCsv(tradesFilePath)
	if tradesErr != nil {
		log.Fatal(tradesErr)
	}

	_, tradesYaxe := prepareTradesData(tradesRecords)

	line := plotChart(ohlcXaxe, ohlcYaxe, tradesYaxe)

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

func plotChart(ohlcXaxe []string, ohlcYaxe map[string][]opts.LineData, tradesYaxe []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | BTC-USDT",
			Subtitle: "OHLCV of 2022-01-01",
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
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
				Snap: true,
			},
		}),
		// AXIS
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:  "Price",
			Type:  "value",
			Show:  true,
			Scale: true,
			//GridIndex: 0, // y index 0 // not required
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Name:  "Size",
		Type:  "value",
		Show:  true,
		Scale: true,
		//GridIndex: 1, // y index 1 // not required
	})

	line.SetXAxis(ohlcXaxe)

	line.AddSeries(openLabel, ohlcYaxe[openLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	//line.AddSeries(openLabel, ohlcYaxe[openLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true, YAxisIndex: 0})) // YAxisIndex not required if referring to index 0
	line.AddSeries(closeLabel, ohlcYaxe[closeLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(lowLabel, ohlcYaxe[lowLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(highLabel, ohlcYaxe[highLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	line.AddSeries(sizeLabel, tradesYaxe, charts.WithLineChartOpts(opts.LineChart{Smooth: true, YAxisIndex: 1}))

	return line
}

func prepareTradesData(records [][]string) ([]string, []opts.LineData) {
	start := 0
	if strings.Contains(records[0][0], "TIMESTAMP") {
		start = 1
	}
	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	for _, record := range records[start:] {
		// TIMESTAMP,TRADE_ID,PRICE,SIDE,SIZE,BUYER_ORDER_ID,SELLER_ORDER_ID,COMPONENT,BUCKET
		sizeVal, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}
		x = append(x, record[1])
		y = append(y, opts.LineData{Value: sizeVal, YAxisIndex: 1})
	}

	return x, y
}

func prepareOhlcData(records [][]string) ([]string, map[string][]opts.LineData) {
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
		//y[openLabel] = append(y[openLabel], opts.LineData{Value: openVal, YAxisIndex: 0}) // YAxisIndex not required if referring to index 0
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
