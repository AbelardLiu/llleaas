package common

type Event struct {
	Type string	`json:"type"`
	Message string `json:"message"`
}

type RegisterEvent struct {
	FaasId string `json:"faas_id,omitempty"`
	Description string `json:"description,omitempty"`
	NodeType string `json:"node_type,omitempty"`
	Spec FaaSSpec `json:"spec, omitempty"`
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
