package function

import (
	"github.com/gorilla/mux"
	"lll.github.com/llleaas/pkg/hermes/common"
	"net/http"
)

type BasicFunctionManager struct {
	Functions map[string] common.FaaSInstance
}

func NewBasicFunctionManager() *BasicFunctionManager {
	return &BasicFunctionManager{
		Functions: make(map[string] common.FaaSInstance),
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
	panic("implement me")
}

func (m *BasicFunctionManager) Deploy(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (m *BasicFunctionManager) Invoke(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}




