// Package core provides the main entry point for the Hakoniwa game.
// This file contains the main game loop translated from Perl lib/Hako/Main.pm
//
// Ref: perl/lib/Hako/Main.pm:65-163
package core

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// RunMain is the main entry point for the game
// Ref: perl/lib/Hako/Main.pm:65
func RunMain(r *http.Request) {
	// Acquire lock
	if !hakoLock() {
		// Lock failed
		tempHeader()
		tempLockFail()
		tempFooter()
		return
	}
	defer unlock()

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Read cookie
	cookieInput(r)

	// Parse CGI input
	cgiInput(r)

	// Read island data
	currentID := variable.CurrentID
	if currentID == "" {
		currentID = "0"
	}
	if !ReadIslandsFile(currentID) {
		tempHeader()
		tempNoDataFile()
		tempFooter()
		return
	}

	// Mode dispatch
	// Phase 1: Simplified - just show error
	tempHeader()
	out("<H1>Phase 1: 実装中</H1>\n")
	out("Mode: " + variable.MainMode + "<BR>\n")
	tempFooter()

	// TODO: Implement mode handlers
	// switch variable.MainMode {
	// case "top":
	//     top.TopPageMain()
	// case "print":
	//     mapview.PrintIslandMain()
	// case "owner":
	//     mapview.OwnerMain()
	// case "command":
	//     mapview.CommandMain()
	// case "comment":
	//     mapview.CommentMain()
	// case "new":
	//     turn.NewIslandMain()
	// case "change":
	//     turn.ChangeMain()
	// case "Hdebugturn":
	//     turn.TurnMain()
	// case "lbbs":
	//     mapview.LocalBbsMain()
	// default:
	//     top.TopPageMain()
	// }
}

// cgiInput parses CGI form data
// Ref: perl/lib/Hako/Main.pm:454
func cgiInput(r *http.Request) {
	// Phase 1: Basic implementation
	if err := r.ParseForm(); err != nil {
		return
	}

	// Get mode from form
	variable.MainMode = r.FormValue("mode")
	if variable.MainMode == "" {
		variable.MainMode = "top"
	}

	// Get island ID
	variable.CurrentID = r.FormValue("ISLANDID")
	variable.CurrentName = r.FormValue("ISLANDNAME")

	// Get passwords
	variable.OldPassword = r.FormValue("OLDPASS")
	variable.InputPassword = r.FormValue("PASSWORD")
	variable.InputPassword2 = r.FormValue("PASSWORD2")

	// Get message
	variable.Message = r.FormValue("MESSAGE")

	// Get command parameters
	// TODO: Parse command parameters

	// Phase 1: Stub - set default mode
	if variable.MainMode == "" {
		variable.MainMode = "top"
	}
}

// cookieInput reads cookie values
// Ref: perl/lib/Hako/Main.pm:592
func cookieInput(r *http.Request) {
	// Phase 1: Basic implementation
	variable.DefaultPassword = ""

	cookie, err := r.Cookie("OWNISLANDID")
	if err == nil {
		variable.DefaultID = cookie.Value
	}

	cookie, err = r.Cookie("OWNISLANDPASSWORD")
	if err == nil {
		variable.DefaultPassword = cookie.Value
	}

	cookie, err = r.Cookie("TARGETISLANDID")
	if err == nil {
		variable.DefaultTarget = cookie.Value
	}

	// TODO: Read other cookie values
}

// cookieOutput sets cookie values
// Ref: perl/lib/Hako/Main.pm:627
func cookieOutput(w http.ResponseWriter) {
	// Phase 1: Basic implementation
	expires := time.Now().Add(30 * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "OWNISLANDID",
		Value:   variable.DefaultID,
		Expires: expires,
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "OWNISLANDPASSWORD",
		Value:   variable.DefaultPassword,
		Expires: expires,
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "TARGETISLANDID",
		Value:   variable.DefaultTarget,
		Expires: expires,
		Path:    "/",
	})

	// TODO: Set other cookie values
}
