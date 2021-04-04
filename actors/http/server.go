package http

import "github.com/labstack/echo/v4"

type Server struct {
	echo *echo.Echo
	*ThingsController
}

func NewServer(tc *ThingsController) *Server {
	s := &Server{
		echo:             echo.New(),
		ThingsController: tc,
	}

	RegisterHandlers(s.echo, s)
	return s
}

func (s *Server) Start() {
	s.echo.Server.ListenAndServe()
}

func httpError(ctx echo.Context, code int, m string) error {
	err := ctx.JSON(code, Error{
		Message: m,
	})

	return err
}
