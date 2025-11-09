// Package mapview provides map viewing and island management functions for the Hakoniwa game.
// This file is translated from Perl lib/Hako/Map.pm
//
// Ref: perl/lib/Hako/Map.pm
package mapview

import (
	"fmt"

	"github.com/neguse/hakoniwa/internal/hako/core"
	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/types"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// out outputs a string to the output buffer
func out(s string) {
	variable.OutputBuffer.WriteString(s)
}

// PrintIslandMain displays the island viewing page
// Ref: perl/lib/Hako/Map.pm:38
func PrintIslandMain() {
	// Release lock
	core.Unlock()

	// Get island number from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num

	// Get island name
	variable.CurrentName = variable.Islands[variable.CurrentNumber].Name

	// Display tourist view
	tempPrintIslandHead()
	islandInfo()
	islandMap(0)

	// Local BBS
	if hconst.UseLbbs != 0 {
		tempLbbsHead()
		tempLbbsInput()
		tempLbbsContents()
	}

	// Recent events
	tempRecent(0)
}

// OwnerMain displays the island development page
// Ref: perl/lib/Hako/Map.pm:75
func OwnerMain() {
	// Release lock
	core.Unlock()

	// Set mode
	variable.MainMode = "owner"

	// Get island from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num
	island := variable.Islands[variable.CurrentNumber]
	variable.CurrentName = island.Name

	// Check password
	if !core.CheckPassword(island.Password, variable.InputPassword) {
		core.TempWrongPassword()
		return
	}

	// Development view
	tempOwner()

	// Local BBS
	if hconst.UseLbbs != 0 {
		tempLbbsHead()
		tempLbbsInputOW()
		tempLbbsContents()
	}

	// Recent events
	tempRecent(1)
}

// CommandMain processes command input
// Ref: perl/lib/Hako/Map.pm:114
func CommandMain() {
	// Get island from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.Unlock()
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num
	island := variable.Islands[variable.CurrentNumber]
	variable.CurrentName = island.Name

	// Check password
	if !core.CheckPassword(island.Password, variable.InputPassword) {
		core.Unlock()
		core.TempWrongPassword()
		return
	}

	// Phase 1: Simplified command processing
	// Handle basic command operations
	if variable.CommandMode == "delete" {
		// Delete command
		slideFront(island.Commands, variable.CommandPlanNumber)
		tempCommandDelete()
	} else {
		// Add/insert command
		if variable.CommandMode == "insert" {
			slideBack(island.Commands, variable.CommandPlanNumber)
		}
		tempCommandAdd()

		// Register command
		if variable.CommandPlanNumber >= 0 && variable.CommandPlanNumber < len(island.Commands) {
			island.Commands[variable.CommandPlanNumber] = types.Command{
				Kind:   variable.CommandKind,
				Target: variable.CommandTarget,
				X:      variable.CommandX,
				Y:      variable.CommandY,
				Arg:    variable.CommandArg,
			}
		}
	}

	// Write data
	core.WriteIslandsFile()

	// Go to owner mode
	OwnerMain()
}

// CommentMain processes comment input
// Ref: perl/lib/Hako/Map.pm:208
func CommentMain() {
	// Get island from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.Unlock()
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num
	island := variable.Islands[variable.CurrentNumber]
	variable.CurrentName = island.Name

	// Check password
	if !core.CheckPassword(island.Password, variable.InputPassword) {
		core.Unlock()
		core.TempWrongPassword()
		return
	}

	// Update comment
	island.Comment = core.HtmlEscape(variable.Message)

	// Write data
	core.WriteIslandsFile()

	// Comment update message
	tempComment()

	// Go to owner mode
	OwnerMain()
}

// LocalBbsMain processes local BBS operations
// Ref: perl/lib/Hako/Map.pm:242
func LocalBbsMain() {
	// Get island number from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.Unlock()
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num
	island := variable.Islands[variable.CurrentNumber]

	// Check for missing name/message (except delete mode)
	if variable.LbbsMode != 2 {
		if variable.LbbsName == "" || variable.LbbsMessage == "" {
			core.Unlock()
			tempLbbsNoMessage()
			return
		}
	}

	// Check password for non-tourist mode
	if variable.LbbsMode != 0 {
		if !core.CheckPassword(island.Password, variable.InputPassword) {
			core.Unlock()
			core.TempWrongPassword()
			return
		}
	}

	lbbs := island.Lbbs

	// Mode handling
	if variable.LbbsMode == 2 {
		// Delete mode
		slideBackLbbsMessage(&lbbs, variable.CommandPlanNumber)
		tempLbbsDelete()
	} else {
		// Add mode
		slideLbbsMessage(&lbbs)

		// Write message
		var modeFlag string
		if variable.LbbsMode == 0 {
			modeFlag = "0"
		} else {
			modeFlag = "1"
		}

		variable.LbbsName = fmt.Sprintf("%d：%s", variable.IslandTurn, core.HtmlEscape(variable.LbbsName))
		variable.LbbsMessage = core.HtmlEscape(variable.LbbsMessage)
		if len(lbbs) > 0 {
			lbbs[0] = types.LbbsEntry{
				Message: fmt.Sprintf("%s>%s>%s", modeFlag, variable.LbbsName, variable.LbbsMessage),
			}
		}

		tempLbbsAdd()
	}

	// Update island lbbs
	island.Lbbs = lbbs

	// Write data
	core.WriteIslandsFile()

	// Return to original mode
	if variable.LbbsMode == 0 {
		PrintIslandMain()
	} else {
		OwnerMain()
	}
}

// Helper functions for command/lbbs array manipulation

func slideFront(commands []types.Command, number int) {
	if number >= 0 && number < len(commands)-1 {
		copy(commands[number:], commands[number+1:])
		commands[len(commands)-1] = types.Command{Kind: 0}
	}
}

func slideBack(commands []types.Command, number int) {
	if number >= 0 && number < len(commands)-1 {
		copy(commands[number+1:], commands[number:])
	}
}

func slideLbbsMessage(lbbs *[]types.LbbsEntry) {
	if len(*lbbs) > 0 {
		// Remove last, insert at beginning
		*lbbs = (*lbbs)[:len(*lbbs)-1]
		*lbbs = append([]types.LbbsEntry{{Message: ""}}, *lbbs...)
	}
}

func slideBackLbbsMessage(lbbs *[]types.LbbsEntry, number int) {
	if number >= 0 && number < len(*lbbs) {
		// Remove element at number
		*lbbs = append((*lbbs)[:number], (*lbbs)[number+1:]...)
		// Add empty at end
		*lbbs = append(*lbbs, types.LbbsEntry{Message: "0>>"})
	}
}

// islandInfo displays island information table
// Ref: perl/lib/Hako/Map.pm:342
func islandInfo() {
	island := variable.Islands[variable.CurrentNumber]

	// Display info
	rank := variable.CurrentNumber + 1
	farm := island.Farm
	factory := island.Factory
	mountain := island.Mountain

	farmStr := "保有せず"
	if farm != 0 {
		farmStr = fmt.Sprintf("%d0%s", farm, hconst.UnitPop)
	}
	factoryStr := "保有せず"
	if factory != 0 {
		factoryStr = fmt.Sprintf("%d0%s", factory, hconst.UnitPop)
	}
	mountainStr := "保有せず"
	if mountain != 0 {
		mountainStr = fmt.Sprintf("%d0%s", mountain, hconst.UnitPop)
	}

	mStr1 := ""
	mStr2 := ""
	if hconst.HideMoneyMode == 1 || variable.MainMode == "owner" {
		mStr1 = fmt.Sprintf("<TH %s nowrap=nowrap><NOBR>%s資金%s</NOBR></TH>",
			hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd)
		mStr2 = fmt.Sprintf("<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>",
			hconst.BgInfoCell, island.Money, hconst.UnitMoney)
	} else if hconst.HideMoneyMode == 2 {
		mTmp := core.AboutMoney(island.Money)
		mStr1 = fmt.Sprintf("<TH %s nowrap=nowrap><NOBR>%s資金%s</NOBR></TH>",
			hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd)
		mStr2 = fmt.Sprintf("<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>",
			hconst.BgInfoCell, mTmp)
	}

	out(fmt.Sprintf(`<CENTER>
<TABLE BORDER>
<TR>
<TH %s nowrap=nowrap><NOBR>%s順位%s</NOBR></TH>
<TH %s nowrap=nowrap><NOBR>%s人口%s</NOBR></TH>
%s
<TH %s nowrap=nowrap><NOBR>%s食料%s</NOBR></TH>
<TH %s nowrap=nowrap><NOBR>%s面積%s</NOBR></TH>
<TH %s nowrap=nowrap><NOBR>%s農場規模%s</NOBR></TH>
<TH %s nowrap=nowrap><NOBR>%s工場規模%s</NOBR></TH>
<TH %s nowrap=nowrap><NOBR>%s採掘場規模%s</NOBR></TH>
</TR>
<TR>
<TD %s align=middle nowrap=nowrap><NOBR>%s%d%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>
%s
<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
</TR>
</TABLE></CENTER>
`,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		mStr1,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgNumberCell, hconst.TagNumberBegin, rank, hconst.TagNumberEnd,
		hconst.BgInfoCell, island.Pop, hconst.UnitPop,
		mStr2,
		hconst.BgInfoCell, island.Food, hconst.UnitFood,
		hconst.BgInfoCell, island.Area, hconst.UnitArea,
		hconst.BgInfoCell, farmStr,
		hconst.BgInfoCell, factoryStr,
		hconst.BgInfoCell, mountainStr,
	))
}

// islandMap displays the island map
// mode: 0=tourist view, 1=owner view
// Ref: perl/lib/Hako/Map.pm:402
func islandMap(mode int) {
	island := variable.Islands[variable.CurrentNumber]

	out("<CENTER><TABLE BORDER><TR><TD>\n")

	// Get terrain data
	land := island.Land
	landValue := island.LandValue

	// Get commands for owner mode
	var comStr [][]string
	if variable.MainMode == "owner" {
		comStr = make([][]string, hconst.IslandSize)
		for i := 0; i < hconst.IslandSize; i++ {
			comStr[i] = make([]string, hconst.IslandSize)
		}

		for i := 0; i < len(island.Commands) && i < hconst.CommandMax; i++ {
			j := i + 1
			com := island.Commands[i]
			if com.Kind < 20 {
				comStr[com.X][com.Y] += fmt.Sprintf(" [%d]%s", j, hconst.ComName[com.Kind])
			}
		}
	}

	// Output coordinate bar (top)
	out("<IMG SRC=\"xbar.gif\" width=400 height=16><BR>")

	// Output terrain for each cell
	for y := 0; y < hconst.IslandSize; y++ {
		for x := 0; x < hconst.IslandSize; x++ {
			l := land[x][y]
			lv := landValue[x][y]
			var comStrXY string
			if comStr != nil {
				comStrXY = comStr[x][y]
			}
			landString(l, lv, x, y, mode, comStrXY)
		}
		out("<BR>\n")
	}

	// Output coordinate bar (bottom)
	out("<IMG SRC=\"xbar.gif\" width=400 height=16><BR>")
	out("</TD></TR></TABLE></CENTER>\n")
}

// landString outputs HTML for a single terrain cell
// Ref: perl/lib/Hako/Map.pm:460
func landString(l, lv, x, y, mode int, comStr string) {
	point := fmt.Sprintf("(%d,%d)", x, y)
	var image, alt string

	switch l {
	case hconst.LandSea:
		if lv == 1 {
			image = "land14.gif"
			alt = "海(浅瀬)"
		} else {
			image = "land0.gif"
			alt = "海"
		}

	case hconst.LandWaste:
		if lv == 1 {
			image = "land13.gif"
			alt = "荒地"
		} else {
			image = "land1.gif"
			alt = "荒地"
		}

	case hconst.LandPlains:
		image = "land2.gif"
		alt = "平地"

	case hconst.LandForest:
		if mode == 1 {
			image = "land6.gif"
			alt = fmt.Sprintf("森(%d%s)", lv, hconst.UnitTree)
		} else {
			image = "land6.gif"
			alt = "森"
		}

	case hconst.LandTown:
		var p int
		var n string
		if lv < 30 {
			p = 3
			n = "村"
		} else if lv < 100 {
			p = 4
			n = "町"
		} else {
			p = 5
			n = "都市"
		}
		image = fmt.Sprintf("land%d.gif", p)
		alt = fmt.Sprintf("%s(%d%s)", n, lv, hconst.UnitPop)

	case hconst.LandFarm:
		image = "land7.gif"
		alt = fmt.Sprintf("農場(%d0%s規模)", lv, hconst.UnitPop)

	case hconst.LandFactory:
		image = "land8.gif"
		alt = fmt.Sprintf("工場(%d0%s規模)", lv, hconst.UnitPop)

	case hconst.LandBase:
		if mode == 0 {
			image = "land6.gif"
			alt = "森"
		} else {
			level := core.ExpToLevel(l, lv)
			image = "land9.gif"
			alt = fmt.Sprintf("ミサイル基地 (レベル %d/経験値 %d)", level, lv)
		}

	case hconst.LandSbase:
		if mode == 0 {
			image = "land0.gif"
			alt = "海"
		} else {
			level := core.ExpToLevel(l, lv)
			image = "land12.gif"
			alt = fmt.Sprintf("海底基地 (レベル %d/経験値 %d)", level, lv)
		}

	case hconst.LandDefence:
		image = "land10.gif"
		alt = "防衛施設"

	case hconst.LandHaribote:
		image = "land10.gif"
		if mode == 0 {
			alt = "防衛施設"
		} else {
			alt = "ハリボテ"
		}

	case hconst.LandOil:
		image = "land16.gif"
		alt = "海底油田"

	case hconst.LandMountain:
		if lv > 0 {
			image = "land15.gif"
			alt = fmt.Sprintf("山(採掘場%d0%s規模)", lv, hconst.UnitPop)
		} else {
			image = "land11.gif"
			alt = "山"
		}

	case hconst.LandMonument:
		image = hconst.MonumentImage[lv]
		alt = hconst.MonumentName[lv]

	case hconst.LandMonster:
		kind, name, hp := core.MonsterSpec(lv)
		special := hconst.MonsterSpecial[kind]
		image = hconst.MonsterImage[kind]

		// Hardened?
		if (special == 3 && (variable.IslandTurn%2) == 1) ||
			(special == 4 && (variable.IslandTurn%2) == 0) {
			image = hconst.MonsterImage2[kind]
		}
		alt = fmt.Sprintf("怪獣%s(体力%d)", name, hp)

	default:
		image = "land1.gif"
		alt = "不明"
	}

	// For development view, add click handler for coordinates
	if mode == 1 {
		out(fmt.Sprintf("<A HREF=\"JavaScript:void(0);\" onclick=\"ps(%d,%d)\">", x, y))
	}

	out(fmt.Sprintf("<IMG SRC=\"%s\" ALT=\"%s %s %s\" width=32 height=32 BORDER=0>", image, point, alt, comStr))

	if mode == 1 {
		out("</A>")
	}
}

// Template functions

func tempPrintIslandHead() {
	out(fmt.Sprintf(`<CENTER>
%s%s「%s島」%sへようこそ！！%s<BR>
%s<BR>
</CENTER>
`,
		hconst.TagBigBegin, hconst.TagNameBegin, variable.CurrentName, hconst.TagNameEnd, hconst.TagBigEnd,
		hconst.TempBack,
	))
}

func tempOwner() {
	out(fmt.Sprintf(`<CENTER>
%s%s%s島%s開発計画%s<BR>
%s<BR>
</CENTER>
<SCRIPT Language="JavaScript">
<!--
function ps(x, y) {
    document.forms[0].elements[4].options[x].selected = true;
    document.forms[0].elements[5].options[y].selected = true;
    return true;
}

function ns(x) {
    document.forms[0].elements[2].options[x].selected = true;
    return true;
}

//-->
</SCRIPT>
`,
		hconst.TagBigBegin, hconst.TagNameBegin, variable.CurrentName, hconst.TagNameEnd, hconst.TagBigEnd,
		hconst.TempBack,
	))

	islandInfo()

	// Phase 1: Simplified command form
	out(fmt.Sprintf(`<CENTER>
<TABLE BORDER>
<TR>
<TD %s>
<CENTER>
<FORM action="%s" method=POST>
<INPUT TYPE=submit VALUE="計画送信" NAME=CommandButton%s>
<HR>
<B>パスワード</B></BR>
<INPUT TYPE=password NAME=PASSWORD VALUE="%s">
<HR>
<B>計画フォーム(Phase 1: 簡略版)</B>
</FORM>
</CENTER>
</TD>
</TR>
</TABLE>
</CENTER>
`,
		hconst.BgInputCell,
		hconst.ThisFile,
		variable.Islands[variable.CurrentNumber].ID,
		variable.DefaultPassword,
	))

	islandMap(1)
}

func tempRecent(mode int) {
	out(fmt.Sprintf("<H1>%s近況%s</H1>\n", hconst.TagHeaderBegin, hconst.TagHeaderEnd))
	logPrintLocal(mode)
}

func logPrintLocal(mode int) {
	// Phase 1: Stub - would call logFilePrint for each log file
	out("<P>(近況ログ表示は未実装)</P>\n")
}

func tempLbbsHead() {
	out(fmt.Sprintf("<H1>%s%s島ローカル掲示板%s</H1>\n",
		hconst.TagHeaderBegin, variable.CurrentName, hconst.TagHeaderEnd))
}

func tempLbbsInput() {
	out("<P>ローカル掲示板書き込み(観光者用)(Phase 1: 簡略版)</P>\n")
}

func tempLbbsInputOW() {
	out("<P>ローカル掲示板書き込み(所有者用)(Phase 1: 簡略版)</P>\n")
}

func tempLbbsContents() {
	out("<P>掲示板メッセージ一覧(Phase 1: 簡略版)</P>\n")
}

func tempLbbsNoMessage() {
	out("<H1>名前かメッセージが入力されていません。</H1>\n")
}

func tempLbbsDelete() {
	out("<H1>メッセージを削除しました。</H1>\n")
}

func tempLbbsAdd() {
	out("<H1>メッセージを追加しました。</H1>\n")
}

func tempCommandDelete() {
	out("<H1>計画を削除しました。</H1>\n")
}

func tempCommandAdd() {
	out("<H1>計画を追加しました。</H1>\n")
}

func tempComment() {
	out("<H1>コメントを変更しました。</H1>\n")
}
