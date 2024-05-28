package server

import (
	"context"
	"net"
	"net/http"

	"github.com/DimaGitHahahab/ozon-fintech-posts/internal/subscription"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Server is a server that handles GraphQL schema
type Server struct {
	server *http.Server
}

var schema graphql.Schema

// NewServer creates a new GraphQL server
func NewServer(s *graphql.Schema) *Server {
	schema = *s
	h := handler.New(&handler.Config{
		Schema:     s,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	mux := http.NewServeMux()
	mux.Handle("/root", h)

	subManager := subscription.New(schema)
	mux.HandleFunc("/subscriptions", subManager.SubscriptionsHandler)

	server := &http.Server{
		Handler: mux,
	}

	return &Server{
		server: server,
	}
}

// Run starts HTTP server on the given port
func (s *Server) Run(port string) error {
	s.server.Addr = net.JoinHostPort("", port)

	return s.server.ListenAndServe()
}

// Shutdown stops HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
