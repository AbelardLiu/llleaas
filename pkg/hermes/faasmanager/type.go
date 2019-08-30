package faasmanager

import "net/http"

type Event struct {
	Type string	`json:"type"`
	Message string `json:"message"`
}

type RegisterEvent struct {
	FaasId string `json:"faas_id"`
	Description string `json:"description"`
	NodeType string `json:"node_type"`
}

type DeployMessage struct {
	Function string `json:"function"`
	Version string `json:"version"`
	Address string `json:"address"`  // need to delete
	Runtime string `json:"runtime"`
}

type DeployFunctionEvent struct {
	FaasId string `json:"faas_id"`
	Type string `json:"type"`
	Message DeployMessage `json:"message"`
}

type InvokeMessage struct {
	Function string `json:"function"`
	Version string `json:"version"`
	Address string `json:"address"`  // need to delete
	Runtime string `json:"runtime"`
	Request []byte `json:"request"`
}

type InvokeFunctionEvent struct {
	FaasId string `json:"faas_id"`
	Type string `json:"type"`
	Message InvokeMessage `json:"message"`
}

type Response struct {
	Code int `json:"code""`
	Message string `json:"message"`
}

type FaaSInstance interface {
 	Send(event Event) error
	Recv() (Event, error)
 	Info() string
}

type FaaSManager interface {
	//Start(basepath ) error
	ServeHTTP(w http.ResponseWriter,r *http.Request)

	//Register(id string, faas FaaSInstance) error
	//UnRegister(id string) error
	//
	//Get(id string) (FaaSInstance, error)
	//List() map[string]FaaSInstance
}