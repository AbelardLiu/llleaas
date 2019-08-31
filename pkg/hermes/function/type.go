package function

import "net/http"

type FunctionManager interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Deploy(w http.ResponseWriter, r *http.Request)
	Invoke(w http.ResponseWriter, r *http.Request)
}
