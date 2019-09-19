package common


type FaasResource struct {
	Type string `json:"type, omitempty"`
	Value uint64 `json:"value, omitempty"`
	Unit string `json:"unit, omitempty""`
}

type FaaSSpec struct {
	Id string		`json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	Platform string `json:"platform, omitempty"`
	Resources []FaasResource `json:"resources, omitempty"`
}

type FaaSInstance interface {
	Send(event Event) error
	Recv() (Event, error)
	Info() string
	Spec() FaaSSpec
	Status() string
}

type FaaSGetter interface {
	Get(id string) (FaaSInstance, error)
}
