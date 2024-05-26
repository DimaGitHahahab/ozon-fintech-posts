package server

import (
	"context"
	"net"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Server is a server that handles GraphQL schema
type Server struct {
	server *http.Server
}

// NewServer creates a new GraphQL server
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

// Run starts HTTP server on the given port
func (s *Server) Run(port string) error {
	s.server.Addr = net.JoinHostPort("", port)

	http.Handle("/root", s.server.Handler)

	return s.server.ListenAndServe()
}

// Shutdown stops HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
