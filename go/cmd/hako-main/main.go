// hako-main is the main entry point for the Hakoniwa game server.
// This file is translated from Perl cgi/hako-main.cgi
//
// Ref: perl/cgi/hako-main.cgi
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/neguse/hakoniwa/internal/web"
)

func main() {
	// Parse command line flags
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	log.Printf("箱庭諸島 ver2.30 (Go version)")
	log.Printf("Starting server on %s", *addr)

	// Setup routes
	mux := http.NewServeMux()

	// Main handler with middleware
	mainHandler := http.HandlerFunc(web.Handler)
	mux.Handle("/", web.Chain(mainHandler, web.LoggingMiddleware, web.RecoveryMiddleware))

	// Start server
	log.Fatal(http.ListenAndServe(*addr, mux))
}
