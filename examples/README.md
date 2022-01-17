
# Golang CSV Viewer - Examples

## TODO list

- iamjinlei/go-tachart
- full kline chart like https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush
- left/right tooltip like https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush
- area of evidence like https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush
- min, max, avg of window like https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh
- animations
  - draw line animation like https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh
  - bar animation delay like https://echarts.apache.org/examples/en/editor.html?c=bar-animation-delay
- changing color of x-axis zone like https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh
- custom colors like https://echarts.apache.org/examples/en/editor.html?c=bubble-gradient
- inverted values like
  - https://echarts.apache.org/examples/en/editor.html?c=area-rainfall
  - https://echarts.apache.org/examples/en/editor.html?c=grid-multiple
- put peaks in evidence like https://echarts.apache.org/examples/en/editor.html?c=line-sections
- multiple x-axis like https://echarts.apache.org/examples/en/editor.html?c=multiple-x-axis
- thresholds like https://echarts.apache.org/examples/en/editor.html?c=line-aqi
- filled lines like https://echarts.apache.org/examples/en/editor.html?c=area-stack-gradient

---

## `official`

There are variety of charts, such as `bar`, `line`, `pie`, `radaer`, etc.

```bash
cd official && go run .

# in another terminal window
open http://localhost:8089
```

## `web-server`

```bash
cd web-server && go run main.go

# in another terminal window
open http://localhost:8081
```

## `line`

```bash
cd line && go run main.go

open ohlcv.html
```

## `two Y-Axis`

```bash
cd two-y-axis && go run main.go

open two-y-axis.html
```

## `bar`

```bash
cd bar && go run main.go

open games.html
```

## `line-bar`

```bash
cd line-bar && go run main.go

open trades.html
```

## `kline`

```bash
cd kline && go run main.go

open ohlcv.html
```
