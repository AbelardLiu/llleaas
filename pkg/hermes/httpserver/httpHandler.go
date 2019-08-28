package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"html/template"
	"io/ioutil"
	"lll.github.com/llleaas/pkg/common/log"
	"net/http"
	"strings"
)

type Event struct {
	Type string	`json:"type"`
	Message string `json:"message"`
}

type HttpHandler struct {
	Name string
	Router *mux.Router
	Connections map[string] *websocket.Conn
}

const HERMES_PATH = "/hermes"

func NewHttpHandler(name string) (*HttpHandler) {
	return &HttpHandler{
		Name: name,
		Router: mux.NewRouter(),
		Connections: make(map[string] *websocket.Conn),
	}
}

func (h *HttpHandler)Start() error {
	h.Router.HandleFunc(HERMES_PATH + "/registry/message", h.Message).Methods(http.MethodPost)
	h.Router.HandleFunc(HERMES_PATH + "/registry/index", h.Index).Methods(http.MethodGet)

	h.Router.Handle(HERMES_PATH + "/registry/upper", websocket.Handler(h.Upper))

	return nil
}

func (h *HttpHandler)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func (h *HttpHandler)Index(w http.ResponseWriter,r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func (h *HttpHandler)Message(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.GetLogger().Errorf("Http handler message read message error: %v", err)
		return
	}

	var message Event
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.GetLogger().Errorf("Http handler message unmarshal json error: %v", err)
		return
	}

	websocket.JSON.Send(h.Connections["123"], message);

	res, _ := json.Marshal(message)
	w.Write(res)
}

func (h *HttpHandler) Upper(ws *websocket.Conn) {
	var err error
	for {
		var reply Event

		if err = websocket.JSON.Receive(ws, &reply); err != nil {
			fmt.Println(err)
			break
		}

		if reply.Type == "register" {
			h.Connections[reply.Message] = ws
		}

		reply.Message = strings.ToUpper(reply.Message)
		if err = websocket.JSON.Send(ws, reply); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
