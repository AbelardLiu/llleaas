package options

type HermesOption struct {
	Name string
	Ip   string
	Port int32
	Log2std bool
	LogLevel string
}

func NewHermesOption() *HermesOption {
	return &HermesOption{
		Name:     "",
		Ip:       "",
		Port:     0,
		Log2std:  false,
		LogLevel: "",
	}
}
