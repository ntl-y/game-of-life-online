package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

const port = "3000"

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(h http.Handler) error {
	s.httpServer = &http.Server{Addr: ":" + port, Handler: h}
	logrus.Printf("server started on port %s \n", port)
	return s.httpServer.ListenAndServe()
}
