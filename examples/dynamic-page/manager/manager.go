package manager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bygui86/go-csv-view/examples/dynamic-page/statics"
	"github.com/bygui86/go-csv-view/examples/dynamic-page/viewer"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Manager struct {
	address         string
	pagePath        string
	staticsPath     string
	srv             *http.Server
	ctx             context.Context
	cancel          context.CancelFunc
	shutdownTimeout int
	mux             *http.ServeMux
	page            *components.Page
	viewers         []*viewer.Viewer
}

func NewManager(address, pagePath string,
	ctx context.Context, shutdownTimeout int,
	viewers ...*viewer.Viewer) *Manager {

	zap.S().Debugf("New Manager - address %s, pagePath %s, viewers %d, shutdownTimeout %d",
		address, pagePath, len(viewers), shutdownTimeout)

	return (&Manager{}).
		setupGeneral(address, pagePath, shutdownTimeout).
		setupHttpServer().
		setupCtx(ctx).
		setupPage().
		registerViewers(viewers...).
		setupMux()
}

func (m *Manager) setupGeneral(address, pagePath string, shutdownTimeout int) *Manager {
	templates.PageTpl = statics.PageTemplate
	m.address = address
	m.pagePath = pagePath
	m.staticsPath = m.pagePath + "/statics"
	m.shutdownTimeout = shutdownTimeout

	zap.S().Debugf("Setup general - address %s, pagePath %s, staticsPath %s, shutdownTimeout %d",
		m.address, m.pagePath, m.staticsPath, m.shutdownTimeout)

	return m
}

func (m *Manager) setupHttpServer() *Manager {
	m.srv = &http.Server{
		Addr:         m.address,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		//MaxHeaderBytes: 1 << 20,
	}

	zap.S().Debugf("Setup HTTP server - address %s", m.address)

	return m
}

func (m *Manager) setupCtx(ctx context.Context) *Manager {
	if ctx != nil {
		m.ctx, m.cancel = context.WithCancel(ctx)
	} else {
		m.ctx, m.cancel = context.WithCancel(context.Background())
	}

	zap.S().Debug("Setup ctx")

	return m
}

func (m *Manager) setupPage() *Manager {
	m.page = components.NewPage()
	m.page.PageTitle = "Dynamic page example"
	m.page.AssetsHost = fmt.Sprintf("http://%s%s/", m.address, m.staticsPath)
	//m.page.Assets.JSAssets.Add("echarts.min.js") // TODO why not required?
	m.page.Assets.JSAssets.Add("jquery.min.js")
	//m.page.Assets.JSAssets.Add(""westeros.js") // TODO why not required?
	//m.page.Assets.JSAssets.Add(""macarons.js") // TODO why not required?

	zap.S().Debugf("Setup page - title %s, assetsHosts %s",
		m.page.PageTitle, m.page.AssetsHost)

	return m
}

func (m *Manager) registerViewers(viewers ...*viewer.Viewer) *Manager {
	m.viewers = append(m.viewers, viewers...)
	zap.S().Infof("Registering %d viewers", len(m.viewers))

	for _, v := range viewers {
		zap.S().Infof("Registering viewer - name %s, addressPath %s, graphID %s",
			v.Name, v.AddressPath, v.Graph.ChartID)

		m.page.AddCharts(v.Graph)
	}

	return m
}

func (m *Manager) setupMux() *Manager {
	m.mux = http.NewServeMux()

	m.mux.HandleFunc(m.pagePath, m.pageHandler)

	zap.S().Debug("Listening on %s", m.pagePath)

	for _, v := range m.viewers {
		//m.mux.HandleFunc(fmt.Sprintf(viewPath, v.Name), v.Handler)
		m.mux.HandleFunc(v.AddressPath, v.Handler)
		zap.S().Debug("Listening on %s", v.AddressPath)
	}

	echartsPath := fmt.Sprintf("%s/%s", m.staticsPath, "echarts.min.js")
	jqueryPath := fmt.Sprintf("%s/%s", m.staticsPath, "jquery.min.js")
	themeWesteros := fmt.Sprintf("%s/%s", m.staticsPath, "themes/westeros.js")
	themeMacarons := fmt.Sprintf("%s/%s", m.staticsPath, "themes/macarons.js")

	m.mux.HandleFunc(echartsPath, echartJsHandler)
	m.mux.HandleFunc(jqueryPath, jqueryJsHandler)
	m.mux.HandleFunc(themeWesteros, westerosJsHandler)
	m.mux.HandleFunc(themeMacarons, macaronsJsHandler)

	zap.S().Debug("Listening on %s", echartsPath)
	zap.S().Debug("Listening on %s", jqueryPath)
	zap.S().Debug("Listening on %s", themeWesteros)
	zap.S().Debug("Listening on %s", themeMacarons)

	m.srv.Handler = cors.AllowAll().Handler(m.mux)

	zap.S().Debug("Setup mux")

	return m
}

func (m *Manager) Start() error {
	zap.S().Info("Start...")

	for _, v := range m.viewers {
		v.Start()
	}

	return m.srv.ListenAndServe()
}

func (m *Manager) Stop() {
	zap.S().Info("Stop...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.shutdownTimeout)*time.Second)
	defer cancel()

	for _, v := range m.viewers {
		v.Stop()
	}

	err := m.srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down HTTP server: %s", err.Error())
	}

	time.Sleep(time.Duration(m.shutdownTimeout) * time.Second)

	m.cancel()
}

func (m *Manager) pageHandler(w http.ResponseWriter, _ *http.Request) {
	zap.S().Info("Page handler")

	err := m.page.Render(w)
	if err != nil {
		zap.S().Fatalf("page rendering failed: %s", err.Error())
	}
}
