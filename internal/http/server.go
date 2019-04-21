package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// Server is a structure that contains the pieces that make up a server.
type Server struct {
	Router      *chi.Mux
	httpsServer *http.Server
	httpServer  *http.Server
}

// Run kicks everything off by setting the routes and launching the HTTP
// and HTTPS servers in goroutines. It returns an error if something fails.
func (srv *Server) Run() error {
	// srv.setRoutes()

	errs := make(chan error, 1)

	go srv.newHTTPServer(errs)
	go srv.newHTTPSServer("localhost.pem", "localhost-key.pem", errs)

	return <-errs
}

// setRoutes configures the routing and loads any middleware.
// func (srv *Server) setRoutes() {
// 	router := chi.NewRouter()

// 	router.Use(middleware.RealIP)
// 	router.Use(middleware.Logger)
// 	router.Use(middleware.Recoverer)
// 	router.Use(middleware.DefaultCompress)

// 	// Non-API routes
// 	// router.Mount("/", (&PageHandler{}).Routes())
// 	router.Mount("/movies", srv.MovieHandler.Routes())

// 	// API (v1) routes
// 	// router.Route("/api/v1", func(r chi.Router) {
// 	// 	r.Mount("/movies", (&APIMoviesResource{}).Routes())
// 	// })
// }

// newHTTPServer configures and starts the HTTP server.
// It sends any errors via a channel if the server fails to start.
func (srv *Server) newHTTPServer(errs chan<- error) {
	srv.httpServer = &http.Server{
		Addr:         ":8080", // TODO(tim): Make this configurable
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.Host)
			if err != nil {
				errs <- err
			}
			url := r.URL
			url.Host = net.JoinHostPort(host, "8081") // TODO(tim): Make this configurable
			url.Scheme = "https"
			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		}),
	}

	errs <- srv.httpServer.ListenAndServe()
}

// newHTTPSServer configures and starts the HTTP server.
// It sends any errors via a channel if the server fails to start.
func (srv *Server) newHTTPSServer(certFile, keyFile string, errs chan<- error) {
	srv.httpsServer = &http.Server{
		Addr:         ":8081", // TODO(tim): Make this configurable
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      srv.Router,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			MinVersion:               tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	errs <- srv.httpsServer.ListenAndServeTLS(certFile, keyFile)
}
