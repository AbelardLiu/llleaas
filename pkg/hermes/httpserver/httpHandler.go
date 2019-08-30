package httpserver

import (
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"lll.github.com/llleaas/cmd/hermes/app/options"
	"lll.github.com/llleaas/pkg/common/log"
	"lll.github.com/llleaas/pkg/hermes/faasmanager"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)



type HttpHandler struct {
	Name string
	Router *mux.Router
	FaasManager faasmanager.FaaSManager
	Option *options.HermesOption
	Connections map[string] *websocket.Conn
	FaasProxy *httputil.ReverseProxy
}

const HERMES_PATH = "/hermes"


func startFaasManager(name string, option *options.HermesOption) faasmanager.FaaSManager {
	faasmgr := faasmanager.NewWebsocketManager(name, option)
	err := faasmgr.Start(HERMES_PATH)
	if err != nil {
		log.GetLogger().Fatalf("http handler start faas manager error: %v", err)
		panic(err)
	}

	return faasmgr
}

func NewHttpHandler(name string, option *options.HermesOption) (*HttpHandler) {
	faasmgr := startFaasManager(name, option)

	return &HttpHandler{
		Name: name,
		Option: option,
		Router: mux.NewRouter(),
		FaasManager: faasmgr,
		Connections: make(map[string] *websocket.Conn),
		FaasProxy: nil,
	}
}

func (h *HttpHandler)Start() error {
	// create reverse proxy to default faas
	faasUrl, err := url.Parse(h.Option.DefaultFaasUrl)
	if err != nil {
		log.GetLogger().Errorf("hermes http handler start default faas url parse error: %v", err)
		return err
	}
	h.FaasProxy = httputil.NewSingleHostReverseProxy(faasUrl)

	return nil
}

func (h *HttpHandler)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	if strings.HasPrefix(r.URL.Path, HERMES_PATH) {
		// hermes path
		h.FaasManager.ServeHTTP(w, r)
		return
	}

	// faas maanger path
	h.FaasProxy.ServeHTTP(w, r)

}