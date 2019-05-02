// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package offersweb

import (
	"context"
	"net"
	"net/http"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Error is satellite referral error type
var Error = errs.Class("satellite referral error")

// Config contains configuration for referral offersweb server
type Config struct {
	Address string `help:"server address of the frontend app" default:"127.0.0.1:8090"`
}

// Server represents referral offersweb server
type Server struct {
	log *zap.Logger

	config Config
	// service *referral.Service

	listener net.Listener
	server   http.Server
}

// NewServer creates new instance of offersweb server
func NewServer(logger *zap.Logger, config Config, listener net.Listener) *Server {
	server := Server{
		log:      logger,
		config:   config,
		listener: listener,
		// service:  service,
	}

	logger.Debug("Starting offersweb UI...")

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(server.appHandler))

	server.server = http.Server{
		Handler: mux,
	}

	return &server
}

// appHandler is web app http handler function
func (s *Server) appHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("This is a test"))
	// http.ServeFile(w, req, filepath.Join(s.config.StaticDir, "dist", "public", "index.html"))
}

func (s *Server) serveError(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("bad request"))
}

// Run starts the server that host admin web app and api endpoint
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return s.server.Shutdown(nil)
	})
	group.Go(func() error {
		defer cancel()
		return s.server.Serve(s.listener)
	})

	return group.Wait()
}

// Close closes server and underlying listener
func (s *Server) Close() error {
	return s.server.Close()
}