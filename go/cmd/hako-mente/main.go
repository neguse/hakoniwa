// hako-mente is the maintenance tool entry point for the Hakoniwa game.
// This file is translated from Perl cgi/hako-mente.cgi
//
// Ref: perl/cgi/hako-mente.cgi
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/neguse/hakoniwa/internal/web"
)

func main() {
	// Parse command line flags
	addr := flag.String("addr", ":8081", "HTTP server address")
	flag.Parse()

	log.Printf("箱庭諸島 Maintenance Tool (Go version)")
	log.Printf("Starting maintenance server on %s", *addr)

	// Setup routes
	mux := http.NewServeMux()

	// Maintenance handler with middleware
	menteHandler := http.HandlerFunc(web.MaintenanceHandler)
	mux.Handle("/", web.Chain(menteHandler, web.LoggingMiddleware, web.RecoveryMiddleware))

	// Start server
	log.Fatal(http.ListenAndServe(*addr, mux))
}
