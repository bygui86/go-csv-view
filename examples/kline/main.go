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
)

func main() {
	records, loadErr := loadCsv(csvFilePath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	xAxe, ohlcYaxe, volYaxe := prepareOhlcvData(records)

	simpleChart := plotSimpleChart(xAxe, ohlcYaxe)
	overlapChart := plotOverlapChart(xAxe, ohlcYaxe, volYaxe)
	//bundleChart := plotBundleChart(xAxe, ohlcYaxe, volYaxe) // TODO see https://github.com/go-echarts/go-echarts/issues/169#issuecomment-1014285680

	volumeLineChart := plotVolumeLine(xAxe, volYaxe)

	//pageErr := createHtml(htmlFilePath, simpleChart)
	//pageErr := createHtml(htmlFilePath, overlapChart)
	//pageErr := createHtml(htmlFilePath, volumeLineChart)
	//pageErr := createHtml(htmlFilePath, simpleChart, volumeLineChart)
	//pageErr := createHtml(htmlFilePath, overlapChart, volumeLineChart)
	pageErr := createHtml(htmlFilePath, overlapChart, simpleChart, volumeLineChart)
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

// TODO try to specify multiple YAxis
// TODO convert to line + bars and extract in another example
func plotOverlapChart(xAxe []string, ohlcYAxe []opts.KlineData, volYaxe []opts.LineData) *charts.Kline {
	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | OHLCV | BTC-USDT | 2022-01-01",
			Subtitle: "OHLCV full",
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

	kline.ExtendYAxis(opts.YAxis{
		Name:  "Volume",
		Type:  "value",
		Show:  true,
		Scale: true,
		//GridIndex: 1, // y index 1 // not required
	})

	kline.SetXAxis(xAxe).AddSeries("ohlc", ohlcYAxe)

	kline.Overlap(plotVolumeLine(xAxe, volYaxe)) // Supported charts: Bar/BoxPlot/Line/Scatter/EffectScatter/Kline/HeatMap

	return kline
}

func plotVolumeLine(xAxe []string, yAxe []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | OHLCV | BTC-USDT | 2022-01-01",
			Subtitle: "VOLUME only",
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
		//AXIS
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),

		//charts.WithYAxisOpts(opts.YAxis{
		//	Name:      "Price",
		//	Type:      "value",
		//	Show:      true,
		//	Scale:     true,
		//	GridIndex: 0, // y index 0 // not required
		//}),

		charts.WithYAxisOpts(opts.YAxis{
			//Name: "hidden-left-y-axe",
			//Type: "value",
			Show: false,
			//Scale: true,
			//GridIndex: 0, // y index 0 // not required
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Name:  "Volume",
		Type:  "value",
		Show:  true,
		Scale: true,
		//GridIndex: 1, // y index 1 // not required
	})

	//line.SetXAxis(xAxe).AddSeries("volume", yAxe, charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.SetXAxis(xAxe).AddSeries("volume", yAxe, charts.WithLineChartOpts(opts.LineChart{Smooth: true, YAxisIndex: 1}))

	return line
}

func plotSimpleChart(xAxe []string, yAxe []opts.KlineData) *charts.Kline {
	kline := charts.NewKLine()
	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | OHLCV | BTC-USDT | 2022-01-01",
			Subtitle: "OHLC only",
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
			//Name:  "OHLC",
			//Type:  "value",
			Show:  true,
			Scale: true,
			//GridIndex: 0, // y index 0 // not required
		}),
	)

	kline.SetXAxis(xAxe).AddSeries("ohlc", yAxe)

	return kline
}

func prepareOhlcvData(records [][]string) ([]string, []opts.KlineData, []opts.LineData) {
	start := 0
	if strings.Contains(records[0][0], "CLOSED_AT") {
		start = 1
	}
	x := make([]string, 0)
	ohlcY := make([]opts.KlineData, 0)
	volY := make([]opts.LineData, 0)
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
		volumeVal, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			continue
		}

		x = append(x, record[1])
		ohlcY = append(ohlcY, opts.KlineData{
			// [open, close, lowest, highest]
			Value: [4]float64{openVal, closeVal, lowVal, highVal},
		})
		volY = append(volY, opts.LineData{Value: volumeVal, YAxisIndex: 1})
	}

	return x, ohlcY, volY
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
