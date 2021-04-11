package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tclemos/go-web-service-example/adapters/http/api"
	"github.com/tclemos/go-web-service-example/adapters/http/controllers"
	"github.com/tclemos/go-web-service-example/adapters/logger"
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
	host := os.Getenv("THING_APP_HTTP_SERVER_HOST")
	port := os.Getenv("THING_APP_HTTP_SERVER_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	logger.Infof("Server address: %s. Listening...", address)
	s.echo.Start(address)
}

// (GET /ping)
func (s *Server) Ping(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "pong")
}
