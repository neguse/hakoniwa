// Package main provides the main game server entry point.
// This is the HTTP server for Hakoniwa game, equivalent to Perl's hako-main.cgi
//
// Ref: perl/cgi/hako-main.cgi
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neguse/hakoniwa/internal/hako/core"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Clear output buffer
	variable.OutputBuffer.Reset()

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	// Run main game logic
	core.RunMain(r)

	// Write output
	w.Write(variable.OutputBuffer.Bytes())
}

func main() {
	http.HandleFunc("/", mainHandler)

	port := ":8080"
	fmt.Printf("Starting Hakoniwa server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
