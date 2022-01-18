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
	csvFilePath  = "trades.csv"
	htmlFilePath = "trades.html"
)

func main() {
	records, loadErr := loadCsv(csvFilePath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	xAxe, lineYAxe, barYAxe := prepareData(records)

	line := plotLine(xAxe, lineYAxe)
	bar := plotBar(xAxe, barYAxe, nil)
	lineBar := plotBar(xAxe, barYAxe, line)

	pageErr := createHtml(htmlFilePath, lineBar, line, bar)
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

func plotBar(xAxe []string, yAxe []opts.BarData, overlapChart charts.Overlaper) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | BTC-USDT",
			Subtitle: "TRADES of 2022-01-01",
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
			Name:  "Size",
			Type:  "value",
			Show:  true,
			Scale: true,
			//GridIndex: 0, // y index 0 // not required
		}),
	)

	bar.ExtendYAxis(opts.YAxis{
		Name:  "Price",
		Type:  "value",
		Show:  true,
		Scale: true,
		//GridIndex: 1, // y index 1 // not required
	})

	//bar.SetXAxis(xAxe).AddSeries("size", yAxe)
	bar.SetXAxis(xAxe).AddSeries("size", yAxe, charts.WithRippleEffectOpts(opts.RippleEffect{Period: 5}))

	if overlapChart != nil {
		bar.Overlap(overlapChart)
	}

	return bar
}

func plotLine(xAxe []string, yAxe []opts.LineData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Binance | BTC-USDT",
			Subtitle: "TRADES of 2022-01-01",
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
			// HIDDEN
			Show: false,
			//GridIndex: 0, // y index 0 // not required
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Name:  "Price",
		Type:  "value",
		Show:  true,
		Scale: true,
		//GridIndex: 1, // y index 1 // not required
	})

	line.SetXAxis(xAxe).AddSeries("price", yAxe,
		charts.WithLineChartOpts(opts.LineChart{Smooth: true, YAxisIndex: 1}),
		//charts.WithRippleEffectOpts(opts.RippleEffect{Period: 5, BrushType: "stroke"}),
		charts.WithRippleEffectOpts(opts.RippleEffect{Period: 5}),
	)

	return line
}

func prepareData(records [][]string) ([]string, []opts.LineData, []opts.BarData) {
	start := 0
	if strings.Contains(records[0][0], "TIMESTAMP") {
		start = 1
	}
	x := make([]string, 0)
	lineY := make([]opts.LineData, 0)
	barY := make([]opts.BarData, 0)
	for _, record := range records[start:] {
		// TIMESTAMP,TRADE_ID,PRICE,SIDE,SIZE,BUYER_ORDER_ID,SELLER_ORDER_ID,COMPONENT,BUCKET
		priceVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		sizeVal, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}
		x = append(x, record[0])
		lineY = append(lineY, opts.LineData{Value: priceVal, YAxisIndex: 1})
		barY = append(barY, opts.BarData{Value: sizeVal})
	}

	return x, lineY, barY
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
