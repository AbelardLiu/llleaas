package options

type HermesOption struct {
	Name string
	Ip   string
	Port int32
	DefaultFaasUrl string
	Log2std bool
	LogLevel string
}

func NewHermesOption() *HermesOption {
	return &HermesOption{
		Name:     "",
		Ip:       "",
		Port:     0,
		DefaultFaasUrl: "",
		Log2std:  false,
		LogLevel: "",
	}
}
