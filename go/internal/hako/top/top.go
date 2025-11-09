// Package top provides top page display functions for the Hakoniwa game.
// This file is translated from Perl lib/Hako/Top.pm
//
// Ref: perl/lib/Hako/Top.pm
package top

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neguse/hakoniwa/internal/hako/core"
	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// out outputs a string to the output buffer
func out(s string) {
	variable.OutputBuffer.WriteString(s)
}

// TopPageMain displays the top page
// Ref: perl/lib/Hako/Top.pm:32
func TopPageMain() {
	// Release lock
	core.Unlock()

	// Output template
	tempTopPage()
}

// tempTopPage outputs the top page HTML
// Ref: perl/lib/Hako/Top.pm:42
func tempTopPage() {
	// Title
	out(fmt.Sprintf("%s%s%s\n", hconst.TagTitleBegin, hconst.Title, hconst.TagTitleEnd))

	// Debug mode turn button
	if hconst.Debug {
		out(fmt.Sprintf(`<FORM action="%s" method="POST">
<INPUT TYPE="submit" VALUE="ターンを進める" NAME="TurnButton">
</FORM>
`, hconst.ThisFile))
	}

	// Money column header (optional based on hide mode)
	mStr1 := ""
	if hconst.HideMoneyMode != 0 {
		mStr1 = fmt.Sprintf("<TH %s align=center nowrap=nowrap><NOBR>%s資金%s</NOBR></TH>",
			hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd)
	}

	// Island list for select
	islandList := getIslandList(variable.DefaultID)

	// Main form and table header
	out(fmt.Sprintf(`<H1>%sターン%d%s</H1>

<HR>
<H1>%s自分の島へ%s</H1>
<FORM action="%s" method="POST">
あなたの島の名前は？<BR>
<SELECT NAME="ISLANDID">
%s</SELECT><BR>

パスワードをどうぞ！！<BR>
<INPUT TYPE="password" NAME="PASSWORD" VALUE="%s" SIZE=32 MAXLENGTH=32><BR>
<INPUT TYPE="submit" VALUE="開発しに行く" NAME="OwnerButton"><BR>
</FORM>

<HR>

<H1>%s諸島の状況%s</H1>
<P>
島の名前をクリックすると、<B>観光</B>することができます。
</P>
<TABLE BORDER>
<TR>
<TH %s align=center nowrap=nowrap><NOBR>%s順位%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s島%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s人口%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s面積%s</NOBR></TH>
%s
<TH %s align=center nowrap=nowrap><NOBR>%s食料%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s農場規模%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s工場規模%s</NOBR></TH>
<TH %s align=center nowrap=nowrap><NOBR>%s採掘場規模%s</NOBR></TH>
</TR>
`,
		hconst.TagHeaderBegin, variable.IslandTurn, hconst.TagHeaderEnd,
		hconst.TagHeaderBegin, hconst.TagHeaderEnd,
		hconst.ThisFile,
		islandList,
		variable.DefaultPassword,
		hconst.TagHeaderBegin, hconst.TagHeaderEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		mStr1,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
		hconst.BgTitleCell, hconst.TagTHBegin, hconst.TagTHEnd,
	))

	// Island rows
	for ii := 0; ii < variable.IslandNumber; ii++ {
		j := ii + 1
		island := variable.Islands[ii]

		id := island.ID
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

		// Island name with absent counter
		name := ""
		if island.Absent == 0 {
			name = fmt.Sprintf("%s%s島%s", hconst.TagNameBegin, island.Name, hconst.TagNameEnd)
		} else {
			name = fmt.Sprintf("%s%s島(%d)%s", hconst.TagName2Begin, island.Name, island.Absent, hconst.TagName2End)
		}

		// Parse prize field: "flags,monsters,turns"
		prize := parsePrize(island.Prize)

		// Money display
		mStr1 := ""
		if hconst.HideMoneyMode == 1 {
			mStr1 = fmt.Sprintf("<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>",
				hconst.BgInfoCell, island.Money, hconst.UnitMoney)
		} else if hconst.HideMoneyMode == 2 {
			mTmp := core.AboutMoney(island.Money)
			mStr1 = fmt.Sprintf("<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>",
				hconst.BgInfoCell, mTmp)
		}

		// Output island row
		out(fmt.Sprintf(`<TR>
<TD %s ROWSPAN=2 align=center nowrap=nowrap><NOBR>%s%d%s</NOBR></TD>
<TD %s ROWSPAN=2 align=left nowrap=nowrap>
<NOBR>
<A STYlE="text-decoration:none" HREF="%s?Sight=%s">
%s
</A>
</NOBR><BR>
%s
</TD>
<TD %s align=right nowrap=nowrap>
<NOBR>%d%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>
%s
<TD %s align=right nowrap=nowrap><NOBR>%d%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
<TD %s align=right nowrap=nowrap><NOBR>%s</NOBR></TD>
</TR>
<TR>
<TD %s COLSPAN=7 align=left nowrap=nowrap><NOBR>%sコメント：%s%s</NOBR></TD>
</TR>
`,
			hconst.BgNumberCell, hconst.TagNumberBegin, j, hconst.TagNumberEnd,
			hconst.BgNameCell,
			hconst.ThisFile, id,
			name,
			prize,
			hconst.BgInfoCell, island.Pop, hconst.UnitPop,
			hconst.BgInfoCell, island.Area, hconst.UnitArea,
			mStr1,
			hconst.BgInfoCell, island.Food, hconst.UnitFood,
			hconst.BgInfoCell, farmStr,
			hconst.BgInfoCell, factoryStr,
			hconst.BgInfoCell, mountainStr,
			hconst.BgCommentCell, hconst.TagTHBegin, hconst.TagTHEnd, island.Comment,
		))
	}

	out("</TABLE>\n\n<HR>\n")

	// New island section
	out(fmt.Sprintf("<H1>%s新しい島を探す%s</H1>\n", hconst.TagHeaderBegin, hconst.TagHeaderEnd))

	if variable.IslandNumber < hconst.MaxIsland {
		out(fmt.Sprintf(`<FORM action="%s" method="POST">
どんな名前をつける予定？<BR>
<INPUT TYPE="text" NAME="ISLANDNAME" SIZE=32 MAXLENGTH=32>島<BR>
パスワードは？<BR>
<INPUT TYPE="password" NAME="PASSWORD" SIZE=32 MAXLENGTH=32><BR>
念のためパスワードをもう一回<BR>
<INPUT TYPE="password" NAME="PASSWORD2" SIZE=32 MAXLENGTH=32><BR>

<INPUT TYPE="submit" VALUE="探しに行く" NAME="NewIslandButton">
</FORM>
`, hconst.ThisFile))
	} else {
		out("        島の数が最大数です・・・現在登録できません。\n")
	}

	// Change name/password section
	out(fmt.Sprintf(`<HR>
<H1>%s島の名前とパスワードの変更%s</H1>
<P>
(注意)名前の変更には%d%sかかります。
</P>
<FORM action="%s" method="POST">
どの島ですか？<BR>
<SELECT NAME="ISLANDID">
%s
</SELECT>
<BR>
どんな名前に変えますか？(変更する場合のみ)<BR>
<INPUT TYPE="text" NAME="ISLANDNAME" SIZE=32 MAXLENGTH=32>島<BR>
パスワードは？(必須)<BR>
<INPUT TYPE="password" NAME="OLDPASS" SIZE=32 MAXLENGTH=32><BR>
新しいパスワードは？(変更する時のみ)<BR>
<INPUT TYPE="password" NAME="PASSWORD" SIZE=32 MAXLENGTH=32><BR>
念のためパスワードをもう一回(変更する時のみ)<BR>
<INPUT TYPE="password" NAME="PASSWORD2" SIZE=32 MAXLENGTH=32><BR>

<INPUT TYPE="submit" VALUE="変更する" NAME="ChangeInfoButton">
</FORM>

<HR>

<H1>%s最近の出来事%s</H1>
`,
		hconst.TagHeaderBegin, hconst.TagHeaderEnd,
		hconst.CostChangeName, hconst.UnitMoney,
		hconst.ThisFile,
		islandList,
		hconst.TagHeaderBegin, hconst.TagHeaderEnd,
	))

	// Recent events log
	logPrintTop()

	// Discovery history
	out(fmt.Sprintf("<H1>%s発見の記録%s</H1>\n", hconst.TagHeaderBegin, hconst.TagHeaderEnd))
	historyPrint()
}

// parsePrize parses the prize field and returns HTML for prize display
// Prize format: "flags,monsters,turns"
// Ref: perl/lib/Hako/Top.pm:124-163
func parsePrize(prizeField int) string {
	// Convert to string for parsing
	prizeStr := fmt.Sprintf("%d", prizeField)

	// Parse: flags,monsters,turns
	parts := strings.Split(prizeStr, ",")
	flags := 0
	monsters := 0
	turns := ""

	if len(parts) >= 1 {
		flags, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 {
		monsters, _ = strconv.Atoi(parts[1])
	}
	if len(parts) >= 3 {
		turns = parts[2]
	}

	prize := ""

	// Turn prizes
	turnParts := strings.Split(turns, ",")
	for _, tp := range turnParts {
		if tp != "" {
			prize += fmt.Sprintf(`<IMG SRC="prize0.gif" ALT="%s%s" WIDTH=16 HEIGHT=16> `,
				tp, hconst.Prize[0])
		}
	}

	// Flag prizes (prizes 1-9)
	f := 1
	for i := 1; i < 10; i++ {
		if (flags & f) != 0 {
			prize += fmt.Sprintf(`<IMG SRC="prize%d.gif" ALT="%s" WIDTH=16 HEIGHT=16> `,
				i, hconst.Prize[i])
		}
		f *= 2
	}

	// Monster prizes
	f = 1
	max := -1
	mNameList := ""
	for i := 0; i < hconst.MonsterNumber; i++ {
		if (monsters & f) != 0 {
			mNameList += fmt.Sprintf("[%s] ", hconst.MonsterName[i])
			max = i
		}
		f *= 2
	}
	if max != -1 {
		prize += fmt.Sprintf(`<IMG SRC="%s" ALT="%s" WIDTH=16 HEIGHT=16> `,
			hconst.MonsterImage[max], mNameList)
	}

	return prize
}

// getIslandList generates island select options
// Ref: perl/lib/Hako/Main.pm:1079
func getIslandList(selectID string) string {
	list := ""
	for i := 0; i < variable.IslandNumber; i++ {
		name := variable.Islands[i].Name
		id := variable.Islands[i].ID
		selected := ""
		if id == selectID {
			selected = "SELECTED"
		}
		list += fmt.Sprintf("<OPTION VALUE=\"%s\" %s>%s島\n", id, selected, name)
	}
	return list
}

// logPrintTop prints top N turns of logs
// Ref: perl/lib/Hako/Top.pm:265
func logPrintTop() {
	for i := 0; i < hconst.TopLogTurn; i++ {
		logFilePrint(i, "0", 0)
	}
}

// logFilePrint prints a specific log file
// Ref: perl/lib/Hako/Main.pm:1030
func logFilePrint(fileNumber int, id string, mode int) {
	filename := fmt.Sprintf("%s/hakojima.log%d", hconst.DirName, fileNumber)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse: m,turn,id1,id2,message
		parts := strings.SplitN(line, ",", 5)
		if len(parts) < 5 {
			continue
		}

		m, _ := strconv.Atoi(parts[0])
		turn := parts[1]
		id1 := parts[2]
		id2 := parts[3]
		message := parts[4]

		// Secret handling
		mStr := ""
		if m == 1 {
			if mode == 0 || id1 != id {
				// No permission to see secret
				continue
			}
			mStr = "<B>(機密)</B>"
		}

		// Filter by island
		if id != "0" {
			if id != id1 && id != id2 {
				continue
			}
		}

		// Output
		out(fmt.Sprintf("<NOBR>%sターン%s%s%s：%s</NOBR><BR>\n",
			hconst.TagNumberBegin, turn, mStr, hconst.TagNumberEnd, message))
	}
}

// historyPrint prints discovery history
// Ref: perl/lib/Hako/Top.pm:273
func historyPrint() {
	filename := fmt.Sprintf("%s/hakojima.his", hconst.DirName)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	// Read all lines
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Reverse order
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]

		// Parse: turn,message
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			continue
		}

		turn := parts[0]
		message := parts[1]

		out(fmt.Sprintf("<NOBR>%sターン%s%s：%s</NOBR><BR>\n",
			hconst.TagNumberBegin, turn, hconst.TagNumberEnd, message))
	}
}
