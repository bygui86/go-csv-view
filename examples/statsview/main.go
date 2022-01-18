package main

import (
	"time"

	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
)

func main() {
	// set configurations before calling `statsview.New()` method
	viewer.SetConfiguration(
		viewer.WithTheme(viewer.ThemeWesteros),
		viewer.WithAddr("localhost:8090"), // default localhost:18066
	)

	mgr := statsview.New()

	// Start() runs a HTTP server at `localhost:18066` by default.
	go mgr.Start()

	// Stop() will shutdown the http server gracefully
	// mgr.Stop()

	// busy working
	time.Sleep(time.Minute)
}

// Visit your browser at http://localhost:18066/debug/statsview
// Or debug as always via http://localhost:18066/debug/pprof, http://localhost:18066/debug/pprof/heap, ...
