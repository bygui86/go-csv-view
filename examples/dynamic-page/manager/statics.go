package manager

import (
	"log"
	"net/http"

	"github.com/bygui86/go-csv-view/examples/dynamic-page/statics"
)

func echartJsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(statics.EchartJS))
	if err != nil {
		log.Fatal(err)
	}
}

func jqueryJsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(statics.JqueryJS))
	if err != nil {
		log.Fatal(err)
	}
}

func westerosJsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(statics.WesterosJS))
	if err != nil {
		log.Fatal(err)
	}
}

func macaronsJsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte(statics.MacaronsJS))
	if err != nil {
		log.Fatal(err)
	}
}
