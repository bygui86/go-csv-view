
# Golang CSV Viewer - go-echarts examples

## Not yet supported by go-echarts (or not well documented)

- full kline chart like https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush

- left/right tooltip like https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush

- markarea (zone of evidence) like
  - https://echarts.apache.org/examples/en/editor.html?c=line-sections
  - https://echarts.apache.org/examples/en/editor.html?c=area-rainfall
  - https://echarts.apache.org/examples/en/editor.html?c=candlestick-brush

- inverse y-axis like 
  - https://echarts.apache.org/examples/en/editor.html?c=grid-multiple
  - https://echarts.apache.org/examples/en/editor.html?c=area-rainfall

- multiple charts, single datazoom, single panel like https://echarts.apache.org/examples/en/editor.html?c=grid-multiple

- multiple x-axis like https://echarts.apache.org/examples/en/editor.html?c=multiple-x-axis

- linear gradient like https://echarts.apache.org/examples/en/editor.html?c=area-stack-gradient

- animations
  - draw line animation like https://echarts.apache.org/examples/en/editor.html?c=candlestick-sh
  - bar animation delay like https://echarts.apache.org/examples/en/editor.html?c=bar-animation-delay

---

## `go-echarts officials`

There are variety of charts, such as `bar`, `line`, `pie`, `radaer`, etc.

```bash
cd go-echarts-official && go run .

# in another terminal window
open http://localhost:8089
```

## `go-echarts web-server`

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

## `bar`

```bash
cd bar && go run main.go

open games.html
```

## `two Y-Axis`

```bash
cd two-y-axis && go run main.go

open two-y-axis.html
```

## `line-bar`

```bash
cd line-bar && go run main.go

open trades.html
```

## `tree`

```bash
cd tree && go run main.go

open tree.html
```

## `kline`

```bash
cd kline && go run main.go

open ohlcv.html
```

## `kline` using go-tachart

```bash
cd go-tachart-official && go run main.go

open kline.html
```

## `kline` using unofficial modified go-echarts

```bash
cd kline-custom && go run main.go

open ohlcv.html
```

## `statsview`

```bash
cd statsview && go run main.go

# in another terminal window
open http://localhost:8090/debug/statsview
```

## `dynamic-page`

```bash
cd dynamic-page && go run main.go

# in another terminal window
open http://localhost:8080/page
```
