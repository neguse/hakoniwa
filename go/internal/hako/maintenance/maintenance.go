// Package maintenance provides maintenance tool functions for the Hakoniwa game.
// This file is translated from Perl lib/Hako/Maintenance.pm
//
// Ref: perl/lib/Hako/Maintenance.pm
package maintenance

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// Maintenance variables
var (
	mainMode  string
	inputPass string
	deleteID  string
	currentID string
	ctYear    int
	ctMon     int
	ctDate    int
	ctHour    int
	ctMin     int
	ctSec     int
)

// out outputs a string to the output buffer
func out(s string) {
	variable.OutputBuffer.WriteString(s)
}

// RunMaintenance runs the maintenance tool
// Ref: perl/lib/Hako/Maintenance.pm:70
func RunMaintenance(r *http.Request) {
	out(`Content-type: text/html

<HTML>
<HEAD>
<META http-equiv="Content-Type" content="text/html; charset=UTF-8">
<TITLE>箱島２ メンテナンスツール</TITLE>
</HEAD>
<BODY>
`)

	cgiInput(r)

	if mainMode == "delete" {
		if passCheck() {
			deleteMode()
		}
	} else if mainMode == "current" {
		if passCheck() {
			currentMode()
		}
	} else if mainMode == "time" {
		if passCheck() {
			timeMode()
		}
	} else if mainMode == "stime" {
		if passCheck() {
			stimeMode()
		}
	} else if mainMode == "new" {
		if passCheck() {
			newMode()
		}
	}

	mainModeDisplay()

	out(`</FORM>
</BODY>
</HTML>
`)
}

// myrmtree removes a directory and all its contents
// Ref: perl/lib/Hako/Maintenance.pm:120
func myrmtree(dirName string) error {
	dir, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		os.Remove(filepath.Join(dirName, file.Name()))
	}

	return os.Remove(dirName)
}

// currentMode restores a backup to current
// Ref: perl/lib/Hako/Maintenance.pm:131
func currentMode() {
	myrmtree(hconst.DirName)
	os.Mkdir(hconst.DirName, os.FileMode(hconst.DirMode))

	backupDir := hconst.DirName + ".bak" + currentID
	dir, err := os.Open(backupDir)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Name() != "." && file.Name() != ".." {
			fileCopy(
				filepath.Join(backupDir, file.Name()),
				filepath.Join(hconst.DirName, file.Name()),
			)
		}
	}
}

// deleteMode deletes current data or backup
// Ref: perl/lib/Hako/Maintenance.pm:143
func deleteMode() {
	if deleteID == "" {
		myrmtree(hconst.DirName)
	} else {
		myrmtree(hconst.DirName + ".bak" + deleteID)
	}
	os.Remove("hakojimalockflock")
}

// newMode creates new data directory
// Ref: perl/lib/Hako/Maintenance.pm:153
func newMode() {
	os.Mkdir(hconst.DirName, os.FileMode(hconst.DirMode))

	// Get current time aligned to unit time
	now := time.Now().Unix()
	now = now - (now % int64(hconst.UnitTime))

	// Create hakojima.dat
	filename := filepath.Join(hconst.DirName, "hakojima.dat")
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "1\n")        // Turn 1
	fmt.Fprintf(file, "%d\n", now)  // Start time
	fmt.Fprintf(file, "0\n")        // Island count
	fmt.Fprintf(file, "1\n")        // Next ID
}

// timeMode changes time using date/time fields
// Ref: perl/lib/Hako/Maintenance.pm:170
func timeMode() {
	ctMon--
	ctYear -= 1900

	// Convert to Unix timestamp
	t := time.Date(ctYear+1900, time.Month(ctMon+1), ctDate, ctHour, ctMin, ctSec, 0, time.Local)
	ctSec = int(t.Unix())

	stimeMode()
}

// stimeMode changes time using seconds
// Ref: perl/lib/Hako/Maintenance.pm:177
func stimeMode() {
	filename := filepath.Join(hconst.DirName, "hakojima.dat")
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	if len(lines) < 2 {
		return
	}

	// Update line 1 (0-indexed) with new time
	lines[1] = fmt.Sprintf("%d", ctSec)

	// Write back
	file, err = os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, line := range lines {
		fmt.Fprintf(file, "%s\n", line)
	}
}

// mainModeDisplay displays the main maintenance interface
// Ref: perl/lib/Hako/Maintenance.pm:191
func mainModeDisplay() {
	// Note: In Go version, this uses variable from hconst, not local thisFile
	// For Phase 1, we'll use a simplified version
	thisFile := "/hako-mente"

	out(fmt.Sprintf(`<FORM action="%s" method="POST">
<H1>箱島２ メンテナンスツール</H1>
<B>パスワード:</B><INPUT TYPE=password SIZE=32 MAXLENGTH=32 NAME=PASSWORD></TD>
`, thisFile))

	// Current data
	if _, err := os.Stat(hconst.DirName); err == nil {
		dataPrint("")
	} else {
		out(`    <HR>
    <INPUT TYPE="submit" VALUE="新しいデータを作る" NAME="NEW">
`)
	}

	// Backup data
	files, err := os.ReadDir("./")
	if err == nil {
		for _, file := range files {
			name := file.Name()
			if strings.HasPrefix(name, hconst.DirName+".bak") {
				suffix := strings.TrimPrefix(name, hconst.DirName+".bak")
				dataPrint(suffix)
			}
		}
	}
}

// dataPrint displays data information
// Ref: perl/lib/Hako/Maintenance.pm:222
func dataPrint(suf string) {
	out("<HR>")

	var filename string
	if suf == "" {
		filename = filepath.Join(hconst.DirName, "hakojima.dat")
		out("<H1>現役データ</H1>")
	} else {
		filename = filepath.Join(hconst.DirName+".bak"+suf, "hakojima.dat")
		out(fmt.Sprintf("<H1>バックアップ%s</H1>", suf))
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lastTurn string
	if scanner.Scan() {
		lastTurn = scanner.Text()
	}

	var lastTime int64
	if scanner.Scan() {
		lastTime, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	}

	timeString := timeToString(lastTime)

	out(fmt.Sprintf(`    <B>ターン%s</B><BR>
    <B>最終更新時間</B>:%s<BR>
    <B>最終更新時間(秒数表示)</B>:1970年1月1日から%d 秒<BR>
    <INPUT TYPE="submit" VALUE="このデータを削除" NAME="DELETE%s">
`, lastTurn, timeString, lastTime, suf))

	if suf == "" {
		t := time.Unix(lastTime, 0)
		year := t.Year()
		mon := int(t.Month())
		date := t.Day()
		hour := t.Hour()
		min := t.Minute()
		sec := t.Second()

		out(fmt.Sprintf(`    <H2>最終更新時間の変更</H2>
    <INPUT TYPE="text" SIZE=4 NAME="YEAR" VALUE="%d">年
    <INPUT TYPE="text" SIZE=2 NAME="MON" VALUE="%d">月
    <INPUT TYPE="text" SIZE=2 NAME="DATE" VALUE="%d">日
    <INPUT TYPE="text" SIZE=2 NAME="HOUR" VALUE="%d">時
    <INPUT TYPE="text" SIZE=2 NAME="MIN" VALUE="%d">分
    <INPUT TYPE="text" SIZE=2 NAME="NSEC" VALUE="%d">秒
    <INPUT TYPE="submit" VALUE="変更" NAME="NTIME"><BR>
    1970年1月1日から<INPUT TYPE="text" SIZE=32 NAME="SSEC" VALUE="%d">秒
    <INPUT TYPE="submit" VALUE="秒指定で変更" NAME="STIME">

`, year, mon, date, hour, min, sec, lastTime))
	} else {
		out(fmt.Sprintf(`	<INPUT TYPE="submit" VALUE="このデータを現役に" NAME="CURRENT%s">
`, suf))
	}
}

// timeToString formats a Unix timestamp
// Ref: perl/lib/Hako/Maintenance.pm:279
func timeToString(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return fmt.Sprintf("%d年 %d月 %d日 %d時 %d分 %d秒",
		t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// cgiInput parses CGI input
// Ref: perl/lib/Hako/Maintenance.pm:289
func cgiInput(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	mainMode = ""

	// Check for DELETE button
	for key := range r.Form {
		if strings.HasPrefix(key, "DELETE") {
			mainMode = "delete"
			deleteID = strings.TrimPrefix(key, "DELETE")
			break
		}
		if strings.HasPrefix(key, "CURRENT") {
			mainMode = "current"
			currentID = strings.TrimPrefix(key, "CURRENT")
			break
		}
	}

	if r.FormValue("NEW") != "" {
		mainMode = "new"
	} else if r.FormValue("NTIME") != "" {
		mainMode = "time"
		ctYear, _ = strconv.Atoi(r.FormValue("YEAR"))
		ctMon, _ = strconv.Atoi(r.FormValue("MON"))
		ctDate, _ = strconv.Atoi(r.FormValue("DATE"))
		ctHour, _ = strconv.Atoi(r.FormValue("HOUR"))
		ctMin, _ = strconv.Atoi(r.FormValue("MIN"))
		ctSec, _ = strconv.Atoi(r.FormValue("NSEC"))
	} else if r.FormValue("STIME") != "" {
		mainMode = "stime"
		ctSec, _ = strconv.Atoi(r.FormValue("SSEC"))
	}

	inputPass = r.FormValue("PASSWORD")
}

// fileCopy copies a file
// Ref: perl/lib/Hako/Maintenance.pm:345
func fileCopy(src, dist string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	distFile, err := os.Create(dist)
	if err != nil {
		return err
	}
	defer distFile.Close()

	_, err = io.Copy(distFile, srcFile)
	return err
}

// passCheck checks master password
// Ref: perl/lib/Hako/Maintenance.pm:357
func passCheck() bool {
	if inputPass == hconst.MasterPassword {
		return true
	}

	out(`   <FONT SIZE=7>パスワードが違います。</FONT>
`)
	return false
}
