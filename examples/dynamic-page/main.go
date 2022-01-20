package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bygui86/go-csv-view/examples/dynamic-page/manager"
	"github.com/bygui86/go-csv-view/examples/dynamic-page/viewer"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

const (
	logEncoding = "console"
	//logEncoding = "json"

	logLevel = zapcore.InfoLevel
	//logLevel = zapcore.DebugLevel

	address         = "localhost:8080"
	pagePath        = "/page"
	viewPath        = pagePath + "/view/%s"
	interval        = 2000 // milliseconds
	shutdownTimeout = 10   // seconds
)

func main() {
	setupLogger()

	zap.L().Info("Starting dynamic page")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	viewerName := "stack"
	stackViewer := viewer.NewViewer(
		viewerName,
		address,
		fmt.Sprintf(viewPath, viewerName),
		ctx,
		interval,
		shutdownTimeout,
	)

	zap.L().Info("Viewers created")

	mng := manager.NewManager(
		address,
		pagePath,
		ctx,
		shutdownTimeout,
		stackViewer,
	)

	zap.L().Info("Manager created")

	err := mng.Start()
	if err != nil {
		zap.S().Fatalf("manager start failed: %s", err.Error())
	}
}

func setupLogger() {
	logger, err := zap.Config{
		Encoding:         logEncoding,
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			MessageKey:   "message",
		},
	}.Build()
	if err != nil {
		log.Fatalf("logger setup failed: %s", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
