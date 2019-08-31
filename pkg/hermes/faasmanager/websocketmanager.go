package faasmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"html/template"
	"io/ioutil"
	"lll.github.com/llleaas/cmd/hermes/app/options"
	"lll.github.com/llleaas/pkg/common/log"
	faascommon "lll.github.com/llleaas/pkg/hermes/common"
	"lll.github.com/llleaas/pkg/hermes/function"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

type WebsocketFaas struct {
	FaasSpec   faascommon.FaaSSpec
	Connection *websocket.Conn
}

type WebsocketManager struct {
	Name string
	Router *mux.Router
	Option *options.HermesOption
	FaaSInstances map[string] faascommon.FaaSInstance
	FunctionManager function.FunctionManager
	FaasMux sync.Mutex
	FaasProxy *httputil.ReverseProxy
}

func NewWebsocketFaas(spec faascommon.FaaSSpec, ws *websocket.Conn) *WebsocketFaas {
	return &WebsocketFaas{
		FaasSpec:   spec,
		Connection: ws,
	}
}

func (f *WebsocketFaas)Send(event faascommon.Event) error {
	if err := websocket.JSON.Send(f.Connection, event); err != nil {
		log.GetLogger().Errorf("websocket faas json send error: %v", err)
		return err
	}

	return nil
}

func (f *WebsocketFaas)Recv() (faascommon.Event ,error) {
	var event faascommon.Event
	if err := websocket.JSON.Receive(f.Connection, &event); err != nil {
		log.GetLogger().Errorf("websocket faas json send error: %v", err)
		return faascommon.Event{}, err
	}

	return event, nil
}

func (f *WebsocketFaas)Info() string {
	res, err := json.Marshal(f.Spec)
	if err != nil {
		log.GetLogger().Errorf("websocket faas instance info error: %v", err)
		return ""
	}

	return string(res)
}

func (f *WebsocketFaas)Spec() faascommon.FaaSSpec {
	return f.FaasSpec
}


func NewWebsocketManager(name string, option *options.HermesOption) *WebsocketManager{
	return &WebsocketManager{
		Name: name,
		Option: option,
		Router: mux.NewRouter(),
		FaaSInstances: make(map[string] faascommon.FaaSInstance),
		FunctionManager: nil,
		FaasMux: sync.Mutex{},
		FaasProxy: nil,
	}
}

func (m *WebsocketManager) Start(basepath string) error {
	m.Router.HandleFunc(basepath + "/registry/index", m.Index).Methods(http.MethodGet)
	m.Router.HandleFunc(basepath + "/registry/message/{faas_id}", m.Message).Methods(http.MethodPost)
	m.Router.HandleFunc(basepath + "/registry/faas", m.ListFaas).Methods(http.MethodGet)
	m.Router.HandleFunc(basepath + "/registry/faas/{faas_id}", m.GetFaas).Methods(http.MethodGet)

	m.Router.Handle(basepath + "/registry/upper", websocket.Handler(m.Upper))

	// register function handler
	mgr := function.NewBasicFunctionManager()
	err := mgr.RegisterHandler(basepath, m.Router)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager start register function handler error: %v", err)
		return err
	}

	return nil
}

func (m *WebsocketManager) Register(id string, faas faascommon.FaaSInstance) error {
	m.FaasMux.Lock()
	defer m.FaasMux.Unlock()

	m.FaaSInstances[id] = faas

	return nil
}

func (m *WebsocketManager) UnRegister(id string) error {
	m.FaasMux.Lock()
	defer m.FaasMux.Unlock()

	delete(m.FaaSInstances, id)

	return nil
}

func (m *WebsocketManager) GetFaas(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	faas_id := "unknown"
	if val, ok := params["faas_id"]; ok {
		faas_id = val
	}
	faasInstance,err := m.Get(faas_id)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager get faas get instance error: $v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	spec := faasInstance.Spec()
	res, err := json.Marshal(spec)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager get faas marshal spec to json error: $v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
	return
}

func (m *WebsocketManager) ListFaas(w http.ResponseWriter, r *http.Request) {
	faasInstancs := m.List()

	res, err := json.Marshal(faasInstancs)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager list faas marshal spec to json error: $v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
	return
}

func (m *WebsocketManager) Get(id string) (faascommon.FaaSInstance, error) {
	m.FaasMux.Lock()
	defer m.FaasMux.Unlock()

	faas, ok := m.FaaSInstances[id]
	if !ok {
		return &WebsocketFaas{},errors.New("not find faas instance")
	} else {
		return faas, nil
	}
}

func (m *WebsocketManager) List() map[string]faascommon.FaaSInstance {
	m.FaasMux.Lock()
	defer m.FaasMux.Unlock()

	res := make(map[string]faascommon.FaaSInstance)

	for k,v := range m.FaaSInstances {
		res[k] = v
	}

	return res
}

func (m *WebsocketManager)ServeHTTP(w http.ResponseWriter,r *http.Request){
	m.Router.ServeHTTP(w, r)
}

func (m *WebsocketManager)Index(w http.ResponseWriter,r *http.Request) {
	if r.Method != "GET" {
		return
	}

	log.GetLogger().Info("hello")

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func (m *WebsocketManager)Message(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.GetLogger().Errorf("Http handler message read message error: %v", err)
		return
	}

	var message faascommon.Event
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.GetLogger().Errorf("Http handler message unmarshal json error: %v", err)
		return
	}

	params := mux.Vars(r)

	faas_id := "unknown"
	if val, ok := params["faas_id"]; ok {
		faas_id = val
	}
	log.GetLogger().Debugf("hello : %v", faas_id)
	faasInstance, err := m.Get(faas_id)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager message get faas instance error: %v", err)
		return
	}
	err = faasInstance.Send(message)
	if err != nil {
		log.GetLogger().Errorf("websocket faas manager message faas instance send error: %v", err)
		return
	}

	res, _ := json.Marshal(message)
	w.Write(res)
}

func (m *WebsocketManager) Upper(ws *websocket.Conn) {
	var err error
	for {
		var event faascommon.Event

		if err = websocket.JSON.Receive(ws, &event); err != nil {
			log.GetLogger().Errorf("websocket faas manager recv error: %v ", err)
			break
		}

		if event.Type == "register" {
			var regMsg faascommon.RegisterEvent

			err  = json.Unmarshal([]byte(event.Message), &regMsg)
			if err != nil {
				log.GetLogger().Errorf("websocket faas manager json unmarshal register event error: %v", err)
				continue
			}

			faasSpec := faascommon.FaaSSpec{
				Id: regMsg.FaasId,
				Description: regMsg.Description,
			}

			faasInstance := NewWebsocketFaas(faasSpec, ws)
			err = m.Register(faasSpec.Id, faasInstance)
			if err != nil {
				log.GetLogger().Errorf("websocket faas manager register faas instance error: %v", err)
				continue
			}

			event.Type = "response"
			msg := faascommon.Response{
				Code: 0,
				Message: "faas instance " + faasSpec.Id + " register successful",
			}

			msgBytes, _ := json.Marshal(msg)
			event.Message = string(msgBytes)
			if err = faasInstance.Send(event); err != nil {
				log.GetLogger().Errorf("websocket faas manager register faas instance response error: %v", err)
				continue
			}
			continue
		} else if ( event.Type == "response") {
			continue
		} else {
			event.Type = "data"
			event.Message = strings.ToUpper(event.Message)
			if err = websocket.JSON.Send(ws, event); err != nil {
				fmt.Println(err)
				continue
			}
		}


	}
}

