package common

type FaaSSpec struct {
	Id string
	Description string
}

type FaaSInstance interface {
	Send(event Event) error
	Recv() (Event, error)
	Info() string
	Spec() FaaSSpec
}
