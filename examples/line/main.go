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

	lineChartA := plotLineA(xAxe, yAxe)
	lineChartB := plotLineB(xAxe, yAxe)

	pageErr := createHtml(htmlFilePath, lineChartA, lineChartB)
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

func plotLineA(xAxe []string, yAxe map[string][]opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | BTC-USDT",
			Subtitle: "OHLCV of 2022-01-01",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Anything you want",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "DataView",
					Lang:  []string{"data view", "turn off", "refresh"},
				},
			}},
		),
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
			SplitArea: &opts.SplitArea{
				Show: true,
			},
		}),

		charts.WithVisualMapOpts(
		//opts.VisualMap{Type: "piecewise", Min: 0.0, Max: 46800.0, InRange: &opts.VisualMapInRange{Color: []string{"#93CE07"}}},      // green
		//opts.VisualMap{Type: "piecewise", Min: 46800.1, Max: 47200.00, InRange: &opts.VisualMapInRange{Color: []string{"#FBDB0F"}}}, // yellow
		//opts.VisualMap{Type: "piecewise", Min: 47200.1, Max: 47400.00, InRange: &opts.VisualMapInRange{Color: []string{"#FC7D02"}}}, // orange
		//opts.VisualMap{Type: "piecewise", Min: 47400.1, Max: 47600.00, InRange: &opts.VisualMapInRange{Color: []string{"#FD0100"}}}, // red
		//opts.VisualMap{Type: "piecewise", Min: 47600.1, Max: 48000.00, InRange: &opts.VisualMapInRange{Color: []string{"#AA069F"}}}, // purple

		//opts.VisualMap{Type: "piecewise", Range: []float32{47501.0, 50000.0}, InRange: &opts.VisualMapInRange{Color: []string{"#AA069F"}}}, // purple
		),

		charts.WithColorsOpts(opts.Colors{"green", "blue", "pink", "orange"}),
	)

	line.SetXAxis(xAxe)

	line.AddSeries(openLabel, yAxe[openLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(closeLabel, yAxe[closeLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.AddSeries(lowLabel, yAxe[lowLabel],
		charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "red"}), // instead of pink
	)
	line.AddSeries(highLabel, yAxe[highLabel], charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}

func plotLineB(xAxe []string, yAxe map[string][]opts.LineData) *charts.Line {
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
		// AXIS
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
			SplitArea: &opts.SplitArea{
				Show: true,
			},
		}),
	)

	line.SetXAxis(xAxe)

	line.AddSeries(openLabel+"_smooth", yAxe[openLabel]).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: true}),
			charts.WithMarkLineNameTypeItemOpts(
				opts.MarkLineNameTypeItem{Type: "max"},
				opts.MarkLineNameTypeItem{Type: "min"},
				opts.MarkLineNameTypeItem{Type: "average"},
			),
			charts.WithMarkPointStyleOpts(
				//opts.MarkPointStyle{Symbol: []string{"circle", "pin", "arrow"}},
				opts.MarkPointStyle{Symbol: []string{"pin"}},
			),
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Type: "max"},
				opts.MarkPointNameTypeItem{Type: "min"},
				opts.MarkPointNameTypeItem{Type: "average"},
			),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color: "orange",
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Color:   "green",
				Opacity: 0.4,
			}),
		)

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
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	records, readErr := reader.ReadAll()
	if readErr != nil {
		return nil, readErr
	}
	return records, nil
}
