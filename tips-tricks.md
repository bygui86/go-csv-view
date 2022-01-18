
# go-echarts - Tips & tricks

## scaling axis values based on zoom

Example:   https://echarts.apache.org/examples/en/editor.html?c=area-simple

Go elements:

  - opts.XAxis.Scale
  - opts.YAxis.Scale

## legend

Example:   https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh

Go element:   opts.Legend

## tooltip

Example:   https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush

Go element:   opts.Tooltip

## markline

Example:   https://echarts.apache.org/examples/en/editor.html?c=line-marker

Go elements:

  - opts.MarkLineNameTypeItem
  - MarkLineNameYAxisItem
  - MarkLineNameXAxisItem

## markpoint

Example:   https://echarts.apache.org/examples/en/editor.html?c=line-marker

Go elements:

  - opts.MarkPointStyle
  - opts.MarkPointNameTypeItem
  - opts.MarkPointNameCoordItem

## visualmap (color thresholds)

`/!\ WARN` buggy feature or really complex configurations

Example:   https://echarts.apache.org/examples/en/editor.html?c=line-aqi

Go element:   opts.VisualMap

## flip colors of y-axis zones

Example:   https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh

Go elements:

- opts.XAxis.SplitArea
- opts.YAxis.SplitArea

## custom colors

Example:   https://echarts.apache.org/examples/en/editor.html?c=multiple-x-axis

Go element:   opts.Colors

## toolbox

Example:   https://echarts.apache.org/examples/en/editor.html?c=line-marker

Go element:   opts.Toolbox
