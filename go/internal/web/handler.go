// Package web provides HTTP handlers for the Hakoniwa game.
// This file is translated from Perl lib/Hako/Main.pm (run_main function)
//
// Ref: perl/lib/Hako/Main.pm
package web

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/neguse/hakoniwa/internal/hako/core"
	"github.com/neguse/hakoniwa/internal/hako/maintenance"
	"github.com/neguse/hakoniwa/internal/hako/mapview"
	"github.com/neguse/hakoniwa/internal/hako/top"
	"github.com/neguse/hakoniwa/internal/hako/turn"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// Handler handles HTTP requests for the Hakoniwa game
// Ref: perl/lib/Hako/Main.pm:64 (run_main)
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Lock the data files
	if !core.HakoLock() {
		// Lock failed
		tempHeader(w)
		core.TempLockFail()
		tempFooter()
		writeOutput(w)
		return
	}

	// Parse cookies
	cookieInput(r)

	// Parse CGI parameters
	cgiInput(r)

	// Read island data
	if !core.ReadIslandsFile(variable.CurrentID) {
		core.Unlock()
		tempHeader(w)
		core.TempNoDataFile()
		tempFooter()
		writeOutput(w)
		return
	}

	// Initialize template (clear output buffer)
	variable.OutputBuffer.Reset()
	core.TempInitialize()

	// Output cookie
	cookieOutput(w)

	// Output header
	tempHeader(w)

	// Route by main mode
	switch variable.MainMode {
	case "turn":
		// Turn progression
		turn.TurnMain()

	case "new":
		// New island creation
		turn.NewIslandMain()

	case "print":
		// Tourist mode
		mapview.PrintIslandMain()

	case "owner":
		// Development mode
		mapview.OwnerMain()

	case "command":
		// Command input mode
		mapview.CommandMain()

	case "comment":
		// Comment input mode
		mapview.CommentMain()

	case "lbbs":
		// Local BBS mode
		mapview.LocalBbsMain()

	case "change":
		// Information change mode
		turn.ChangeMain()

	default:
		// Default: Top page mode
		top.TopPageMain()
	}

	// Output footer
	tempFooter()

	// Write output buffer to response
	writeOutput(w)
}

// MaintenanceHandler handles maintenance tool requests
// Ref: perl/cgi/hako-mente.cgi
func MaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Parse CGI parameters
	cgiInput(r)

	// Initialize template (clear output buffer)
	variable.OutputBuffer.Reset()
	core.TempInitialize()

	// Output header
	tempHeader(w)

	// Run maintenance mode
	maintenance.RunMaintenance(r)

	// Output footer
	tempFooter()

	// Write output buffer to response
	writeOutput(w)
}

// cookieInput reads cookies from the request
// Ref: perl/lib/Hako/Main.pm (cookieInput)
func cookieInput(r *http.Request) {
	// Phase 1: Simplified - read ISLANDID cookie
	if cookie, err := r.Cookie("ISLANDID"); err == nil {
		variable.DefaultID = cookie.Value
	}

	if cookie, err := r.Cookie("DEVELOPPE"); err == nil {
		variable.DefaultPassword = cookie.Value
	}
}

// cookieOutput writes cookies to the response
// Ref: perl/lib/Hako/Main.pm (cookieOutput)
func cookieOutput(w http.ResponseWriter) {
	// Phase 1: Simplified - set ISLANDID cookie if island was accessed
	if variable.CurrentID != "" && variable.CurrentID != variable.DefaultID {
		http.SetCookie(w, &http.Cookie{
			Name:     "ISLANDID",
			Value:    variable.CurrentID,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365, // 1 year
			HttpOnly: true,
		})
	}

	// Set development password cookie if provided
	if variable.InputPassword != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "DEVELOPPE",
			Value:    variable.InputPassword,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365, // 1 year
			HttpOnly: true,
		})
	}
}

// cgiInput reads CGI parameters from the request
// Ref: perl/lib/Hako/Main.pm (cgiInput)
func cgiInput(r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return
	}

	// Read main mode
	variable.MainMode = r.FormValue("mode")

	// Read island ID
	if id := r.FormValue("ISLANDID"); id != "" {
		variable.CurrentID = id
	} else {
		variable.CurrentID = variable.DefaultID
	}

	// Read island name
	variable.CurrentName = r.FormValue("ISLANDNAME")

	// Read passwords
	variable.OldPassword = r.FormValue("OLDPASS")
	variable.InputPassword = r.FormValue("PASSWORD")
	variable.InputPassword2 = r.FormValue("PASSWORD2")

	// Read comment
	variable.InputComment = r.FormValue("COMMENT")

	// Read BBS inputs
	variable.InputBbsName = r.FormValue("BBSNAME")
	variable.InputBbsMessage = r.FormValue("MESSAGE")

	// Read command inputs
	// Phase 1: Simplified command parsing
	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("COMKIND%d", i)
		if r.FormValue(key) != "" {
			// Store command data in variables
			// Full implementation would parse all command parameters
		}
	}
}

// tempHeader outputs HTTP headers and HTML header
// Ref: perl/lib/Hako/Main.pm (tempHeader)
func tempHeader(w http.ResponseWriter) {
	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache")

	// Output HTML header (will be written later by writeOutput)
	variable.OutputBuffer.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>箱庭諸島</title>
</head>
<body>
`)
}

// tempFooter outputs HTML footer
// Ref: perl/lib/Hako/Main.pm (tempFooter)
func tempFooter() {
	variable.OutputBuffer.WriteString(`
</body>
</html>
`)
}

// writeOutput writes the accumulated output buffer to the response
func writeOutput(w http.ResponseWriter) {
	io.WriteString(w, variable.OutputBuffer.String())
}
