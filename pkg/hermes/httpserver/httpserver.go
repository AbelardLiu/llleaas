package httpserver

import (
	"context"
	"lll.github.com/llleaas/cmd/hermes/app/options"
	"lll.github.com/llleaas/pkg/common/log"
	"net/http"
	"os"
	"strconv"
)

type HttpServer struct {
	Name string
	Ip string
	Port int32
	Option *options.HermesOption
	Server *http.Server
}

func NewHttpServer(name string, ip string, port int32, option *options.HermesOption) *HttpServer {
	return &HttpServer{
		Name: name,
		Ip: ip,
		Port: port,
		Option: option,
		Server: nil,
	}
}

func (s *HttpServer)Start(stopCh <- chan struct{}) error {
	handler := NewHttpHandler(s.Name, s.Option)
	handler.Start()

	listenAddress := s.Ip + ":" + strconv.Itoa(int(s.Port))
	log.GetLogger().Infof("start listening server %v", listenAddress)
	s.Server = &http.Server{
		Addr: listenAddress,
		Handler: handler,
		MaxHeaderBytes: 1 << 20,
	}

	go s.Stop(stopCh)

	if err := s.Server.ListenAndServe(); err != nil {
		log.GetLogger().Errorf("http server start listen and serve error: %v", err)
		return err
	}

	return nil
}

func (s *HttpServer)Stop(stopCh <- chan struct{}) error {
	select {
	case <-stopCh:
		log.GetLogger().Info("http server stop")
		s.Server.Shutdown(context.Background())
		os.Exit(0)
	}
	return nil
}
