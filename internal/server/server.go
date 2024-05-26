package server

import (
	"context"
	"net"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Server struct {
	server *http.Server
}

func NewServer(schema *graphql.Schema) *Server {
	server := &http.Server{
		Handler: handler.New(&handler.Config{
			Schema:   schema,
			Pretty:   true,
			GraphiQL: true,
		}),
	}

	return &Server{
		server: server,
	}
}

func (s *Server) Run(port string) error {
	s.server.Addr = net.JoinHostPort("", port)

	http.Handle("/root", s.server.Handler)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
