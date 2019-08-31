package faasmanager

import "net/http"


type FaaSManager interface {
	//Start(basepath ) error
	ServeHTTP(w http.ResponseWriter,r *http.Request)

	//Register(id string, faas FaaSInstance) error
	//UnRegister(id string) error
	//
	//Get(id string) (FaaSInstance, error)
	//List() map[string]FaaSInstance
}