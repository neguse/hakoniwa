// Package main provides the maintenance tool entry point.
// This is the HTTP server for Hakoniwa maintenance, equivalent to Perl's hako-mente.cgi
//
// Ref: perl/cgi/hako-mente.cgi
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neguse/hakoniwa/internal/hako/maintenance"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

func menteHandler(w http.ResponseWriter, r *http.Request) {
	// Clear output buffer
	variable.OutputBuffer.Reset()

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	// Run maintenance tool
	maintenance.RunMaintenance(r)

	// Write output
	w.Write(variable.OutputBuffer.Bytes())
}

func main() {
	http.HandleFunc("/", menteHandler)

	port := ":8081"
	fmt.Printf("Starting Hakoniwa maintenance server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
