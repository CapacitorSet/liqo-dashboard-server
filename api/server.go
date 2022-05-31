package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type APIServer struct {
	client.Client

	EchoServer    *echo.Echo
	ListenAddress string
}

func (s APIServer) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.EchoServer.Shutdown(ctx)
	}()
	return s.EchoServer.Start(s.ListenAddress)
}

func NewAPIServer(k8sClient client.Client, address string) *APIServer {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	srv := &APIServer{Client: k8sClient, EchoServer: e, ListenAddress: address}
	RegisterHandlers(e, srv)

	return srv
}
