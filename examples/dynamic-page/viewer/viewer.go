package viewer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/bygui86/go-csv-view/examples/dynamic-page/statics"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"go.uber.org/zap"
)

const (
	defaultWidth  = "800px"
	defaultHeight = "600px"
	//defaultTheme  = types.ThemeWesteros
	defaultTheme = types.ThemeMacarons

	defaultInterval  = 2000
	defaultMaxPoints = 30
)

type Viewer struct { // INFO fields globally visible because used by manager.Manager
	Name        string
	Address     string
	AddressPath string
	Graph       *charts.Line
	updater     *Updater
}

// ViewerTemplate defines fields used in the default view template
// WARN: changing names or types, please remember to change also the content of statics.ViewTemplate const
type ViewerTemplate struct { // INFO fields globally visible because used in statics.ViewTemplate
	Interval  int
	MaxPoints int
	Address   string
	ViewPath  string
	ViewID    string
}

type Metrics struct {
	Values []float64 `json:"values"`
	Time   string    `json:"time"`
}

func NewViewer(name, address, addressPath string,
	ctx context.Context, interval, shutdownTimeout int) *Viewer {

	zap.S().Debugf("New Viewer - name %s, address %s, addressPath %s, interval %d, shutdownTimeout %d",
		name, address, addressPath, interval, shutdownTimeout)

	return (&Viewer{}).
		setupGeneral(name, address, addressPath).
		setupGraph().
		setupUpdater(ctx, interval, shutdownTimeout)
}

func (v *Viewer) setupGeneral(name, address, addressPath string) *Viewer {
	v.Name = name
	v.Address = address
	v.AddressPath = addressPath

	zap.S().Debugf("Setup general - name %s, address %s, addressPath %s",
		name, address, addressPath)

	return v
}

func (v *Viewer) setupGraph() *Viewer {
	v.Graph = charts.NewLine()
	v.Graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Line sample"}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  defaultWidth,
			Height: defaultHeight,
			Theme:  defaultTheme,
		}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Time"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Size", AxisLabel: &opts.AxisLabel{Formatter: "{value} MB"}}),
	)
	v.Graph.SetXAxis([]string{}).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	// v.Name NOT ENOUGH!!
	v.Graph.AddJSFuncs(v.generateViewTemplate())

	v.Graph.
		AddSeries("Sys", []opts.LineData{}).
		AddSeries("Inuse", []opts.LineData{}).
		AddSeries("MSpan Sys", []opts.LineData{}).
		AddSeries("MSpan Inuse", []opts.LineData{})

	zap.S().Debugf("Setup graph")

	return v
}

func (v *Viewer) setupUpdater(ctx context.Context, interval, shutdownTimeout int) *Viewer {
	v.updater = newUpdater(ctx, interval, shutdownTimeout)

	zap.S().Debugf("Setup updater")

	return v
}

func (v *Viewer) Start() {
	zap.S().Infof("%s Viewer start...", v.Name)

	v.updater.Start()
}

func (v *Viewer) Stop() {
	zap.S().Infof("%s Viewer stop...", v.Name)

	v.updater.Stop()
}

func (v *Viewer) generateViewTemplate() string {
	zap.S().Debugf("Generate view template")

	tpl, tplErr := template.New("view").Parse(statics.ViewTemplate)
	if tplErr != nil {
		log.Fatalf("template parsing failed: %s" + tplErr.Error())
	}

	buf := bytes.Buffer{}
	execErr := tpl.Execute(
		&buf,
		&ViewerTemplate{ // TODO what about & ?
			Interval:  defaultInterval,
			MaxPoints: defaultMaxPoints,
			Address:   v.Address,
			ViewPath:  v.AddressPath,
			ViewID:    v.Graph.ChartID,
		},
	)
	if execErr != nil {
		log.Fatalf("template execution failed: %s" + execErr.Error())
	}

	return buf.String()
}

func (v *Viewer) Handler(w http.ResponseWriter, _ *http.Request) {
	zap.S().Infof("%s Viewer handler", v.Name)

	v.updater.ViewerRefresh() // INFO let the updater know last time the view was refreshed

	metrics := Metrics{
		Values: []float64{
			fixedPrecision(float64(goMemStats.MemStats.StackSys)/1024/1024, 2),   // "Sys" series
			fixedPrecision(float64(goMemStats.MemStats.StackInuse)/1024/1024, 2), // "Inuse" series
			fixedPrecision(float64(goMemStats.MemStats.MSpanSys)/1024/1024, 2),   // "MSpan Sys" series
			fixedPrecision(float64(goMemStats.MemStats.MSpanInuse)/1024/1024, 2), // "MSpan Inuse" series
		},
		Time: goMemStats.PointTime,
	}

	zap.S().Debugf("%s Viewer new metrics: Sys %.2f, Inuse %.2f, MSpan Sys %.2f, MSpan Inuse %.2f, time %s",
		v.Name, metrics.Values[0], metrics.Values[1], metrics.Values[2], metrics.Values[3], metrics.Time)

	metricsBytes, jsonErr := json.Marshal(metrics)
	if jsonErr != nil {
		zap.S().Errorf("view json marhal failed: %s", jsonErr.Error())
		return
	}

	_, wrErr := w.Write(metricsBytes)
	if wrErr != nil {
		zap.S().Errorf("view rendering failed: %s", wrErr.Error())
	}
}

func fixedPrecision(n float64, p int) float64 {
	var r float64
	switch p {
	case 2:
		r, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", n), 64)
	case 6:
		r, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", n), 64)
	}
	return r
}
