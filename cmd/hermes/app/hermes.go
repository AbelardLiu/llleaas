package app

import (
	"fmt"
	"lll.github.com/llleaas/pkg/hermes/server"
	"os"

	"github.com/spf13/cobra"

	"lll.github.com/llleaas/cmd/hermes/app/options"
)

func NewHermesCommand(stopCh <- chan struct{}) *cobra.Command {
	s := options.NewHermesOption()

	cmd := &cobra.Command{
		Use:                        "Hermes service",
		Long:                       "Hermes service is a service for passing message between all things in llleaas system",
		RunE:                        func (cmd *cobra.Command, args []string) error {
			fmt.Println("starting hermes service!")
			return server.Run(s, stopCh)
		},
	}

	cmd.SetArgs(os.Args)

	fs := cmd.Flags()

	fs.StringP("ip", "", "0.0.0.0", "The hermes listen ip")
	fs.Int32P("port", "", 51010, "The listen port of hermes")
	fs.StringP("default-faas-url", "", "http://192.168.0.192:31000", "The default faas url of hermes")
	//fs.StringP("default-faas-url", "", "http://127.0.0.1:8500", "The default faas url of hermes")
	fs.StringP("name", "n", "hermes", "The instance name of service")
	fs.BoolP("log2std", "", true, "log to std out, used for debugging")
	fs.StringP("loglevel", "", "info", "Log Level('debug','info', 'warn', 'fatal', 'trace')")

	fs.Parse(os.Args)
	s.Name, _ = fs.GetString("name")
	s.Ip, _   = fs.GetString("ip")
	s.Port, _ = fs.GetInt32("port")
	s.DefaultFaasUrl, _ = fs.GetString("default-faas-url")
	s.Log2std, _ = fs.GetBool("log2std")
	s.LogLevel, _ = fs.GetString("loglevel")

	fmt.Print(s)

	return cmd
}

