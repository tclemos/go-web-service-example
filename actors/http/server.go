package http

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) AddController(c Controller) error {
	return nil
}

func (s *Server) Start() {

}
