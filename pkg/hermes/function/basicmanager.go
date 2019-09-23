package function

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"lll.github.com/llleaas/cmd/hermes/app/options"
	"lll.github.com/llleaas/pkg/common/log"
	"lll.github.com/llleaas/pkg/hermes/common"
	"net/http"
	"sync"
)

type DeploySpec struct {
	FunctionName string `json:"functionName"`
	FunctionVersion string `json:"functionversion"`
	FunctionRuntime string `json:"functionRuntime"`
	FaaSId string `json:"faasId"`
}

type InvokeSpec struct {
	FunctionName string `json:"functionName"`
	FunctionVersion string `json:"functionversion"`
	FaaSId string `json:"faasId"`
	Request string `json:"request, omitempty"`
}

type BasicFunctionManager struct {
	Option *options.HermesOption
	Getter common.FaaSGetter
	FunctionMux sync.Mutex
	Functions map[string] common.FaaSSpec   // functionName-functionVersion: faasInstance
}

func NewBasicFunctionManager(getter common.FaaSGetter, option *options.HermesOption) *BasicFunctionManager {
	return &BasicFunctionManager{
		Getter:getter,
		FunctionMux: sync.Mutex{},
		Functions: make(map[string] common.FaaSSpec),
		Option: option,
	}
}

func (m *BasicFunctionManager) RegisterHandler(basepath string, r *mux.Router) error {
	r.HandleFunc(basepath + "/function/{name}/version/{version}", m.Get).Methods(http.MethodGet)
	r.HandleFunc(basepath + "/function/", m.List).Methods(http.MethodGet)
	r.HandleFunc(basepath + "/function/deploy/{name}/version/{version}", m.Deploy).Methods(http.MethodPost)
	r.HandleFunc(basepath + "/function/invoke/{name}/version/{version}", m.Invoke).Methods(http.MethodPost)

	return nil
}

func (m *BasicFunctionManager) Get(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (m *BasicFunctionManager) List(w http.ResponseWriter, r *http.Request) {


}

func (m *BasicFunctionManager) Deploy(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	funcName := "unknown"
	funcVersion := "unknown"
	if val, ok := params["name"]; ok {
		funcName = val
	}
	if val, ok := params["version"]; ok {
		funcVersion = val
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.GetLogger().Errorf("basic function manager deploy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var deploySpec DeploySpec
	err = json.Unmarshal(body, &deploySpec)

	if ( deploySpec.FunctionName != funcName || deploySpec.FunctionVersion != funcVersion ) {
		log.GetLogger().Errorf("basic function manager function check error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	faasId := deploySpec.FaaSId
	functionRuntime := deploySpec.FunctionRuntime

	faasInstance, err := m.Getter.Get(faasId)
	if err != nil {
		log.GetLogger().Errorf("basic function manager get faas instance error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := m.makeDeployEvent(faasId, funcName, funcVersion, functionRuntime)
	if err != nil {
		log.GetLogger().Errorf("basic function manager make deploy event error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = faasInstance.Send(event)
	if err != nil {
		log.GetLogger().Errorf("basic function manager send event error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: recv response

	// update record function - faas map
	m.FunctionMux.Lock()
	defer m.FunctionMux.Unlock()

	functionId := funcName + "-" + funcVersion
	m.Functions[functionId] = faasInstance.Spec()
}

func (m *BasicFunctionManager) Invoke(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	funcName := "unknown"
	funcVersion := "unknown"
	if val, ok := params["name"]; ok {
		funcName = val
	}
	if val, ok := params["version"]; ok {
		funcVersion = val
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.GetLogger().Errorf("basic function manager deploy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var invokeSpec InvokeSpec
	err = json.Unmarshal(body, &invokeSpec)

	if ( invokeSpec.FunctionName != funcName || invokeSpec.FunctionVersion != funcVersion ) {
		log.GetLogger().Errorf("basic function manager function check error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	faasId := invokeSpec.FaaSId

	faasInstance, err := m.Getter.Get(faasId)
	if err != nil {
		log.GetLogger().Errorf("basic function manager get faas instance error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event, err := m.makeInvokeEvent(faasId, funcName, funcVersion, faasInstance.Spec().Platform, invokeSpec.Request)
	if err != nil {
		log.GetLogger().Errorf("basic function manager make deploy event error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = faasInstance.Send(event)
	if err != nil {
		log.GetLogger().Errorf("basic function manager send event error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: recv response

}

func (m *BasicFunctionManager)makeDeployEvent(faasId string, functionName string, functionVersion string, functionRuntime string) (common.Event, error) {
	address := "http://122.112.130.166:51010/function/code-zip/" + functionName + "/version/" + functionVersion
	deployMessage := common.DeployMessage{
		Function: functionName,
		Version: functionVersion,
		Address: address,
		Runtime: functionRuntime,
	}

	deployEvent := common.DeployFunctionEvent{
		FaasId: faasId,
		Type: "deployment",
		Message: deployMessage,
	}

	deployEventData, err := json.Marshal(deployEvent)
	if err != nil {
		log.GetLogger().Errorf("basic function manager json marshal deploy event error: %v", err)
		return common.Event{}, err
	}

	event := common.Event{
		Type: "deploy",
		Message: string(deployEventData),
	}

	return event, nil
}

func (m *BasicFunctionManager)makeInvokeEvent(faasId string, functionName string, functionVersion string,  runtime string, request string) (common.Event, error) {
	address := "http://122.112.130.166:51010/function/code-zip/" + functionName + "/version/" + functionVersion
	invokeMessage := common.InvokeMessage{
		Function: functionName,
		Version: functionVersion,
		Address: address,
		Runtime: runtime,
		Request: []byte(request),
	}

	log.GetLogger().Infof("invoke message %v", invokeMessage)

	invokeEvent := common.InvokeFunctionEvent{
		FaasId: faasId,
		Type: "invoke",
		Message: invokeMessage,
	}

	invokeEventData, err := json.Marshal(invokeEvent)
	if err != nil {
		log.GetLogger().Errorf("basic function manager json marshal deploy event error: %v", err)
		return common.Event{}, err
	}

	event := common.Event{
		Type: "invoke",
		Message: string(invokeEventData),
	}

	return event, nil
}


