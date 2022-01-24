package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/charts"
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/components"
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/opts"
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/types"
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

	dataset := prepareOhlcvData(records)

	kline := plotChart(dataset)

	pageErr := createHtml(htmlFilePath, kline)
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

func plotChart(dataset [][]interface{}) *charts.Kline {
	kline := charts.NewKLine()

	kline.Dataset = opts.Dataset{Source: dataset}

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			//Title: "Binance | OHLCV | BTC-USDT | 2022-01-01",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			//Theme:  "dark",
			Theme: types.ThemeVintage,
			//Theme:  types.ThemeWesteros,
			//Theme:  types.ThemeWonderland,
			//Theme:  types.ThemeRoma,
			//Theme:  types.ThemeEssos,
			//Width:  "100%",
			//Height: "100%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:         true,
			SelectedMode: "multiple",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Pieces: []opts.Piece{
				{Value: 1, Color: "rgba(0, 218, 60, 0.7)"},
				{Value: -1, Color: "rgba(236, 0, 0, 0.7)"},
			},
			Dimension:   6,
			SeriesIndex: 1,
			Show:        false,
		}),
		charts.WithGridOpts(
			opts.Grid{Bottom: "210", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "210", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "120", Left: "50", Right: "10"},
			opts.Grid{Height: "80", Bottom: "40", Left: "50", Right: "10"},
		),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:        true,
			Trigger:     "axis",
			AxisPointer: &opts.AxisPointer{Type: "line"},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      30,
			End:        70,
			XAxisIndex: []int{0, 1, 2, 3},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      30,
			End:        70,
			XAxisIndex: []int{0, 1, 2, 3},
		}),
		// AXIS
		charts.WithXAxisOpts(opts.XAxis{
			Type:        "category",
			SplitNumber: 20,
			GridIndex:   0,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale:     true,
			GridIndex: 0,
		}),
	)

	kline.ExtendXAxis(opts.XAxis{Type: "category", SplitNumber: 20, GridIndex: 1,
		AxisTick:  &opts.AxisTick{Show: false},
		AxisLabel: &opts.AxisLabel{Show: false},
	})
	kline.ExtendYAxis(opts.YAxis{
		Scale: true, GridIndex: 1, SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{Show: false},
		AxisLine:  &opts.AxisLine{Show: false},
		AxisTick:  &opts.AxisTick{Show: false},
		SplitLine: &opts.SplitLine{Show: false},
	})
	kline.ExtendXAxis(opts.XAxis{Type: "category", SplitNumber: 20, GridIndex: 2,
		AxisTick:  &opts.AxisTick{Show: false},
		AxisLabel: &opts.AxisLabel{Show: false},
	})
	kline.ExtendYAxis(opts.YAxis{
		Scale: true, GridIndex: 2, SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{Show: false},
		AxisLine:  &opts.AxisLine{Show: false},
		AxisTick:  &opts.AxisTick{Show: false},
		SplitLine: &opts.SplitLine{Show: false},
	})
	kline.ExtendXAxis(opts.XAxis{
		Type: "category", SplitNumber: 20, GridIndex: 3,
		AxisTick:  &opts.AxisTick{Show: false},
		AxisLine:  &opts.AxisLine{Show: false},
		AxisLabel: &opts.AxisLabel{Show: false},
	})
	kline.ExtendYAxis(opts.YAxis{
		Scale: true, GridIndex: 3, SplitNumber: 2,
		AxisLabel: &opts.AxisLabel{Show: true},
		AxisLine:  &opts.AxisLine{Show: true},
		AxisTick:  &opts.AxisTick{Show: true},
		SplitLine: &opts.SplitLine{Show: true},
	})

	kline.AddSeries("candlestick", nil,
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "#00da3c",
			Color0:       "#ec0000",
			BorderColor:  "#008F28",
			BorderColor0: "#8A0000"}),
		charts.WithEncodeOpts(opts.Encode{
			X: 0,
			Y: []int{1, 2, 3, 4}},
		),
	)

	volumeBarChart := charts.NewBar()
	volumeBarChart.AddSeries("Volume", nil,
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "#7fbe9e"}),
		charts.WithBarChartOpts(opts.BarChart{Type: "bar", XAxisIndex: 1, YAxisIndex: 1}),
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 5}))

	ema30LineChart := charts.NewLine()
	ema30LineChart.SetGlobalOptions(charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20, GridIndex: 0}), charts.WithYAxisOpts(opts.YAxis{Scale: true, GridIndex: 0}))
	ema30LineChart.AddSeries("EMA30", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 8}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(69, 140, 255, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 0, YAxisIndex: 0}))

	ema200LineChart := charts.NewLine()
	ema200LineChart.SetGlobalOptions(charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20, GridIndex: 0}), charts.WithYAxisOpts(opts.YAxis{Scale: true, GridIndex: 0}))
	ema200LineChart.AddSeries("EMA200", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 13}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(255, 174, 69, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 0, YAxisIndex: 0}))

	ema10LineChart := charts.NewLine()
	ema10LineChart.AddSeries("EMA10", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 7}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(69, 246, 255, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))

	bbLowerLineChart := charts.NewLine()
	bbLowerLineChart.AddSeries("EMA10", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 15}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(255, 252, 89, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))
	bbUpperLineChart := charts.NewLine()
	bbUpperLineChart.AddSeries("EMA10", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 16}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(255, 252, 89, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))

	buysChart := charts.NewScatter()
	buysChart.AddSeries("Buy", nil,
		charts.WithScatterChartOpts(opts.ScatterChart{Symbol: "circle"}),
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 9}),
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "#00b500"}))
	buysChart.AddSeries("Sell", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 10}),
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "#ff0000"}))

	// TODO group all together
	kline.Overlap(volumeBarChart)
	kline.Overlap(ema10LineChart)
	kline.Overlap(ema30LineChart)
	kline.Overlap(ema200LineChart)
	kline.Overlap(buysChart)
	kline.Overlap(bbLowerLineChart)
	kline.Overlap(bbUpperLineChart)

	macdChart := charts.NewLine()
	macdChart.AddSeries("MACD", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 11}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 2, YAxisIndex: 2}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "#f00"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))
	macd9Chart := charts.NewLine()
	macd9Chart.AddSeries("MACD signal", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 12}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 2, YAxisIndex: 2}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "#0f0"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))

	// TODO group all together
	kline.Overlap(macdChart)
	kline.Overlap(macd9Chart)

	rsiChart := charts.NewLine()
	rsiChart.AddSeries("RSI", nil,
		charts.WithEncodeOpts(opts.Encode{X: 0, Y: 14}),
		charts.WithLineChartOpts(opts.LineChart{XAxisIndex: 3, YAxisIndex: 3}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "rgba(169, 84, 255, 0.5)"}),
		charts.WithItemStyleOpts(opts.ItemStyle{Opacity: 0.01}))

	// TODO group all together
	kline.Overlap(rsiChart)

	return kline
}

func prepareOhlcvData(records [][]string) [][]interface{} {
	start := 0
	if strings.Contains(records[0][0], "CLOSED_AT") {
		start = 1
	}
	dataset := make([][]interface{}, 0)
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

		dataset = append(dataset, []interface{}{
			record[1],
			openVal,
			closeVal,
			lowVal,
			highVal,
			volumeVal,
			nil,                       // getSign(d, idx),
			openVal - (openVal / 10),  // ema10[idx],
			openVal - (openVal / 30),  // ema30[idx],
			nil,                       // buy,
			nil,                       // sell,
			lowVal,                    // macdHist[idx],
			highVal,                   // signalHist[idx],
			openVal - (openVal / 200), // ema200[idx],
			closeVal,                  // rsiHist[idx],
			nil,                       // bbUpperHist[idx],
			nil,                       // bbLowerHist[idx],
		})
	}

	return dataset
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
