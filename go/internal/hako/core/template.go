// Package core provides HTML template functions for the Hakoniwa game.
// This file contains template output functions translated from Perl lib/Hako/Main.pm
//
// Ref: perl/lib/Hako/Main.pm:1071-1175
package core

import (
	"fmt"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
)

// tempHeader outputs HTML header
// Ref: perl/lib/Hako/Main.pm:1100
func tempHeader() {
	out("Content-type: text/html; charset=UTF-8\n\n")
	out("<HTML>\n")
	out("<HEAD>\n")
	out(fmt.Sprintf("<TITLE>%s</TITLE>\n", hconst.Title))
	out("</HEAD>\n")
	out(fmt.Sprintf("<BODY %s TEXT=\"%s\">\n", hconst.HTMLBody, hconst.NormalColor))
}

// tempFooter outputs HTML footer
// Ref: perl/lib/Hako/Main.pm:1118
func tempFooter() {
	out("<HR>\n")
	out(hconst.TempBack)
	out("<BR>\n")
	out("</BODY></HTML>\n")
}

// TempLockFail outputs lock failure message (exported for external use)
// Ref: perl/lib/Hako/Main.pm:1133
func TempLockFail() {
	tempLockFail()
}

// tempLockFail outputs lock failure message (internal)
func tempLockFail() {
	out("<H1>ただいま、データ処理中です。</H1>\n")
	out("しばらくしてから、リロードしてください。<P>\n")
}

// TempUnlock outputs forced unlock message (exported for external use)
// Ref: perl/lib/Hako/Main.pm:1144
func TempUnlock() {
	tempUnlock()
}

// tempUnlock outputs forced unlock message (internal)
func tempUnlock() {
	out("<H1>強制ロック解除</H1>\n")
	out("前回のCGI実行が異常終了しました。ロックを強制解除しました。<P>\n")
}

// TempNoDataFile outputs no data file message (exported for external use)
// Ref: perl/lib/Hako/Main.pm:1154
func TempNoDataFile() {
	tempNoDataFile()
}

// tempNoDataFile outputs no data file message (internal)
func tempNoDataFile() {
	out("<H1>データファイルが存在しません。</H1>\n")
}

// TempInitialize initializes template system (exported for external use)
// Ref: perl/lib/Hako/Main.pm (tempInitialize)
func TempInitialize() {
	// Phase 1: No special initialization needed
	// Template variables would be initialized here in Phase 2
}

// TempWrongPassword outputs wrong password message (exported for external use)
// Ref: perl/lib/Hako/Main.pm:1161
func TempWrongPassword() {
	out("<H1>パスワードが違います。</H1>\n")
}

// TempProblem outputs generic problem message (exported for external use)
// Ref: perl/lib/Hako/Main.pm:1168
func TempProblem() {
	out("<H1>何か問題が発生しました。</H1>\n")
	out("しばらくしてから、リロードしてください。<P>\n")
}
