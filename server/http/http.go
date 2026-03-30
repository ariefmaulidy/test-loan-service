package httpserver

import (
	"net/http"

	loanhandler "github.com/ariefmaulidy/test-loan-service/handler/loan"
	"github.com/ariefmaulidy/test-loan-service/server"
)

type Server struct {
	mux *http.ServeMux
}

func New(deps *server.CommonDependency) *Server {
	mux := http.NewServeMux()

	loanH := loanhandler.New(deps.LoanUsecase)
	loanH.Register(mux)

	return &Server{mux: mux}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
