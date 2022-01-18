package charts

import (
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/opts"
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/render"
	"github.com/bygui86/go-csv-view/alaingilbert-go-echarts/types"
)

// Bar3D represents a 3D bar chart.
type Bar3D struct {
	Chart3D
}

// Type returns the chart type.
func (Bar3D) Type() string { return types.ChartBar3D }

// NewBar3D creates a new 3D bar chart.
func NewBar3D() *Bar3D {
	c := &Bar3D{}
	c.initBaseConfiguration()
	c.Renderer = render.NewChartRender(c, c.Validate)
	c.initChart3D()
	return c
}

// AddSeries adds the new series.
func (c *Bar3D) AddSeries(name string, data []opts.Chart3DData, options ...SeriesOpts) *Bar3D {
	c.addSeries(types.ChartScatter3D, name, data, options...)
	return c
}
