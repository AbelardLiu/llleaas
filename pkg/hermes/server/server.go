package server

import (
	"lll.github.com/llleaas/cmd/hermes/app/options"
	"lll.github.com/llleaas/pkg/common/log"
	hermesHttp "lll.github.com/llleaas/pkg/hermes/httpserver"
	"os"
)

func Run(option *options.HermesOption, stopCh <- chan struct{}) error {
	log.UseStdOut()
	log.UseCaller()
	log.GetLogger().Info("server is started!")

	httpServer := hermesHttp.NewHttpServer("hermes-http-server", option.Ip, option.Port, option)

	go httpServer.Start(stopCh)

	if err := waitAndExit(stopCh); err != nil {
		log.GetLogger().Fatal("wait and exit error")
		os.Exit(1)
	}

	return nil
}

func waitAndExit(stopCh <- chan struct{}) error {
	select {
	case <- stopCh:
		log.GetLogger().Info("hermes exit! See you next time ^_^")
		os.Exit(0)
	}

	return nil
}
