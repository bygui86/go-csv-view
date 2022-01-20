package viewer

import (
	"context"
	"runtime"
	"time"

	"go.uber.org/zap"
)

const (
	defaultTimeFormat = "15:04:05"
)

// TODO merge together with Viewer
type Updater struct {
	viewerRefresh   int64 // TODO ? - viewerRefresh tick from Viewer
	ctx             context.Context
	cancel          context.CancelFunc
	interval        int // milliseconds
	ticker          *time.Ticker
	pollingRun      bool
	shutdownTimeout int
}

type statistics struct { // INFO fields globally visible because updated by runtime.ReadMemStats and used by Viewer
	MemStats  *runtime.MemStats
	PointTime string // INFO: this is the time of the last point added to the graph
}

var goMemStats = &statistics{MemStats: &runtime.MemStats{}} // TODO put inside Updater or Viewer?

func newUpdater(ctx context.Context, interval, shutdownTimeout int) *Updater {

	zap.S().Debugf("New Updater - interval %d, shutdownTimeout %d",
		interval, shutdownTimeout)

	return (&Updater{}).
		setupGeneral(interval, shutdownTimeout).
		setupCtx(ctx)
}

func (u *Updater) setupGeneral(interval int, shutdownTimeout int) *Updater {
	u.interval = interval
	u.shutdownTimeout = shutdownTimeout
	u.pollingRun = false

	zap.S().Debugf("Setup general - interval %d, shutdownTimeout %d",
		interval, shutdownTimeout)

	return u
}

func (u *Updater) setupCtx(ctx context.Context) *Updater {
	if ctx != nil {
		u.ctx, u.cancel = context.WithCancel(ctx)
	} else {
		u.ctx, u.cancel = context.WithCancel(context.Background())
	}

	zap.S().Debugf("Setup ctx")

	return u
}

func (u *Updater) Start() {
	if !u.pollingRun {
		zap.S().Debug("Start...")
		go u.polling()
	} else {
		zap.S().Warn("Already polling")
	}
}

func (u *Updater) Stop() {
	if u.pollingRun {
		zap.S().Debug("Stop...")
		_, cancel := context.WithTimeout(context.Background(), time.Duration(u.shutdownTimeout)*time.Second)
		defer cancel()

		u.ticker.Stop()
		u.ctx.Done()

		time.Sleep(time.Duration(u.shutdownTimeout) * time.Second)

		u.cancel()
	} else {
		zap.S().Warn("Nothing to stop")
	}
}

func (u *Updater) polling() {
	zap.S().Info("Polling...")

	u.ticker = time.NewTicker(time.Duration(u.interval) * time.Millisecond)
	u.pollingRun = true

	for {
		select {
		case <-u.ticker.C:
			if u.viewerRefresh > time.Now().Unix() { // INFO fetch new values only if the last time the view was refreshed is later than "now"
				zap.S().Debug("Fetching new metrics")
				runtime.ReadMemStats(goMemStats.MemStats)
				goMemStats.PointTime = time.Now().Format(defaultTimeFormat)
			} else {
				zap.S().Debug("Not yet time to fetch new metrics")
			}
		case <-u.ctx.Done():
			zap.S().Warn("Stop polling")
			u.pollingRun = false
			return
		}
	}
}

// TODO move to Viewer
// TODO at every refresh, we can fetch new data completely avoiding the polling() routine
// ViewerRefresh is used by Viewer to update last time it was refreshed
func (u *Updater) ViewerRefresh() {
	u.viewerRefresh = time.Now().Unix() + int64(float64(u.interval)/1000.0)*2

	zap.S().Debugf("New Viewer refresh time: %d", u.viewerRefresh)
}
