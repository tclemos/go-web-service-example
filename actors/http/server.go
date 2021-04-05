package http

import (
	"github.com/labstack/echo/v4"
	"github.com/tclemos/go-web-service-example/actors/http/api"
	"github.com/tclemos/go-web-service-example/actors/http/controllers"
)

type Server struct {
	echo *echo.Echo
	*controllers.ThingsController
}

func NewServer(tc *controllers.ThingsController) *Server {
	s := &Server{
		echo:             echo.New(),
		ThingsController: tc,
	}

	api.RegisterHandlers(s.echo, s)
	return s
}

func (s *Server) Start() {
	s.echo.Server.ListenAndServe()
}
