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

	// Branch by mode
	commands := island.Commands

	if variable.CommandMode == "delete" {
		// Delete command
		slideFront(commands, variable.CommandPlanNumber)
		tempCommandDelete()
	} else if variable.CommandKind == hconst.ComAutoPrepare || variable.CommandKind == hconst.ComAutoPrepare2 {
		// フル整地、フル地ならし
		// 座標配列を作る
		core.MakeRandomPointArray()
		land := island.Land

		// コマンドの種類決定
		kind := hconst.ComPrepare
		if variable.CommandKind == hconst.ComAutoPrepare2 {
			kind = hconst.ComPrepare2
		}

		i := 0
		j := 0
		for j < hconst.IslandSize*hconst.IslandSize && i < hconst.CommandMax {
			x := variable.Rpx[j]
			y := variable.Rpy[j]
			if land[x][y] == hconst.LandWaste {
				slideBack(commands, variable.CommandPlanNumber)
				commands[variable.CommandPlanNumber] = types.Command{
					Kind:   kind,
					Target: "",
					X:      x,
					Y:      y,
					Arg:    0,
				}
				i++
			}
			j++
		}
		tempCommandAdd()
	} else if variable.CommandKind == hconst.ComAutoDelete {
		// 全消し
		for i := 0; i < hconst.CommandMax; i++ {
			slideFront(commands, variable.CommandPlanNumber)
		}
		tempCommandDelete()
	} else {
		// Normal command
		if variable.CommandMode == "insert" {
			slideBack(commands, variable.CommandPlanNumber)
		}
		tempCommandAdd()

		// Register command
		if variable.CommandPlanNumber >= 0 && variable.CommandPlanNumber < len(commands) {
			commands[variable.CommandPlanNumber] = types.Command{
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
		// 偶数行目なら番号を出力
		if (y % 2) == 0 {
			out(fmt.Sprintf("<IMG SRC=\"space%d.gif\" width=16 height=32>", y))
		}

		for x := 0; x < hconst.IslandSize; x++ {
			l := land[x][y]
			lv := landValue[x][y]
			var comStrXY string
			if comStr != nil {
				comStrXY = comStr[x][y]
			}
			landString(l, lv, x, y, mode, comStrXY)
		}

		// 奇数行目なら番号を出力
		if (y % 2) == 1 {
			out(fmt.Sprintf("<IMG SRC=\"space%d.gif\" width=16 height=32>", y))
		}

		out("<BR>\n")
	}

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
	island := variable.Islands[variable.CurrentNumber]

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

	out(fmt.Sprintf(`<CENTER>
<TABLE BORDER>
<TR>
<TD %s >
<CENTER>
<FORM action="%s" method=POST>
<INPUT TYPE=submit VALUE="計画送信" NAME=CommandButton%s>
<HR>
<B>パスワード</B></BR>
<INPUT TYPE=password NAME=PASSWORD VALUE="%s">
<HR>
<B>計画番号</B><SELECT NAME=NUMBER>
`,
		hconst.BgInputCell,
		hconst.ThisFile,
		island.ID,
		variable.DefaultPassword,
	))

	// 計画番号
	for i := 0; i < hconst.CommandMax; i++ {
		j := i + 1
		out(fmt.Sprintf("<OPTION VALUE=%d>%d\n", i, j))
	}

	out(`</SELECT><BR>
<HR>
<B>開発計画</B><BR>
<SELECT NAME=COMMAND>
`)

	// コマンド
	for i := 0; i < len(hconst.ComList); i++ {
		kind := hconst.ComList[i]
		cost := hconst.ComCost[kind]
		costStr := ""
		if cost == 0 {
			costStr = "無料"
		} else if cost < 0 {
			costStr = fmt.Sprintf("%d%s", -cost, hconst.UnitFood)
		} else {
			costStr = fmt.Sprintf("%d%s", cost, hconst.UnitMoney)
		}

		selected := ""
		if variable.DefaultKind != 0 && kind == variable.DefaultKind {
			selected = "SELECTED"
		}
		out(fmt.Sprintf("<OPTION VALUE=%d %s>%s(%s)\n", kind, selected, hconst.ComName[kind], costStr))
	}

	out(`</SELECT>
<HR>
<B>座標(</B>
<SELECT NAME=POINTX>

`)

	for i := 0; i < hconst.IslandSize; i++ {
		if variable.DefaultX != 0 && i == variable.DefaultX {
			out(fmt.Sprintf("<OPTION VALUE=%d SELECTED>%d\n", i, i))
		} else {
			out(fmt.Sprintf("<OPTION VALUE=%d>%d\n", i, i))
		}
	}

	out(`</SELECT>, <SELECT NAME=POINTY>
`)

	for i := 0; i < hconst.IslandSize; i++ {
		if variable.DefaultY != 0 && i == variable.DefaultY {
			out(fmt.Sprintf("<OPTION VALUE=%d SELECTED>%d\n", i, i))
		} else {
			out(fmt.Sprintf("<OPTION VALUE=%d>%d\n", i, i))
		}
	}

	out(`</SELECT><B>)</B>
<HR>
<B>数量</B><SELECT NAME=AMOUNT>
`)

	// 数量
	for i := 0; i < 100; i++ {
		out(fmt.Sprintf("<OPTION VALUE=%d>%d\n", i, i))
	}

	out(fmt.Sprintf(`</SELECT>
<HR>
<B>目標の島</B><BR>
<SELECT NAME=TARGETID>
%s<BR>
</SELECT>
<HR>
<B>動作</B><BR>
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=insert CHECKED>挿入
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=write>上書き<BR>
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=delete>削除
<HR>
<INPUT TYPE=submit VALUE="計画送信" NAME=CommandButton%s>

</CENTER>
</FORM>
</TD>
<TD %s>
`,
		variable.TargetList,
		island.ID,
		hconst.BgMapCell,
	))

	islandMap(1) // 島の地図、所有者モード

	out(fmt.Sprintf(`</TD>
<TD %s>
`, hconst.BgCommandCell))

	// 入力済みコマンド表示
	for i := 0; i < hconst.CommandMax; i++ {
		tempCommand(i, island.Commands[i])
	}

	out(`
</TD>
</TR>
</TABLE>
</CENTER>
<HR>
<CENTER>
`)

	out(fmt.Sprintf(`%sコメント更新%s<BR>
<FORM action="%s" method="POST">
コメント<INPUT TYPE=text NAME=MESSAGE SIZE=80><BR>
パスワード<INPUT TYPE=password NAME=PASSWORD VALUE="%s">
<INPUT TYPE=submit VALUE="コメント更新" NAME=MessageButton%s>
</FORM>
</CENTER>
`,
		hconst.TagBigBegin, hconst.TagBigEnd,
		hconst.ThisFile,
		variable.DefaultPassword,
		island.ID,
	))
}

// tempCommand displays a single command in the command list
// Ref: perl/lib/Hako/Map.pm:837
func tempCommand(number int, command types.Command) {
	kind := command.Kind
	target := command.Target
	x := command.X
	y := command.Y
	arg := command.Arg

	name := fmt.Sprintf("%s%s%s", hconst.TagComNameBegin, hconst.ComName[kind], hconst.TagComNameEnd)
	point := fmt.Sprintf("%s(%d,%d)%s", hconst.TagNameBegin, x, y, hconst.TagNameEnd)

	targetName := variable.IDToName[target]
	if targetName != "" {
		targetName = fmt.Sprintf("%s%s島%s", hconst.TagNameBegin, targetName, hconst.TagNameEnd)
	} else if target != "" {
		targetName = fmt.Sprintf("%s無人%s", hconst.TagNameBegin, hconst.TagNameEnd)
	}

	value := arg * hconst.ComCost[kind]
	if value == 0 {
		value = hconst.ComCost[kind]
	}

	valueStr := ""
	if value < 0 {
		valueStr = fmt.Sprintf("%d%s", -value, hconst.UnitFood)
	} else {
		valueStr = fmt.Sprintf("%d%s", value, hconst.UnitMoney)
	}
	valueStr = fmt.Sprintf("%s%s%s", hconst.TagNameBegin, valueStr, hconst.TagNameEnd)

	j := fmt.Sprintf("%02d：", number+1)

	out(fmt.Sprintf(`<A STYlE="text-decoration:none" HREF="JavaScript:void(0);" onClick="ns(%d)"><NOBR>%s%s%s<FONT COLOR="%s">`,
		number, hconst.TagNumberBegin, j, hconst.TagNumberEnd, hconst.NormalColor))

	if kind == hconst.ComDoNothing || kind == hconst.ComGiveup {
		out(name)
	} else if kind == hconst.ComMissileNM || kind == hconst.ComMissilePP ||
		kind == hconst.ComMissileST || kind == hconst.ComMissileLD {
		// ミサイル系
		n := "無制限"
		if arg != 0 {
			n = fmt.Sprintf("%d発", arg)
		}
		out(fmt.Sprintf("%s%sへ%s(%s%s%s)", targetName, point, name, hconst.TagNameBegin, n, hconst.TagNameEnd))
	} else if kind == hconst.ComSendMonster {
		// 怪獣派遣
		out(fmt.Sprintf("%sへ%s", targetName, name))
	} else if kind == hconst.ComSell {
		// 食料輸出
		out(fmt.Sprintf("%s%s", name, valueStr))
	} else if kind == hconst.ComPropaganda {
		// 誘致活動
		out(name)
	} else if kind == hconst.ComMoney || kind == hconst.ComFood {
		// 援助
		out(fmt.Sprintf("%sへ%s%s", targetName, name, valueStr))
	} else if kind == hconst.ComDestroy {
		// 掘削
		if arg != 0 {
			out(fmt.Sprintf("%sで%s(予算%s)", point, name, valueStr))
		} else {
			out(fmt.Sprintf("%sで%s", point, name))
		}
	} else if kind == hconst.ComFarm || kind == hconst.ComFactory || kind == hconst.ComMountain {
		// 回数付き
		if arg == 0 {
			out(fmt.Sprintf("%sで%s", point, name))
		} else {
			out(fmt.Sprintf("%sで%s(%d回)", point, name, arg))
		}
	} else {
		// 座標付き
		out(fmt.Sprintf("%sで%s", point, name))
	}

	out("</FONT></NOBR></A><BR>")
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
	out(fmt.Sprintf(`<HR>
<CENTER>
%s%s%s島%s観光者通信%s<BR>
</CENTER>
`,
		hconst.TagBigBegin, hconst.TagNameBegin, variable.CurrentName, hconst.TagNameEnd, hconst.TagBigEnd))
}

func tempLbbsInput() {
	out(fmt.Sprintf(`<CENTER>
<FORM action="%s" method="POST">
<TABLE BORDER>
<TR>
<TH>名前</TH>
<TH>内容</TH>
<TH>動作</TH>
</TR>
<TR>
<TD><INPUT TYPE="text" SIZE=32 MAXLENGTH=32 NAME="LBBSNAME" VALUE="%s"></TD>
<TD><INPUT TYPE="text" SIZE=80 NAME="LBBSMESSAGE"></TD>
<TD><INPUT TYPE="submit" VALUE="記帳する" NAME="LbbsButtonSS%s"></TD>
</TR>
</TABLE>
</FORM>
</CENTER>
`,
		hconst.ThisFile,
		variable.DefaultName,
		variable.CurrentID,
	))
}

func tempLbbsInputOW() {
	out(fmt.Sprintf(`<CENTER>
<FORM action="%s" method="POST">
<TABLE BORDER>
<TR>
<TH>名前</TH>
<TH COLSPAN=2>内容</TH>
</TR>
<TR>
<TD><INPUT TYPE="text" SIZE=32 MAXLENGTH=32 NAME="LBBSNAME" VALUE="%s"></TD>
<TD COLSPAN=2><INPUT TYPE="text" SIZE=80 NAME="LBBSMESSAGE"></TD>
</TR>
<TR>
<TH>パスワード</TH>
<TH COLSPAN=2>動作</TH>
</TR>
<TR>
<TD><INPUT TYPE=password SIZE=32 MAXLENGTH=32 NAME=PASSWORD VALUE="%s"></TD>
<TD align=right>
<INPUT TYPE="submit" VALUE="記帳する" NAME="LbbsButtonOW%s">
</TD>
<TD align=right>
番号
<SELECT NAME=NUMBER>
`,
		hconst.ThisFile,
		variable.DefaultName,
		variable.DefaultPassword,
		variable.CurrentID,
	))

	// 発言番号
	for i := 0; i < hconst.LbbsMax; i++ {
		j := i + 1
		out(fmt.Sprintf("<OPTION VALUE=%d>%d\n", i, j))
	}

	out(fmt.Sprintf(`</SELECT>
<INPUT TYPE="submit" VALUE="削除する" NAME="LbbsButtonDL%s">
</TD>
</TR>
</TABLE>
</FORM>
</CENTER>
`, variable.CurrentID))
}

func tempLbbsContents() {
	lbbs := variable.Islands[variable.CurrentNumber].Lbbs

	out(`<CENTER>
<TABLE BORDER>
<TR>
<TH>番号</TH>
<TH>記帳内容</TH>
</TR>
`)

	for i := 0; i < hconst.LbbsMax && i < len(lbbs); i++ {
		line := lbbs[i].Message
		// Parse format: "mode>name>message"
		// mode: 0=tourist, 1=owner
		var mode, name, message string
		if len(line) > 0 {
			// Simple parsing
			parts := splitN(line, ">", 3)
			if len(parts) >= 3 {
				mode = parts[0]
				name = parts[1]
				message = parts[2]

				j := i + 1
				out(fmt.Sprintf("<TR><TD align=center>%s%d%s</TD>", hconst.TagNumberBegin, j, hconst.TagNumberEnd))

				if mode == "0" {
					// 観光者
					out(fmt.Sprintf("<TD>%s%s > %s%s</TD></TR>", hconst.TagLbbsSSBegin, name, message, hconst.TagLbbsSSEnd))
				} else {
					// 島主
					out(fmt.Sprintf("<TD>%s%s > %s%s</TD></TR>", hconst.TagLbbsOWBegin, name, message, hconst.TagLbbsOWEnd))
				}
			}
		}
	}

	out(`</TD></TR></TABLE></CENTER>
`)
}

// splitN splits a string by separator, limiting to n parts
func splitN(s, sep string, n int) []string {
	result := []string{}
	for i := 0; i < n-1; i++ {
		idx := -1
		for j := 0; j < len(s); j++ {
			if s[j:j+len(sep)] == sep {
				idx = j
				break
			}
		}
		if idx == -1 {
			result = append(result, s)
			return result
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	result = append(result, s)
	return result
}

func tempLbbsNoMessage() {
	out(fmt.Sprintf("%s名前または内容の欄が空欄です。%s%s",
		hconst.TagBigBegin, hconst.TagBigEnd, hconst.TempBack))
}

func tempLbbsDelete() {
	out(fmt.Sprintf("%s記帳内容を削除しました%s<HR>\n",
		hconst.TagBigBegin, hconst.TagBigEnd))
}

func tempLbbsAdd() {
	out(fmt.Sprintf("%s記帳を行いました%s<HR>\n",
		hconst.TagBigBegin, hconst.TagBigEnd))
}

func tempCommandDelete() {
	out(fmt.Sprintf("%sコマンドを削除しました%s<HR>\n",
		hconst.TagBigBegin, hconst.TagBigEnd))
}

func tempCommandAdd() {
	out(fmt.Sprintf("%sコマンドを登録しました%s<HR>\n",
		hconst.TagBigBegin, hconst.TagBigEnd))
}

func tempComment() {
	out(fmt.Sprintf("%sコメントを更新しました%s<HR>\n",
		hconst.TagBigBegin, hconst.TagBigEnd))
}
