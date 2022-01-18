package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/iamjinlei/go-tachart/tachart"
)

const (
	csvFilePath  = "ohlcv.csv"
	htmlFilePath = "kline.html"
)

func main() {
	records, loadErr := loadCsv(csvFilePath)
	if loadErr != nil {
		log.Fatal(loadErr)
	}

	cdls := prepareData(records)

	events := []tachart.Event{
		{
			Type:        tachart.Short,
			Label:       cdls[50].Label,
			Description: "This is a demo event description. Randomly pick this candle to go short on " + cdls[5].Label,
		},
	}

	cfg := tachart.NewConfig().
		SetChartWidth(1080).
		SetChartHeight(800).
		AddOverlay(
			tachart.NewSMA(5),
			tachart.NewSMA(50),
		).
		AddIndicator(
			tachart.NewMACD(12, 26, 9),
		).
		UseRepoAssets() // serving assets file from current repo, avoid network access

	c := tachart.New(*cfg)
	err := c.GenStatic(cdls, events, htmlFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareData(records [][]string) []tachart.Candle {
	start := 0
	if strings.Contains(records[0][0], "CLOSED_AT") {
		start = 1
	}
	candles := make([]tachart.Candle, 0)
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

		candles = append(candles,
			tachart.Candle{Label: record[1], O: openVal, C: closeVal, L: lowVal, H: highVal, V: volumeVal})
	}

	return candles
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
