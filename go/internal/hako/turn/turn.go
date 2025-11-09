// Package turn provides turn processing functions for the Hakoniwa game.
// This file is translated from Perl lib/Hako/Turn.pm (3,640 lines)
//
// Ref: perl/lib/Hako/Turn.pm
package turn

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"

	"github.com/neguse/hakoniwa/internal/hako/core"
	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/top"
	"github.com/neguse/hakoniwa/internal/hako/types"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// Surrounding 2-hex coordinates
// Ref: perl/lib/Hako/Turn.pm:33-34
var (
	ax = []int{0, 1, 1, 1, 0, -1, 0, 1, 2, 2, 2, 1, 0, -1, -1, -2, -1, -1, 0}
	ay = []int{0, -1, 0, 1, 1, 0, -1, -2, -1, 0, 1, 2, 2, 2, 1, 0, -1, -2, -2}
)

// out outputs a string to the output buffer
func out(s string) {
	variable.OutputBuffer.WriteString(s)
}

//======================================================================
// New Island Creation Mode
//======================================================================

// NewIslandMain handles new island creation
// Ref: perl/lib/Hako/Turn.pm:40
func NewIslandMain() {
	// Check if island limit reached
	if variable.IslandNumber >= hconst.MaxIsland {
		core.Unlock()
		tempNewIslandFull()
		return
	}

	// Check if name exists
	if variable.CurrentName == "" {
		core.Unlock()
		tempNewIslandNoName()
		return
	}

	// Check if name is valid
	invalidPattern := regexp.MustCompile(`[,\?\(\)\<\>\$]|^無人$`)
	if invalidPattern.MatchString(variable.CurrentName) {
		core.Unlock()
		tempNewIslandBadName()
		return
	}

	// Check for duplicate name
	if nameToNumber(variable.CurrentName) != -1 {
		core.Unlock()
		tempNewIslandAlready()
		return
	}

	// Check password exists
	if variable.InputPassword == "" {
		core.Unlock()
		tempNewIslandNoPassword()
		return
	}

	// Check password confirmation
	if variable.InputPassword2 != variable.InputPassword {
		core.Unlock()
		core.TempWrongPassword()
		return
	}

	// Determine new island number
	variable.CurrentNumber = variable.IslandNumber
	variable.IslandNumber++
	newIsland := makeNewIsland()
	variable.Islands = append(variable.Islands, newIsland)

	// Set values
	newIsland.Name = variable.CurrentName
	newIsland.Score = 0
	newIsland.ID = fmt.Sprintf("%d", variable.IslandNextID)
	variable.IslandNextID++
	newIsland.Absent = hconst.GiveupTurn - 3
	newIsland.Comment = "(未登録)"
	newIsland.Password = encode(variable.InputPassword)

	// Calculate population, etc.
	estimate(variable.CurrentNumber)

	// Write data
	core.WriteIsland(newIsland)
	core.WriteIslandsFile()
	logDiscover(variable.CurrentName)

	// Release lock
	core.Unlock()

	// Discovery screen
	tempNewIslandHead(variable.CurrentName)
	// Phase 1: Call internal functions (will be exported later)
	// mapview.IslandInfo()
	// mapview.IslandMap(1)
}

// makeNewIsland creates a new island
// Ref: perl/lib/Hako/Turn.pm:124
func makeNewIsland() *types.Island {
	// Create terrain
	land, landValue := makeNewLand()

	// Create initial commands
	commands := make([]types.Command, hconst.CommandMax)
	for i := 0; i < hconst.CommandMax; i++ {
		commands[i] = types.Command{
			Kind:   hconst.ComDoNothing,
			Target: "",
			X:      0,
			Y:      0,
			Arg:    0,
		}
	}

	// Create initial BBS
	lbbs := make([]types.LbbsEntry, hconst.LbbsMax)
	for i := 0; i < hconst.LbbsMax; i++ {
		lbbs[i] = types.LbbsEntry{Message: "0>>"}
	}

	// Return island
	return &types.Island{
		Land:      land,
		LandValue: landValue,
		Commands:  commands,
		Lbbs:      lbbs,
		Money:     hconst.InitialMoney,
		Food:      hconst.InitialFood,
		Prize:     0,
	}
}

// makeNewLand creates terrain for a new island
// Ref: perl/lib/Hako/Turn.pm:160
func makeNewLand() ([][]int, [][]int) {
	// Initialize as sea
	land := make([][]int, hconst.IslandSize)
	landValue := make([][]int, hconst.IslandSize)
	for x := 0; x < hconst.IslandSize; x++ {
		land[x] = make([]int, hconst.IslandSize)
		landValue[x] = make([]int, hconst.IslandSize)
		for y := 0; y < hconst.IslandSize; y++ {
			land[x][y] = hconst.LandSea
			landValue[x][y] = 0
		}
	}

	// Place wasteland in center 4x4
	center := hconst.IslandSize/2 - 1
	for y := center - 1; y < center+3; y++ {
		for x := center - 1; x < center+3; x++ {
			land[x][y] = hconst.LandWaste
		}
	}

	// Grow land within 8x8 range
	for i := 0; i < 120; i++ {
		x := rand.Intn(8) + center - 3
		y := rand.Intn(8) + center - 3

		if countAround(land, x, y, hconst.LandSea, 7) != 7 {
			// If land around, make shallow
			// Shallow becomes wasteland
			// Wasteland becomes plains
			if land[x][y] == hconst.LandWaste {
				land[x][y] = hconst.LandPlains
				landValue[x][y] = 0
			} else {
				if landValue[x][y] == 1 {
					land[x][y] = hconst.LandWaste
					landValue[x][y] = 0
				} else {
					landValue[x][y] = 1
				}
			}
		}
	}

	// Create forests
	count := 0
	for count < 4 {
		x := rand.Intn(4) + center - 1
		y := rand.Intn(4) + center - 1

		if land[x][y] != hconst.LandForest {
			land[x][y] = hconst.LandForest
			landValue[x][y] = 5 // 500 trees initially
			count++
		}
	}

	// Create towns
	count = 0
	for count < 2 {
		x := rand.Intn(4) + center - 1
		y := rand.Intn(4) + center - 1

		if land[x][y] != hconst.LandTown && land[x][y] != hconst.LandForest {
			land[x][y] = hconst.LandTown
			landValue[x][y] = 5 // 500 people initially
			count++
		}
	}

	// Create mountain
	count = 0
	for count < 1 {
		x := rand.Intn(4) + center - 1
		y := rand.Intn(4) + center - 1

		if land[x][y] != hconst.LandTown && land[x][y] != hconst.LandForest && land[x][y] != hconst.LandMountain {
			land[x][y] = hconst.LandMountain
			landValue[x][y] = 0
			count++
		}
	}

	return land, landValue
}

//======================================================================
// Island Information Change Mode
//======================================================================

// ChangeMain handles island information changes
// Ref: perl/lib/Hako/Turn.pm:288
func ChangeMain() {
	// Get island from ID
	num, ok := variable.IDToNumber[variable.CurrentID]
	if !ok {
		core.Unlock()
		core.TempProblem()
		return
	}
	variable.CurrentNumber = num
	island := variable.Islands[variable.CurrentNumber]
	flag := false

	// Password check
	if variable.OldPassword == hconst.SpecialPassword {
		// Special password
		island.Money = 9999
		island.Food = 9999
	} else if !core.CheckPassword(island.Password, variable.OldPassword) {
		core.Unlock()
		core.TempWrongPassword()
		return
	}

	// Confirmation password
	if variable.InputPassword2 != variable.InputPassword {
		core.Unlock()
		core.TempWrongPassword()
		return
	}

	// Name change
	if variable.CurrentName != "" {
		// Check if name is valid
		invalidPattern := regexp.MustCompile(`[,\?\(\)\<\>]|^無人$`)
		if invalidPattern.MatchString(variable.CurrentName) {
			core.Unlock()
			tempNewIslandBadName()
			return
		}

		// Check for duplicate
		if nameToNumber(variable.CurrentName) != -1 {
			core.Unlock()
			tempNewIslandAlready()
			return
		}

		// Check money
		if island.Money < hconst.CostChangeName {
			core.Unlock()
			tempChangeNoMoney()
			return
		}

		// Charge
		if variable.OldPassword != hconst.SpecialPassword {
			island.Money -= hconst.CostChangeName
		}

		// Change name
		logChangeName(island.Name, variable.CurrentName)
		island.Name = variable.CurrentName
		flag = true
	}

	// Password change
	if variable.InputPassword != "" {
		island.Password = encode(variable.InputPassword)
		flag = true
	}

	if !flag && variable.OldPassword != hconst.SpecialPassword {
		core.Unlock()
		tempChangeNothing()
		return
	}

	// Write data
	core.WriteIslandsFile()
	core.Unlock()

	// Success
	tempChange()
}

//======================================================================
// Turn Processing Mode
//======================================================================

// TurnMain processes a game turn
// Ref: perl/lib/Hako/Turn.pm:387
func TurnMain() {
	// Update last update time
	variable.IslandLastTime += hconst.UnitTime

	// Shift log files backward
	for i := hconst.LogMax - 1; i >= 0; i-- {
		j := i + 1
		s := filepath.Join(hconst.DirName, fmt.Sprintf("hakojima.log%d", i))
		d := filepath.Join(hconst.DirName, fmt.Sprintf("hakojima.log%d", j))
		os.Remove(d)
		os.Rename(s, d)
	}

	// Create coordinate array
	makeRandomPointArray()

	// Turn number
	variable.IslandTurn++

	// Determine order
	order := randomArray(variable.IslandNumber)

	// Income and consumption phase
	for i := 0; i < variable.IslandNumber; i++ {
		estimate(order[i])
		income(variable.Islands[order[i]])

		// Remember population before turn
		variable.Islands[order[i]].OldPop = variable.Islands[order[i]].Pop
	}

	// Command processing
	for i := 0; i < variable.IslandNumber; i++ {
		// Repeat until return value is 1
		for doCommand(variable.Islands[order[i]]) == 0 {
		}
	}

	// Growth and single hex disasters
	for i := 0; i < variable.IslandNumber; i++ {
		doEachHex(variable.Islands[order[i]])
	}

	// Whole island processing
	remainNumber := variable.IslandNumber
	for i := 0; i < variable.IslandNumber; i++ {
		island := variable.Islands[order[i]]
		doIslandProcess(order[i], island)

		// Death check
		if island.Dead {
			island.Pop = 0
			remainNumber--
		} else if island.Pop == 0 {
			island.Dead = true
			remainNumber--

			// Death message
			tmpID := island.ID
			logDead(tmpID, island.Name)
			os.Remove(fmt.Sprintf("island.%s", tmpID))
		}
	}

	// Sort by population
	islandSort()

	// Turn prize processing
	if (variable.IslandTurn % hconst.TurnPrizeUnit) == 0 {
		island := variable.Islands[0]
		logPrize(island.ID, island.Name, fmt.Sprintf("%d%s", variable.IslandTurn, hconst.Prize[0]))
		// Update prize field (Phase 1: simplified)
		island.Prize++
	}

	// Cut island count
	variable.IslandNumber = remainNumber

	// Backup if backup turn
	if (variable.IslandTurn % hconst.BackupTurn) == 0 {
		tmp := hconst.BackupTimes - 1
		myrmtree(fmt.Sprintf("%s.bak%d", hconst.DirName, tmp))

		for i := hconst.BackupTimes - 1; i > 0; i-- {
			j := i - 1
			os.Rename(
				fmt.Sprintf("%s.bak%d", hconst.DirName, j),
				fmt.Sprintf("%s.bak%d", hconst.DirName, i),
			)
		}
		os.Rename(hconst.DirName, fmt.Sprintf("%s.bak0", hconst.DirName))
		os.Mkdir(hconst.DirName, os.FileMode(hconst.DirMode))

		// Restore log files
		for i := 0; i <= hconst.LogMax; i++ {
			os.Rename(
				filepath.Join(fmt.Sprintf("%s.bak0", hconst.DirName), fmt.Sprintf("hakojima.log%d", i)),
				filepath.Join(hconst.DirName, fmt.Sprintf("hakojima.log%d", i)),
			)
		}
		os.Rename(
			filepath.Join(fmt.Sprintf("%s.bak0", hconst.DirName), "hakojima.his"),
			filepath.Join(hconst.DirName, "hakojima.his"),
		)
	}

	// Write to file
	core.WriteIslandsFile()

	// Write logs
	logFlush()

	// Trim history log
	logHistoryTrim()

	// To top page
	top.TopPageMain()
}

//======================================================================
// Helper Functions
//======================================================================

func myrmtree(dirName string) {
	os.RemoveAll(dirName)
}

// income processes income for an island
// Ref: perl/lib/Hako/Turn.pm:518-540
func income(island *types.Island) {
	pop := island.Pop
	farm := island.Farm * 10      // Farm capacity
	factory := island.Factory     // Factory capacity
	mountain := island.Mountain   // Mining site capacity

	// Income phase
	if pop > farm {
		// Farm fully operational
		island.Food += farm
		// Remaining population works in factories/mines
		island.Money += min((pop-farm)/10, factory+mountain)
	} else {
		// All population works in farming
		island.Food += pop
	}

	// Food consumption
	island.Food -= int(float64(pop) * hconst.EatenFood)
}

// doCommand executes a command for an island
// Ref: perl/lib/Hako/Turn.pm:543
func doCommand(island *types.Island) int {
	// Phase 1: Simplified implementation
	// Full implementation would process all 28 command types
	// For now, just process basic commands

	if len(island.Commands) == 0 {
		return 1
	}

	com := island.Commands[0]

	switch com.Kind {
	case hconst.ComDoNothing:
		// 資金繰り
		// Ref: perl/lib/Hako/Turn.pm:570-588
		logDoNothing(island.ID, island.Name, hconst.ComName[com.Kind])
		island.Money += 10
		island.Absent++

		// 自動放棄
		if island.Absent >= hconst.GiveupTurn {
			island.Commands[0] = types.Command{
				Kind:   hconst.ComGiveup,
				Target: "",
				X:      0,
				Y:      0,
				Arg:    0,
			}
		}
		return 1

	case hconst.ComPrepare, hconst.ComPrepare2:
		// 整地/地ならし
		// Ref: perl/lib/Hako/Turn.pm:634-650
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		cost := hconst.ComCost[com.Kind]

		if island.Money < cost {
			logNoMoney(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 1
		}

		if land != hconst.LandWaste && land != hconst.LandSea {
			logLandFail(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 1
		}

		// 地ならし処理
		if com.Kind == hconst.ComPrepare2 {
			island.Prepare2++
			// ターン消費せず
			island.Money -= cost

			if land == hconst.LandSea {
				island.Land[x][y] = hconst.LandSea
				island.LandValue[x][y] = 1 // Shallow
				logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			} else {
				island.Land[x][y] = hconst.LandPlains
				island.LandValue[x][y] = 0
				logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			}
			return 0
		}

		// 通常の整地 - 埋蔵金の可能性
		island.Money -= cost

		if land == hconst.LandSea {
			island.Land[x][y] = hconst.LandSea
			island.LandValue[x][y] = 1 // Shallow
			logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
		} else {
			island.Land[x][y] = hconst.LandPlains
			island.LandValue[x][y] = 0
			logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
		}

		// 埋蔵金判定
		if rand.Intn(1000) < hconst.DisMaizo {
			v := 100 + rand.Intn(901)
			island.Money += v
			logMaizo(island.ID, island.Name, hconst.ComName[com.Kind], v)
		}
		return 1

	case hconst.ComReclaim:
		// Reclaim (埋め立て)
		// Ref: perl/lib/Hako/Turn.pm:652-725
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		lv := island.LandValue[x][y]
		cost := hconst.ComCost[com.Kind]

		// Can only reclaim sea, oil field, or submarine base
		if land != hconst.LandSea && land != hconst.LandOil && land != hconst.LandSbase {
			logLandFail(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 0
		}

		// Check if there's land around (all sea means cannot reclaim)
		seaCount := countAround(island.Land, x, y, hconst.LandSea, 7) +
			countAround(island.Land, x, y, hconst.LandOil, 7) +
			countAround(island.Land, x, y, hconst.LandSbase, 7)

		if seaCount == 7 {
			// All around is sea, cannot reclaim
			logNoLandAround(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 0
		}

		if land == hconst.LandSea && lv == 1 {
			// Shallow sea case - turn into wasteland
			island.Land[x][y] = hconst.LandWaste
			island.LandValue[x][y] = 0
			logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			island.Area++

			// If surrounding sea hexes <= 3, turn them into shallow
			if seaCount <= 4 {
				for i := 1; i < 7; i++ {
					sx := x + ax[i]
					sy := y + ay[i]

					// Adjust position based on row parity
					if (sy%2) == 0 && (y%2) == 1 {
						sx--
					}

					if !isValidCoord(sx, sy) {
						continue
					}

					// Turn surrounding sea into shallow
					if island.Land[sx][sy] == hconst.LandSea {
						island.LandValue[sx][sy] = 1
					}
				}
			}
		} else {
			// Regular sea case - turn into shallow
			island.Land[x][y] = hconst.LandSea
			island.LandValue[x][y] = 1
			logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
		}

		// Deduct money
		island.Money -= cost
		return 1

	case hconst.ComDestroy:
		// Destroy (掘削)
		// Ref: perl/lib/Hako/Turn.pm:726-782
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		lv := island.LandValue[x][y]
		cost := hconst.ComCost[com.Kind]

		// Cannot destroy submarine base, oil field, or monster
		if land == hconst.LandSbase || land == hconst.LandOil || land == hconst.LandMonster {
			logLandFail(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 0
		}

		if land == hconst.LandSea && lv == 0 {
			// Sea - search for oil
			// Determine investment amount
			arg := com.Arg
			if arg == 0 {
				arg = 1
			}
			value := min(arg*cost, island.Money)
			str := fmt.Sprintf("%d%s", value, hconst.UnitMoney)
			p := value / cost
			island.Money -= value

			// Check if oil found
			if p > rand.Intn(100) {
				// Oil found
				logOilFound(island.ID, island.Name, x, y, hconst.ComName[com.Kind], str)
				island.Land[x][y] = hconst.LandOil
				island.LandValue[x][y] = 0
			} else {
				// Oil not found
				logOilFail(island.ID, island.Name, x, y, hconst.ComName[com.Kind], str)
			}
			return 1
		}

		// Turn target into sea. Mountain becomes wasteland. Shallow becomes sea.
		if land == hconst.LandMountain {
			island.Land[x][y] = hconst.LandWaste
			island.LandValue[x][y] = 0
		} else if land == hconst.LandSea {
			island.LandValue[x][y] = 0
		} else {
			island.Land[x][y] = hconst.LandSea
			island.LandValue[x][y] = 1
			island.Area--
		}
		logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)

		// Deduct money
		island.Money -= cost
		return 1

	case hconst.ComSellTree:
		// Sell tree (伐採)
		// Ref: perl/lib/Hako/Turn.pm:783-801
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		lv := island.LandValue[x][y]

		// Can only sell trees in forest
		if land != hconst.LandForest {
			logLandFail(island.ID, island.Name, hconst.ComName[com.Kind], x, y)
			return 0
		}

		// Turn target into plains
		island.Land[x][y] = hconst.LandPlains
		island.LandValue[x][y] = 0
		logLandSuc(island.ID, island.Name, hconst.ComName[com.Kind], x, y)

		// Get selling price
		island.Money += hconst.TreeValue * lv
		return 1

	case hconst.ComPlant, hconst.ComFarm, hconst.ComFactory, hconst.ComBase, hconst.ComMonument, hconst.ComHaribote, hconst.ComDbase:
		// Ground construction commands (地上建設系)
		// Ref: perl/lib/Hako/Turn.pm:802-958
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		lv := island.LandValue[x][y]
		cost := hconst.ComCost[com.Kind]
		comName := hconst.ComName[com.Kind]
		point := fmt.Sprintf("(%d,%d)", x, y)

		// Check if money is sufficient
		if island.Money < cost {
			logNoMoney(island.ID, island.Name, comName, x, y)
			return 1
		}

		// Check if terrain is suitable for construction
		// Can build on: plains, town, or same type of facility
		suitable := false
		if land == hconst.LandPlains || land == hconst.LandTown {
			suitable = true
		} else if land == hconst.LandMonument && com.Kind == hconst.ComMonument {
			suitable = true
		} else if land == hconst.LandFarm && com.Kind == hconst.ComFarm {
			suitable = true
		} else if land == hconst.LandFactory && com.Kind == hconst.ComFactory {
			suitable = true
		} else if land == hconst.LandDefence && com.Kind == hconst.ComDbase {
			suitable = true
		}

		if !suitable {
			// Unsuitable terrain
			logLandFail(island.ID, island.Name, comName, x, y)
			return 0
		}

		// Process by command type
		switch com.Kind {
		case hconst.ComPlant:
			// Plant - turn target into forest
			island.Land[x][y] = hconst.LandForest
			island.LandValue[x][y] = 1 // Minimum trees
			logPBSuc(island.ID, island.Name, comName, point)

		case hconst.ComBase:
			// Missile base - turn target into base
			island.Land[x][y] = hconst.LandBase
			island.LandValue[x][y] = 0 // Experience 0
			logPBSuc(island.ID, island.Name, comName, point)

		case hconst.ComHaribote:
			// Haribote - turn target into fake defense
			island.Land[x][y] = hconst.LandHaribote
			island.LandValue[x][y] = 0
			logHariSuc(island.ID, island.Name, comName, hconst.ComName[hconst.ComDbase], point)

		case hconst.ComFarm:
			// Farm
			if land == hconst.LandFarm {
				// Already a farm - expand
				island.LandValue[x][y] += 2 // +2000 people
				if island.LandValue[x][y] > 50 {
					island.LandValue[x][y] = 50 // Maximum 50000 people
				}
			} else {
				// Turn target into farm
				island.Land[x][y] = hconst.LandFarm
				island.LandValue[x][y] = 10 // Scale = 10000 people
			}
			logLandSuc(island.ID, island.Name, comName, x, y)

		case hconst.ComFactory:
			// Factory
			if land == hconst.LandFactory {
				// Already a factory - expand
				island.LandValue[x][y] += 10 // +10000 people
				if island.LandValue[x][y] > 100 {
					island.LandValue[x][y] = 100 // Maximum 100000 people
				}
			} else {
				// Turn target into factory
				island.Land[x][y] = hconst.LandFactory
				island.LandValue[x][y] = 30 // Scale = 30000 people
			}
			logLandSuc(island.ID, island.Name, comName, x, y)

		case hconst.ComDbase:
			// Defense facility
			if land == hconst.LandDefence {
				// Already defense - set self-destruct
				island.LandValue[x][y] = 1 // Self-destruct set
				lName := landName(land, lv)
				logBombSet(island.ID, island.Name, lName, point)
			} else {
				// Turn target into defense facility
				island.Land[x][y] = hconst.LandDefence
				island.LandValue[x][y] = 0
				logLandSuc(island.ID, island.Name, comName, x, y)
			}

		case hconst.ComMonument:
			// Monument
			if land == hconst.LandMonument {
				// Already monument - launch giant missile
				// Get target
				targetID := com.Target
				tn, ok := variable.IDToNumber[targetID]
				if !ok {
					// Target no longer exists - silently abort
					return 0
				}
				tIsland := variable.Islands[tn]
				tIsland.BigMissile++

				// Turn target into wasteland
				island.Land[x][y] = hconst.LandWaste
				island.LandValue[x][y] = 0
				lName := landName(land, lv)
				logMonFly(island.ID, island.Name, lName, point)
			} else {
				// Turn target into monument
				island.Land[x][y] = hconst.LandMonument
				arg := com.Arg
				if arg >= hconst.MonumentNumber {
					arg = 0
				}
				island.LandValue[x][y] = arg
				logLandSuc(island.ID, island.Name, comName, x, y)
			}
		}

		// Deduct money
		island.Money -= cost

		// For repeating commands (farm, factory), put command back
		if com.Kind == hconst.ComFarm || com.Kind == hconst.ComFactory {
			if com.Arg > 1 {
				com.Arg--
				slideBack(island.Commands, 0)
				island.Commands[0] = com
			}
		}

		return 1

	case hconst.ComMountain:
		// Mining site (採掘場整備)
		// Ref: perl/lib/Hako/Turn.pm:959-990
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		cost := hconst.ComCost[com.Kind]
		comName := hconst.ComName[com.Kind]

		// Check if money is sufficient
		if island.Money < cost {
			logNoMoney(island.ID, island.Name, comName, x, y)
			return 1
		}

		// Can only build on mountain
		if land != hconst.LandMountain {
			logLandFail(island.ID, island.Name, comName, x, y)
			return 0
		}

		// Increase scale by 5 (5000 people), max 200 (200000 people)
		island.LandValue[x][y] += 5
		if island.LandValue[x][y] > 200 {
			island.LandValue[x][y] = 200
		}
		logLandSuc(island.ID, island.Name, comName, x, y)

		// Deduct money
		island.Money -= cost

		// For repeating commands, put command back
		if com.Arg > 1 {
			com.Arg--
			slideBack(island.Commands, 0)
			island.Commands[0] = com
		}
		return 1

	case hconst.ComSbase:
		// Submarine base (海底基地建設)
		// Ref: perl/lib/Hako/Turn.pm:991-1008
		x, y := com.X, com.Y
		if !isValidCoord(x, y) {
			return 1
		}

		land := island.Land[x][y]
		lv := island.LandValue[x][y]
		cost := hconst.ComCost[com.Kind]
		comName := hconst.ComName[com.Kind]

		// Check if money is sufficient
		if island.Money < cost {
			logNoMoney(island.ID, island.Name, comName, x, y)
			return 1
		}

		// Can only build on deep sea (sea with lv == 0)
		if land != hconst.LandSea || lv != 0 {
			logLandFail(island.ID, island.Name, comName, x, y)
			return 0
		}

		// Build submarine base
		island.Land[x][y] = hconst.LandSbase
		island.LandValue[x][y] = 0 // Experience 0

		// Log with hidden coordinates (secret base)
		logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで%s%s%sが行われました。",
			variable.IslandTurn, island.ID, hconst.TagNameBegin, island.Name, "(?, ?)", hconst.TagNameEnd,
			hconst.TagComNameBegin, comName, hconst.TagComNameEnd))

		// Deduct money
		island.Money -= cost
		return 1

	case hconst.ComMissileNM, hconst.ComMissilePP, hconst.ComMissileST, hconst.ComMissileLD:
		// Missile commands (ミサイル系)
		// Ref: perl/lib/Hako/Turn.pm:1009-1495
		// Get target
		targetID := com.Target
		tn, ok := variable.IDToNumber[targetID]
		if !ok {
			// Target no longer exists
			logMsNoTarget(island.ID, island.Name, hconst.ComName[com.Kind])
			return 0
		}

		flag := 0
		arg := com.Arg
		if arg == 0 {
			// 0 means shoot as many as possible
			arg = 10000
		}

		// Preparation
		tIsland := variable.Islands[tn]
		tName := tIsland.Name
		tLand := tIsland.Land
		tLandValue := tIsland.LandValue
		var tx, ty, err int

		// Number of refugees
		boat := 0

		// Error range
		if com.Kind == hconst.ComMissilePP {
			err = 7
		} else {
			err = 19
		}

		cost := hconst.ComCost[com.Kind]
		comName := hconst.ComName[com.Kind]
		point := fmt.Sprintf("(%d,%d)", com.X, com.Y)

		// Loop until money runs out, arg is satisfied, or all bases have fired
		var bx, by, count int
		for arg > 0 && island.Money >= cost {
			// Loop until a base is found
			for count < hconst.PointNumber {
				bx = variable.Rpx[count]
				by = variable.Rpy[count]
				if island.Land[bx][by] == hconst.LandBase || island.Land[bx][by] == hconst.LandSbase {
					break
				}
				count++
			}
			if count >= hconst.PointNumber {
				// Not found, exit
				break
			}

			// At least one base found, set flag
			flag = 1

			// Calculate base level
			level := expToLevel(island.Land[bx][by], island.LandValue[bx][by])

			// Loop within the base
			for level > 0 && arg > 0 && island.Money > cost {
				// Confirmed firing, so consume values
				level--
				arg--
				island.Money -= cost

				// Calculate impact point
				r := rand.Intn(err)
				tx = com.X + ax[r]
				ty = com.Y + ay[r]
				if (ty%2) == 0 && (com.Y%2) == 1 {
					tx--
				}

				// Check if impact point is in range
				if tx < 0 || tx >= hconst.IslandSize || ty < 0 || ty >= hconst.IslandSize {
					// Out of range
					if com.Kind == hconst.ComMissileST {
						// Stealth
						logMsOutS(island.ID, targetID, island.Name, tName, comName, point)
					} else {
						// Normal
						logMsOut(island.ID, targetID, island.Name, tName, comName, point)
					}
					continue
				}

				// Get terrain at impact point
				tL := tLand[tx][ty]
				tLv := tLandValue[tx][ty]
				tLname := landName(tL, tLv)
				tPoint := fmt.Sprintf("(%d,%d)", tx, ty)

				// Phase 1: Defense facility judgment simplified (skip)
				// defence := 0

				// Check "no effect" hex first
				if (tL == hconst.LandSea && tLv == 0) || // Deep sea
					((tL == hconst.LandSea || tL == hconst.LandSbase || tL == hconst.LandMountain) &&
						com.Kind != hconst.ComMissileLD) {
					// If submarine base, pretend it's sea
					if tL == hconst.LandSbase {
						tL = hconst.LandSea
					}
					tLname = landName(tL, tLv)

					// Neutralize
					if com.Kind == hconst.ComMissileST {
						// Stealth
						logMsNoDamageS(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					} else {
						// Normal
						logMsNoDamage(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					}
					continue
				}

				// Branch by missile type
				if com.Kind == hconst.ComMissileLD {
					// Land destruction missile
					if tL == hconst.LandMountain {
						// Mountain (becomes wasteland)
						logMsLDMountain(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)

						// Becomes wasteland
						tLand[tx][ty] = hconst.LandWaste
						tLandValue[tx][ty] = 0
						continue
					} else if tL == hconst.LandSbase {
						// Submarine base
						logMsLDSbase(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					} else if tL == hconst.LandMonster {
						// Monster (Phase 1 simplified)
						logMsLDMonster(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					} else if tL == hconst.LandSea {
						// Shallow
						logMsLDSea1(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					} else {
						// Other
						logMsLDLand(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
					}

					// Experience
					if tL == hconst.LandTown {
						if island.Land[bx][by] == hconst.LandBase || island.Land[bx][by] == hconst.LandSbase {
							// Only if still a base
							island.LandValue[bx][by] += tLv / 20
							if island.LandValue[bx][by] > hconst.MaxExpPoint {
								island.LandValue[bx][by] = hconst.MaxExpPoint
							}
						}
					}

					// Becomes shallow
					tLand[tx][ty] = hconst.LandSea
					tIsland.Area--
					tLandValue[tx][ty] = 1

					// But if oil field, shallow, or submarine base, becomes sea
					if tL == hconst.LandOil || tL == hconst.LandSea || tL == hconst.LandSbase {
						tLandValue[tx][ty] = 0
					}
				} else {
					// Other missiles
					if tL == hconst.LandWaste {
						// Wasteland (no damage)
						if com.Kind == hconst.ComMissileST {
							// Stealth
							logMsWasteS(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
						} else {
							// Normal
							logMsWaste(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
						}
					} else if tL == hconst.LandMonster {
						// Monster (Phase 1 simplified - skip processing, just log)
						// Phase 2 will implement monster handling
					} else {
						// Normal terrain
						if com.Kind == hconst.ComMissileST {
							// Stealth
							logMsNormalS(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
						} else {
							// Normal
							logMsNormal(island.ID, targetID, island.Name, tName, comName, tLname, point, tPoint)
						}
					}

					// Experience
					if tL == hconst.LandTown {
						if island.Land[bx][by] == hconst.LandBase || island.Land[bx][by] == hconst.LandSbase {
							island.LandValue[bx][by] += tLv / 20
							boat += tLv // Add to refugees for normal missiles
							if island.LandValue[bx][by] > hconst.MaxExpPoint {
								island.LandValue[bx][by] = hconst.MaxExpPoint
							}
						}
					}

					// Becomes wasteland
					tLand[tx][ty] = hconst.LandWaste
					tLandValue[tx][ty] = 1 // Impact point

					// But if oil field, becomes sea
					if tL == hconst.LandOil {
						tLand[tx][ty] = hconst.LandSea
						tLandValue[tx][ty] = 0
					}
				}
			}

			// Increment count
			count++
		}

		if flag == 0 {
			// No base found
			logMsNoBase(island.ID, island.Name, comName)
			return 0
		}

		// Phase 1: Refugee processing simplified (skip)

		return 1

	case hconst.ComSendMonster:
		// 怪獣派遣
		// Ref: perl/lib/Hako/Turn.pm:1580-1601
		targetID := com.Target
		tn, ok := variable.IDToNumber[targetID]
		if !ok {
			// ターゲットがすでにない
			logMsNoTarget(island.ID, island.Name, hconst.ComName[com.Kind])
			return 0
		}

		tIsland := variable.Islands[tn]
		tName := tIsland.Name

		// メッセージ
		logMonsSend(island.ID, targetID, island.Name, tName)
		tIsland.MonsterSend++

		cost := hconst.ComCost[com.Kind]
		island.Money -= cost
		return 1

	case hconst.ComSell:
		// 食料輸出
		// Ref: perl/lib/Hako/Turn.pm:1602-1613
		cost := hconst.ComCost[com.Kind]
		arg := com.Arg
		if arg == 0 {
			arg = 1
		}
		value := min(arg*(-cost), island.Food)

		// 輸出ログ
		logSell(island.ID, island.Name, hconst.ComName[com.Kind], value)
		island.Food -= value
		island.Money += value / 10
		return 0

	case hconst.ComFood, hconst.ComMoney:
		// 援助系
		// Ref: perl/lib/Hako/Turn.pm:1614-1647
		targetID := com.Target
		tn, ok := variable.IDToNumber[targetID]
		if !ok {
			// ターゲットがない場合はスキップ
			return 0
		}

		tIsland := variable.Islands[tn]
		tName := tIsland.Name
		cost := hconst.ComCost[com.Kind]
		arg := com.Arg
		if arg == 0 {
			arg = 1
		}

		var value int
		var str string
		if cost < 0 {
			// 食料援助
			value = min(arg*(-cost), island.Food)
			str = fmt.Sprintf("%d%s", value, hconst.UnitFood)
		} else {
			// 資金援助
			value = min(arg*cost, island.Money)
			str = fmt.Sprintf("%d%s", value, hconst.UnitMoney)
		}

		// 援助ログ
		logAid(island.ID, targetID, island.Name, tName, hconst.ComName[com.Kind], str)

		if cost < 0 {
			island.Food -= value
			tIsland.Food += value
		} else {
			island.Money -= value
			tIsland.Money += value
		}
		return 0

	case hconst.ComPropaganda:
		// 誘致活動
		// Ref: perl/lib/Hako/Turn.pm:1648-1655
		cost := hconst.ComCost[com.Kind]
		logPropaganda(island.ID, island.Name, hconst.ComName[com.Kind])
		island.Propaganda = 1
		island.Money -= cost
		return 1

	case hconst.ComGiveup:
		// 放棄
		// Ref: perl/lib/Hako/Turn.pm:1656-1663
		logGiveup(island.ID, island.Name)
		island.Dead = true
		os.Remove(fmt.Sprintf("island.%s", island.ID))
		return 1

	default:
		// Other commands: Phase 1 simplified - skip
		return 1
	}
}

// doEachHex processes each hex on an island
// Ref: perl/lib/Hako/Turn.pm:1669-1920
func doEachHex(island *types.Island) {
	// Monster move tracking (for Phase 2)
	monsterMove := make([][]int, hconst.IslandSize)
	for i := 0; i < hconst.IslandSize; i++ {
		monsterMove[i] = make([]int, hconst.IslandSize)
	}

	// Derived values
	name := island.Name
	id := island.ID
	land := island.Land
	landValue := island.LandValue

	// Population growth seed values
	addpop := 10  // Villages and towns
	addpop2 := 0  // Cities
	if island.Food < 0 {
		// Food shortage
		addpop = -30
	} else if island.Propaganda == 1 {
		// Propaganda active
		addpop = 30
		addpop2 = 3
	}

	// Loop through all hexes
	for i := 0; i < hconst.PointNumber; i++ {
		x := variable.Rpx[i]
		y := variable.Rpy[i]
		landKind := land[x][y]
		lv := landValue[x][y]

		if landKind == hconst.LandTown {
			// Town system
			if addpop < 0 {
				// Food shortage
				lv -= (rand.Intn(-addpop) + 1)
				if lv <= 0 {
					// Revert to plains
					land[x][y] = hconst.LandPlains
					landValue[x][y] = 0
					continue
				}
			} else {
				// Growth
				if lv < 100 {
					lv += rand.Intn(addpop) + 1
					if lv > 100 {
						lv = 100
					}
				} else {
					// Cities grow slower
					if addpop2 > 0 {
						lv += rand.Intn(addpop2) + 1
					}
				}
			}
			if lv > 200 {
				lv = 200
			}
			landValue[x][y] = lv

		} else if landKind == hconst.LandPlains {
			// Plains
			if rand.Intn(5) == 0 {
				// If farms or towns around, this becomes a town
				if countGrow(land, landValue, x, y) {
					land[x][y] = hconst.LandTown
					landValue[x][y] = 1
				}
			}

		} else if landKind == hconst.LandForest {
			// Forest
			if lv < 200 {
				// Grow trees
				landValue[x][y]++
			}

		} else if landKind == hconst.LandDefence {
			// Defense facility
			if lv == 1 {
				// Self-destruct
				lName := landName(landKind, lv)
				logBombFire(id, name, lName, fmt.Sprintf("(%d,%d)", x, y))

				// Wide damage routine (Phase 1: Simplified - just make it wasteland)
				land[x][y] = hconst.LandWaste
				landValue[x][y] = 0
			}

		} else if landKind == hconst.LandOil {
			// Submarine oil field
			lName := landName(landKind, lv)
			value := hconst.OilMoney
			island.Money += value
			str := fmt.Sprintf("%d%s", value, hconst.UnitMoney)

			// Income log
			logOilMoney(id, name, lName, fmt.Sprintf("(%d,%d)", x, y), str)

			// Depletion check
			if rand.Intn(1000) < hconst.OilRatio {
				// Depleted
				logOilEnd(id, name, lName, fmt.Sprintf("(%d,%d)", x, y))
				land[x][y] = hconst.LandSea
				landValue[x][y] = 0
			}

		} else if landKind == hconst.LandMonster {
			// Monster (Phase 1: Skip monster movement processing)
			// Phase 2 will implement full monster movement
		}

		// Fire check
		if ((landKind == hconst.LandTown) && (lv > 30)) ||
			(landKind == hconst.LandHaribote) ||
			(landKind == hconst.LandFactory) {
			if rand.Intn(1000) < hconst.DisFire {
				// Count surrounding forests and monuments
				if (countAround(land, x, y, hconst.LandForest, 7) +
					countAround(land, x, y, hconst.LandMonument, 7)) == 0 {
					// No forests or monuments - fire destroys everything
					l := land[x][y]
					lv := landValue[x][y]
					point := fmt.Sprintf("(%d,%d)", x, y)
					lName := landName(l, lv)
					logFire(id, name, lName, point)
					land[x][y] = hconst.LandWaste
					landValue[x][y] = 0
				}
			}
		}
	}
}

// doIslandProcess processes whole island events
// Ref: perl/lib/Hako/Turn.pm:1956
func doIslandProcess(num int, island *types.Island) {
	// Ref: perl/lib/Hako/Turn.pm:1956-2428
	name := island.Name
	id := island.ID
	land := island.Land
	landValue := island.LandValue

	// 地震判定
	if rand.Intn(1000) < ((island.Prepare2+1)*hconst.DisEarthquake) {
		// 地震発生
		logEarthquake(id, name)

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]
			lv := landValue[x][y]

			if ((landKind == hconst.LandTown) && (lv >= 100)) ||
				(landKind == hconst.LandHaribote) ||
				(landKind == hconst.LandFactory) {
				// 1/4で壊滅
				if rand.Intn(4) == 0 {
					logEQDamage(id, name, getLandName(landKind, lv), fmt.Sprintf("(%d, %d)", x, y))
					land[x][y] = hconst.LandWaste
					landValue[x][y] = 0
				}
			}
		}
	}

	// 食料不足
	if island.Food <= 0 {
		// 不足メッセージ
		logStarve(id, name, 0)
		island.Food = 0

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]
			lv := landValue[x][y]

			if (landKind == hconst.LandFarm) ||
				(landKind == hconst.LandFactory) ||
				(landKind == hconst.LandBase) ||
				(landKind == hconst.LandDefence) {
				// 1/4で壊滅
				if rand.Intn(4) == 0 {
					logSvDamage(id, name, getLandName(landKind, lv), fmt.Sprintf("(%d, %d)", x, y))
					land[x][y] = hconst.LandWaste
					landValue[x][y] = 0
				}
			}
		}
	}

	// 津波判定
	if rand.Intn(1000) < hconst.DisTsunami {
		// 津波発生
		logTsunami(id, name)

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]
			lv := landValue[x][y]

			if (landKind == hconst.LandTown) ||
				(landKind == hconst.LandFarm) ||
				(landKind == hconst.LandFactory) ||
				(landKind == hconst.LandBase) ||
				(landKind == hconst.LandDefence) ||
				(landKind == hconst.LandHaribote) {
				// 1d12 <= (周囲の海 - 1) で崩壊
				seaCount := countAround(land, x, y, hconst.LandOil, 7) +
					countAround(land, x, y, hconst.LandSbase, 7) +
					countAround(land, x, y, hconst.LandSea, 7)
				if rand.Intn(12) < (seaCount - 1) {
					logTsunamiDamage(id, name, getLandName(landKind, lv), fmt.Sprintf("(%d, %d)", x, y))
					land[x][y] = hconst.LandWaste
					landValue[x][y] = 0
				}
			}
		}
	}

	// 怪獣判定
	r := rand.Intn(10000)
	pop := island.Pop
	for {
		if ((r < (hconst.DisMonster * island.Area)) && (pop >= hconst.DisMonsBorder1)) ||
			(island.MonsterSend > 0) {
			// 怪獣出現
			// 種類を決める
			var lv, kind int
			if island.MonsterSend > 0 {
				// 人造
				kind = 0
				island.MonsterSend--
			} else if pop >= hconst.DisMonsBorder3 {
				// level3まで
				kind = rand.Intn(hconst.MonsterLevel3) + 1
			} else if pop >= hconst.DisMonsBorder2 {
				// level2まで
				kind = rand.Intn(hconst.MonsterLevel2) + 1
			} else {
				// level1のみ
				kind = rand.Intn(hconst.MonsterLevel1) + 1
			}

			// lvの値を決める
			lv = kind*10 + hconst.MonsterBHP[kind] + rand.Intn(hconst.MonsterDHP[kind])

			// どこに現れるか決める
			for i := 0; i < hconst.PointNumber; i++ {
				bx := variable.Rpx[i]
				by := variable.Rpy[i]
				if land[bx][by] == hconst.LandTown {
					// 地形名
					lName := getLandName(hconst.LandTown, landValue[bx][by])

					// そのヘックスを怪獣に
					land[bx][by] = hconst.LandMonster
					landValue[bx][by] = lv

					// 怪獣情報
					_, mName, _ := monsterSpec(lv)

					// メッセージ
					logMonsCome(id, name, mName, fmt.Sprintf("(%d, %d)", bx, by), lName)
					break
				}
			}
		}

		if island.MonsterSend <= 0 {
			break
		}
	}

	// 地盤沈下判定
	if (island.Area > hconst.DisFallBorder) && (rand.Intn(1000) < hconst.DisFalldown) {
		// 地盤沈下発生
		logFalldown(id, name)

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]
			lv := landValue[x][y]

			if (landKind != hconst.LandSea) &&
				(landKind != hconst.LandSbase) &&
				(landKind != hconst.LandOil) &&
				(landKind != hconst.LandMountain) {
				// 周囲に海があれば、値を-1に
				if countAround(land, x, y, hconst.LandSea, 7)+
					countAround(land, x, y, hconst.LandSbase, 7) > 0 {
					logFalldownLand(id, name, getLandName(landKind, lv), fmt.Sprintf("(%d, %d)", x, y))
					land[x][y] = -1
					landValue[x][y] = 0
				}
			}
		}

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]

			if landKind == -1 {
				// -1になっている所を浅瀬に
				land[x][y] = hconst.LandSea
				landValue[x][y] = 1
			} else if landKind == hconst.LandSea {
				// 浅瀬は海に
				landValue[x][y] = 0
			}
		}
	}

	// 台風判定
	if rand.Intn(1000) < hconst.DisTyphoon {
		// 台風発生
		logTyphoon(id, name)

		for i := 0; i < hconst.PointNumber; i++ {
			x := variable.Rpx[i]
			y := variable.Rpy[i]
			landKind := land[x][y]
			lv := landValue[x][y]

			if (landKind == hconst.LandFarm) || (landKind == hconst.LandHaribote) {
				// 1d12 <= (6 - 周囲の森) で崩壊
				forestCount := countAround(land, x, y, hconst.LandForest, 7) +
					countAround(land, x, y, hconst.LandMonument, 7)
				if rand.Intn(12) < (6 - forestCount) {
					logTyphoonDamage(id, name, getLandName(landKind, lv), fmt.Sprintf("(%d, %d)", x, y))
					land[x][y] = hconst.LandPlains
					landValue[x][y] = 0
				}
			}
		}
	}

	// 巨大隕石判定
	if rand.Intn(1000) < hconst.DisHugeMeteo {
		// 落下
		x := rand.Intn(hconst.IslandSize)
		y := rand.Intn(hconst.IslandSize)
		point := fmt.Sprintf("(%d, %d)", x, y)

		// メッセージ
		logHugeMeteo(id, name, point)

		// 広域被害ルーチン
		wideDamage(id, name, land, landValue, x, y)
	}

	// 巨大ミサイル判定
	for island.BigMissile > 0 {
		island.BigMissile--

		// 落下
		x := rand.Intn(hconst.IslandSize)
		y := rand.Intn(hconst.IslandSize)
		point := fmt.Sprintf("(%d, %d)", x, y)

		// メッセージ
		logMonDamage(id, name, point)

		// 広域被害ルーチン
		wideDamage(id, name, land, landValue, x, y)
	}

	// 隕石判定
	if rand.Intn(1000) < hconst.DisMeteo {
		first := true
		for (rand.Intn(2) == 0) || first {
			first = false

			// 落下
			x := rand.Intn(hconst.IslandSize)
			y := rand.Intn(hconst.IslandSize)
			landKind := land[x][y]
			lv := landValue[x][y]
			point := fmt.Sprintf("(%d, %d)", x, y)

			if (landKind == hconst.LandSea) && (lv == 0) {
				// 海ポチャ
				logMeteoSea(id, name, getLandName(landKind, lv), point)
			} else if landKind == hconst.LandMountain {
				// 山破壊
				logMeteoMountain(id, name, getLandName(landKind, lv), point)
				land[x][y] = hconst.LandWaste
				landValue[x][y] = 0
				continue
			} else if landKind == hconst.LandSbase {
				logMeteoSbase(id, name, getLandName(landKind, lv), point)
			} else if landKind == hconst.LandMonster {
				logMeteoMonster(id, name, getLandName(landKind, lv), point)
			} else if landKind == hconst.LandSea {
				// 浅瀬
				logMeteoSea1(id, name, getLandName(landKind, lv), point)
			} else {
				logMeteoNormal(id, name, getLandName(landKind, lv), point)
			}
			land[x][y] = hconst.LandSea
			landValue[x][y] = 0
		}
	}

	// 噴火判定
	if rand.Intn(1000) < hconst.DisEruption {
		x := rand.Intn(hconst.IslandSize)
		y := rand.Intn(hconst.IslandSize)
		landKind := land[x][y]
		lv := landValue[x][y]
		point := fmt.Sprintf("(%d, %d)", x, y)
		logEruption(id, name, getLandName(landKind, lv), point)
		land[x][y] = hconst.LandMountain
		landValue[x][y] = 0

		for i := 1; i < 7; i++ {
			sx := x + hconst.Ax[i]
			sy := y + hconst.Ay[i]

			// 行による位置調整
			if ((sy % 2) == 0) && ((y % 2) == 1) {
				sx--
			}

			if (sx < 0) || (sx >= hconst.IslandSize) || (sy < 0) || (sy >= hconst.IslandSize) {
				continue
			}

			// 範囲内の場合
			landKind = land[sx][sy]
			lv = landValue[sx][sy]
			point = fmt.Sprintf("(%d, %d)", sx, sy)

			if (landKind == hconst.LandSea) || (landKind == hconst.LandOil) || (landKind == hconst.LandSbase) {
				// 海の場合
				if lv == 1 {
					// 浅瀬
					logEruptionSea1(id, name, getLandName(landKind, lv), point)
				} else {
					logEruptionSea(id, name, getLandName(landKind, lv), point)
					land[sx][sy] = hconst.LandSea
					landValue[sx][sy] = 1
					continue
				}
			} else if (landKind == hconst.LandMountain) ||
				(landKind == hconst.LandMonster) ||
				(landKind == hconst.LandWaste) {
				continue
			} else {
				// それ以外の場合
				logEruptionNormal(id, name, getLandName(landKind, lv), point)
			}
			land[sx][sy] = hconst.LandWaste
			landValue[sx][sy] = 0
		}
	}

	// 食料があふれてたら換金
	if island.Food > 9999 {
		island.Money += (island.Food - 9999) / 10
		island.Food = 9999
	}

	// 金があふれてたら切り捨て
	if island.Money > 9999 {
		island.Money = 9999
	}

	// 各種の値を計算
	estimate(num)

	// 繁栄、災難賞
	pop = island.Pop
	damage := island.OldPop - pop
	flags := island.Prize

	// 繁栄賞
	if ((flags & 1) == 0) && pop >= 3000 {
		flags |= 1
		logPrize(id, name, hconst.Prize[1])
	} else if ((flags & 2) == 0) && pop >= 5000 {
		flags |= 2
		logPrize(id, name, hconst.Prize[2])
	} else if ((flags & 4) == 0) && pop >= 10000 {
		flags |= 4
		logPrize(id, name, hconst.Prize[3])
	}

	// 災難賞
	if ((flags & 64) == 0) && damage >= 500 {
		flags |= 64
		logPrize(id, name, hconst.Prize[7])
	} else if ((flags & 128) == 0) && damage >= 1000 {
		flags |= 128
		logPrize(id, name, hconst.Prize[8])
	} else if ((flags & 256) == 0) && damage >= 2000 {
		flags |= 256
		logPrize(id, name, hconst.Prize[9])
	}

	island.Prize = flags

	// Absent counter
	if island.Pop > 0 {
		island.Absent = 0
	} else {
		island.Absent++
		if island.Absent > hconst.GiveupTurn {
			// Giveup
			logGiveup(island.ID, island.Name)
			island.Dead = true
		}
	}
}

// islandSort sorts islands by population
// Ref: perl/lib/Hako/Turn.pm:2431
func islandSort() {
	// Simple bubble sort (Phase 1)
	for i := 0; i < variable.IslandNumber-1; i++ {
		for j := 0; j < variable.IslandNumber-i-1; j++ {
			if variable.Islands[j].Pop < variable.Islands[j+1].Pop {
				variable.Islands[j], variable.Islands[j+1] = variable.Islands[j+1], variable.Islands[j]
			}
		}
	}
}

// estimate calculates derived statistics for an island
// Ref: perl/lib/Hako/Main.pm (estimate function)
func estimate(num int) {
	island := variable.Islands[num]

	// Count terrain
	pop := 0
	area := 0
	farm := 0
	factory := 0
	mountain := 0

	for y := 0; y < hconst.IslandSize; y++ {
		for x := 0; x < hconst.IslandSize; x++ {
			land := island.Land[x][y]
			lv := island.LandValue[x][y]

			if land != hconst.LandSea {
				area++
			}

			if land == hconst.LandTown {
				pop += lv
			} else if land == hconst.LandFarm {
				farm += lv
			} else if land == hconst.LandFactory {
				factory += lv
			} else if land == hconst.LandMountain && lv > 0 {
				mountain += lv
			}
		}
	}

	island.Pop = pop
	island.Area = area
	island.Farm = farm
	island.Factory = factory
	island.Mountain = mountain
}

// makeRandomPointArray creates shuffled coordinate arrays
// Sets variable.Rpx, variable.Rpy so that numbers from (0,0) to (size-1, size-1) appear exactly once
// Ref: perl/lib/Hako/Main.pm:1002
func makeRandomPointArray() {
	// Initialize
	variable.Rpx = make([]int, hconst.PointNumber)
	variable.Rpy = make([]int, hconst.PointNumber)

	idx := 0
	for y := 0; y < hconst.IslandSize; y++ {
		for x := 0; x < hconst.IslandSize; x++ {
			variable.Rpx[idx] = x
			variable.Rpy[idx] = y
			idx++
		}
	}

	// Shuffle
	for i := hconst.PointNumber - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		if i == j {
			continue
		}
		variable.Rpx[i], variable.Rpx[j] = variable.Rpx[j], variable.Rpx[i]
		variable.Rpy[i], variable.Rpy[j] = variable.Rpy[j], variable.Rpy[i]
	}
}

// randomArray creates a random permutation of 0..n-1
func randomArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}

	// Fisher-Yates shuffle
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

// countAround counts terrain type around a hex
// Ref: perl/lib/Hako/Main.pm
func countAround(land [][]int, x, y, landType, radius int) int {
	count := 0
	for i := 0; i < radius; i++ {
		nx := x + ax[i]
		ny := y + ay[i]
		if isValidCoord(nx, ny) && land[nx][ny] == landType {
			count++
		}
	}
	return count
}

// countGrow checks if there are farms or towns around a hex
// Ref: perl/lib/Hako/Turn.pm:1922-1953
func countGrow(land [][]int, landValue [][]int, x, y int) bool {
	for i := 1; i < 7; i++ {
		sx := x + ax[i]
		sy := y + ay[i]

		// Adjust position based on row parity
		if (sy%2) == 0 && (y%2) == 1 {
			sx--
		}

		// Range check
		if sx < 0 || sx >= hconst.IslandSize || sy < 0 || sy >= hconst.IslandSize {
			continue
		}

		// Check if town or farm
		if (land[sx][sy] == hconst.LandTown) || (land[sx][sy] == hconst.LandFarm) {
			if landValue[sx][sy] != 1 {
				return true
			}
		}
	}
	return false
}

func isValidCoord(x, y int) bool {
	return x >= 0 && x < hconst.IslandSize && y >= 0 && y < hconst.IslandSize
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// expToLevel calculates base level from experience
// Ref: perl/lib/Hako/Main.pm:975
func expToLevel(kind, exp int) int {
	if kind == hconst.LandBase {
		// Missile base
		for i := hconst.MaxBaseLevel; i > 1; i-- {
			if exp >= hconst.BaseLevelUp[i-2] {
				return i
			}
		}
		return 1
	} else {
		// Submarine base
		for i := hconst.MaxSBaseLevel; i > 1; i-- {
			if exp >= hconst.SBaseLevelUp[i-2] {
				return i
			}
		}
		return 1
	}
}

func nameToNumber(name string) int {
	for i := 0; i < variable.IslandNumber; i++ {
		if variable.Islands[i].Name == name {
			return i
		}
	}
	return -1
}

// slideBack shifts commands backward to make room for a new command at position number
// Ref: perl/lib/Hako/Main.pm:182
func slideBack(commands []types.Command, number int) {
	if number == len(commands)-1 {
		return
	}
	// Shift elements from number to len-2, one position forward
	copy(commands[number+1:], commands[number:len(commands)-1])
}

// landName returns the name of a land type
// Ref: perl/lib/Hako/Turn.pm:3445
func landName(land, lv int) string {
	switch land {
	case hconst.LandSea:
		if lv == 1 {
			return "浅瀬"
		}
		return "海"
	case hconst.LandWaste:
		return "荒地"
	case hconst.LandPlains:
		return "平地"
	case hconst.LandTown:
		if lv < 30 {
			return "村"
		} else if lv < 100 {
			return "町"
		}
		return "都市"
	case hconst.LandForest:
		return "森"
	case hconst.LandFarm:
		return "農場"
	case hconst.LandFactory:
		return "工場"
	case hconst.LandBase:
		return "ミサイル基地"
	case hconst.LandDefence:
		return "防衛施設"
	case hconst.LandMountain:
		return "山"
	case hconst.LandMonument:
		return "記念碑"
	case hconst.LandHaribote:
		return "ハリボテ"
	case hconst.LandSbase:
		return "海底基地"
	case hconst.LandOil:
		return "海底油田"
	default:
		return "不明"
	}
}

func encode(password string) string {
	// Phase 1: Simple encoding (same as core.encode but local)
	if hconst.CryptOn {
		return cryptCompat(password, "h2")
	}
	return password
}

func cryptCompat(password, salt string) string {
	// Phase 1: Placeholder
	return password
}

//======================================================================
// Log Functions
//======================================================================

var logBuffer []string

func logOut(s string) {
	logBuffer = append(logBuffer, s)
}

func logFlush() {
	filename := filepath.Join(hconst.DirName, "hakojima.log0")
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, line := range logBuffer {
		fmt.Fprintln(file, line)
	}

	logBuffer = nil
}

func logHistoryTrim() {
	// Phase 1: Stub
}

// Log message functions (simplified for Phase 1)

func logDiscover(name string) {
	logHistory(fmt.Sprintf("%d,%s島が発見される。", variable.IslandTurn, name))
}

func logChangeName(oldName, newName string) {
	logOut(fmt.Sprintf("0,%d,0,0,%s島が%s島に名称変更されました。", variable.IslandTurn, oldName, newName))
}

func logStarve(id, name string, starve int) {
	logOut(fmt.Sprintf("0,%d,0,0,%s島で食料不足により%d%s餓死。", variable.IslandTurn, name, starve, hconst.UnitPop))
}

func logDead(id, name string) {
	logOut(fmt.Sprintf("0,%d,0,0,%s島が放棄され、無人島になりました。", variable.IslandTurn, name))
}

func logGiveup(id, name string) {
	logOut(fmt.Sprintf("0,%d,0,0,%s島が放棄されました。", variable.IslandTurn, name))
}

func logPrize(id, name, prize string) {
	logOut(fmt.Sprintf("0,%d,0,0,%s島が%sを受賞しました。", variable.IslandTurn, name, prize))
}

func logNoMoney(id, name, command string, x, y int) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s島(%d,%d)で予定されていた%sは、資金不足により中止されました。",
		variable.IslandTurn, id, name, x, y, command))
}

func logLandFail(id, name, command string, x, y int) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s島(%d,%d)で予定されていた%sは、対象の地形が適していないため中止されました。",
		variable.IslandTurn, id, name, x, y, command))
}

func logNoLandAround(id, name, command string, x, y int) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s島(%d,%d)で予定されていた%sは、周囲に陸地がないため中止されました。",
		variable.IslandTurn, id, name, x, y, command))
}

func logLandSuc(id, name, command string, x, y int) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s島(%d,%d)で%sが行われました。",
		variable.IslandTurn, id, name, x, y, command))
}

func logOilFound(id, name string, x, y int, command, str string) {
	point := fmt.Sprintf("(%d,%d)", x, y)
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで<B>%s</B>の予算をつぎ込んだ%s%s%sが行われ、<B>油田が掘り当てられました</B>。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		str, hconst.TagComNameBegin, command, hconst.TagComNameEnd))
}

func logOilFail(id, name string, x, y int, command, str string) {
	point := fmt.Sprintf("(%d,%d)", x, y)
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで<B>%s</B>の予算をつぎ込んだ%s%s%sが行われましたが、油田は見つかりませんでした。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		str, hconst.TagComNameBegin, command, hconst.TagComNameEnd))
}

// Disaster log functions

func logEarthquake(id, name string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで大規模な%s地震%sが発生！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logEQDamage(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は%s地震%sにより壊滅しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logSvDamage(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>に<B>食料を求めて住民が殺到</B>。<B>%s</B>は壊滅しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, lName))
}

func logTsunami(id, name string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s付近で%s津波%s発生！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logTsunamiDamage(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は%s津波%sにより崩壊しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logMonsCome(id, name, mName, point, lName string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sに<B>怪獣%s</B>出現！！%s%s%sの<B>%s</B>が踏み荒らされました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		mName, hconst.TagNameBegin, point, hconst.TagNameEnd, lName))
}

func logFalldown(id, name string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで%s地盤沈下%sが発生しました！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logFalldownLand(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は海の中へ沈みました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logTyphoon(id, name string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sに%s台風%s上陸！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logTyphoonDamage(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は%s台風%sで飛ばされました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logHugeMeteo(id, name, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点に%s巨大隕石%sが落下！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logMonDamage(id, name, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,<B>何かとてつもないもの</B>が%s%s島%s%s地点に落下しました！！",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd))
}

func logMeteoSea(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>に%s隕石%sが落下しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logMeteoMountain(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>に%s隕石%sが落下、<B>%s</B>は消し飛びました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd, lName))
}

func logMeteoSbase(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>に%s隕石%sが落下、<B>%s</B>は崩壊しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd, lName))
}

func logMeteoMonster(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,<B>怪獣%s</B>がいた%s%s島%s%s地点に%s隕石%sが落下、陸地は<B>怪獣%s</B>もろとも水没しました。",
		variable.IslandTurn, id, lName, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd, lName))
}

func logMeteoSea1(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点に%s隕石%sが落下、海底がえぐられました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logMeteoNormal(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点の<B>%s</B>に%s隕石%sが落下、一帯が水没しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logEruption(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点で%s火山が噴火%s、<B>山</B>が出来ました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logEruptionSea1(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点の<B>%s</B>は、%s噴火%sの影響で陸地になりました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logEruptionSea(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点の<B>%s</B>は、%s噴火%sの影響で海底が隆起、浅瀬になりました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logEruptionNormal(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%s地点の<B>%s</B>は、%s噴火%sの影響で壊滅しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		lName, hconst.TagDisasterBegin, hconst.TagDisasterEnd))
}

func logWideDamageSea(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は<B>水没</B>しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logWideDamageSea2(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は跡形もなくなりました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logWideDamageMonsterSea(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの陸地は<B>怪獣%s</B>もろとも水没しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logWideDamageMonster(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>怪獣%s</B>は消し飛びました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logWideDamageWaste(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は一瞬にして<B>荒地</B>と化しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

// logPBSuc logs planting or base construction success
// Ref: perl/lib/Hako/Turn.pm:2739
func logPBSuc(id, name, comName, point string) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで%s%s%sが行われました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
	logOut(fmt.Sprintf("0,%d,%s,0,こころなしか、%s%s島%sの<B>森</B>が増えたようです。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd))
}

// logHariSuc logs haribote (fake defense) success
// Ref: perl/lib/Hako/Turn.pm:2754
func logHariSuc(id, name, comName, comName2, point string) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで%s%s%sが行われました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sで%s%s%sが行われました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName2, hconst.TagComNameEnd))
}

// logBombSet logs self-destruct device set
// Ref: perl/lib/Hako/Turn.pm:2697
func logBombSet(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>の<B>自爆装置がセット</B>されました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

// logMonFly logs monument (giant missile) launch
// Ref: perl/lib/Hako/Turn.pm:2719
func logMonFly(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>が<B>轟音とともに飛び立ちました</B>。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

func logSecret(msg string) {
	logOut("1," + msg)
}

func logHistory(msg string) {
	filename := filepath.Join(hconst.DirName, "hakojima.his")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintln(file, msg)
}

// Missile log functions
// Ref: perl/lib/Hako/Turn.pm:2766-3026

// logMsNoTarget logs that missile target has no people
func logMsNoTarget(id, name, comName string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで予定されていた%s%s%sは、目標の島に人が見当たらないため中止されました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logMsNoBase logs that missile launch failed due to no base
func logMsNoBase(id, name, comName string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで予定されていた%s%s%sは、<B>ミサイル設備を保有していない</B>ために実行できませんでした。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logMsOut logs that missile fell out of range
func logMsOut(id, tID, name, tName, comName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、<B>領域外の海</B>に落ちた模様です。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logMsOutS logs that stealth missile fell out of range
func logMsOutS(id, tID, name, tName, comName, point string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、<B>領域外の海</B>に落ちた模様です。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logMsCaught logs that missile was caught by defense facility
func logMsCaught(id, tID, name, tName, comName, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%s地点上空にて力場に捉えられ、<B>空中爆発</B>しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd))
}

// logMsCaughtS logs that stealth missile was caught by defense facility
func logMsCaughtS(id, tID, name, tName, comName, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%s地点上空にて力場に捉えられ、<B>空中爆発</B>しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd))
}

// logMsNoDamage logs that missile had no effect
func logMsNoDamage(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%sの<B>%s</B>に落ちたので被害がありませんでした。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsNoDamageS logs that stealth missile had no effect
func logMsNoDamageS(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%sの<B>%s</B>に落ちたので被害がありませんでした。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsLDMountain logs land destruction missile hit mountain
func logMsLDMountain(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>%s</B>に命中。<B>%s</B>は消し飛び、荒地と化しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname, tLname))
}

// logMsLDSbase logs land destruction missile hit submarine base
func logMsLDSbase(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sに着水後爆発、同地点にあった<B>%s</B>は跡形もなく吹き飛びました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsLDMonster logs land destruction missile hit monster
func logMsLDMonster(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sに着弾し爆発。陸地は<B>怪獣%s</B>もろとも水没しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsLDSea1 logs land destruction missile hit shallow sea
func logMsLDSea1(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>%s</B>に着弾。海底がえぐられました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsLDLand logs land destruction missile hit other terrain
func logMsLDLand(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>%s</B>に着弾。陸地は水没しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsWaste logs missile hit wasteland
func logMsWaste(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%sの<B>%s</B>に落ちました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsWasteS logs stealth missile hit wasteland
func logMsWasteS(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行いましたが、%s%s%sの<B>%s</B>に落ちました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsMonNoDamage logs missile hit hardened monster (no damage)
func logMsMonNoDamage(id, tID, name, tName, comName, mName, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中、しかし硬化状態だったため効果がありませんでした。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName))
}

// logMsMonNoDamageS logs stealth missile hit hardened monster (no damage)
func logMsMonNoDamageS(id, tID, name, tName, comName, mName, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中、しかし硬化状態だったため効果がありませんでした。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName))
}

// logMsMonKill logs missile killed monster
func logMsMonKill(id, tID, name, tName, comName, mName, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中。<B>怪獣%s</B>は力尽き、倒れました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName, mName))
}

// logMsMonKillS logs stealth missile killed monster
func logMsMonKillS(id, tID, name, tName, comName, mName, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中。<B>怪獣%s</B>は力尽き、倒れました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName, mName))
}

// logMsMonster logs missile damaged monster
func logMsMonster(id, tID, name, tName, comName, mName, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中。<B>怪獣%s</B>は苦しそうに咆哮しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName, mName))
}

// logMsMonsterS logs stealth missile damaged monster
func logMsMonsterS(id, tID, name, tName, comName, mName, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>怪獣%s</B>に命中。<B>怪獣%s</B>は苦しそうに咆哮しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, mName, mName))
}

// logMsMonMoney logs monster corpse value
func logMsMonMoney(tID, mName string, value int) {
	logOut(fmt.Sprintf("0,%d,%s,0,<B>怪獣%s</B>の残骸には、<B>%d%s</B>の値が付きました。",
		variable.IslandTurn, tID, mName, value, hconst.UnitMoney))
}

// logMsNormal logs missile hit normal terrain
func logMsNormal(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>%s</B>に命中、一帯が壊滅しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logMsNormalS logs stealth missile hit normal terrain
func logMsNormalS(id, tID, name, tName, comName, tLname, point, tPoint string) {
	logSecret(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%s%s地点に向けて%s%s%sを行い、%s%s%sの<B>%s</B>に命中、一帯が壊滅しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, point, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		hconst.TagNameBegin, tPoint, hconst.TagNameEnd, tLname))
}

// logDoNothing logs do nothing command
// Ref: perl/lib/Hako/Turn.pm:2646
func logDoNothing(id, name, comName string) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%sで%s%s%sが行われました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logMaizo logs buried treasure discovery
// Ref: perl/lib/Hako/Turn.pm:2661
func logMaizo(id, name, comName string, value int) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで予定されていた%s%s%sが、<B>埋蔵金</B>を発見。<B>%d%s</B>の収益を得ました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd, value, hconst.UnitMoney))
}

// logMonsSend logs monster send
// Ref: perl/lib/Hako/Turn.pm:3027
func logMonsSend(id, tID, name, tName string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%sへ<B>怪獣を派遣</B>しました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, hconst.TagNameEnd))
}

// logSell logs food export
// Ref: perl/lib/Hako/Turn.pm:3045
func logSell(id, name, comName string, value int) {
	logSecret(fmt.Sprintf("0,%d,%s,0,%s%s島%sで%s%s%sが行われ、<B>%d%s</B>の食料が輸出され、<B>%d%s</B>の収益が得られました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd,
		value, hconst.UnitFood, value/10, hconst.UnitMoney))
}

// logAid logs aid (food or money)
// Ref: perl/lib/Hako/Turn.pm:3058
func logAid(id, tID, name, tName, comName, str string) {
	logOut(fmt.Sprintf("0,%d,%s,%s,%s%s島%sが%s%s島%sへ%s%s%sを行い、<B>%s</B>の援助を行いました。",
		variable.IslandTurn, id, tID,
		hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagNameBegin, tName, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd, str))
}

// logPropaganda logs propaganda
// Ref: perl/lib/Hako/Turn.pm:3073
func logPropaganda(id, name, comName string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%sで%s%s%sが行われました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, hconst.TagNameEnd,
		hconst.TagComNameBegin, comName, hconst.TagComNameEnd))
}

// logBombFire logs defense facility self-destruct
// Ref: perl/lib/Hako/Turn.pm:2681
func logBombFire(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>が<B>自爆装置の発動</B>により爆発しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

// logOilMoney logs oil field income
// Ref: perl/lib/Hako/Turn.pm:3084
func logOilMoney(id, name, lName, point, str string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>から、<B>%s</B>の収益が上がりました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName, str))
}

// logOilEnd logs oil field depletion
// Ref: perl/lib/Hako/Turn.pm:3097
func logOilEnd(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>は、<B>枯渇</B>したようです。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

// logFire logs fire disaster
// Ref: perl/lib/Hako/Turn.pm:3421
func logFire(id, name, lName, point string) {
	logOut(fmt.Sprintf("0,%d,%s,0,%s%s島%s%sの<B>%s</B>が<B>火災</B>により壊滅しました。",
		variable.IslandTurn, id, hconst.TagNameBegin, name, point, hconst.TagNameEnd, lName))
}

//======================================================================
// Template Functions
//======================================================================

func tempNewIslandFull() {
	out("<H1>島が一杯です。</H1>\n")
}

func tempNewIslandNoName() {
	out("<H1>島の名前を入力してください。</H1>\n")
}

func tempNewIslandBadName() {
	out("<H1>その名前は使用できません。</H1>\n")
}

func tempNewIslandAlready() {
	out("<H1>その名前は既に使われています。</H1>\n")
}

func tempNewIslandNoPassword() {
	out("<H1>パスワードを入力してください。</H1>\n")
}

func tempNewIslandHead(name string) {
	out(fmt.Sprintf(`<CENTER>
<H1>%s%s島%sを発見しました！！</H1>
%s<BR>
</CENTER>
`,
		hconst.TagBigBegin, name, hconst.TagBigEnd,
		hconst.TempBack,
	))
}

func tempChange() {
	out("<H1>変更が完了しました。</H1>\n")
}

func tempChangeNoMoney() {
	out("<H1>資金が不足しています。</H1>\n")
}

func tempChangeNothing() {
	out("<H1>変更する項目がありません。</H1>\n")
}

// getLandName returns the display name for a land type
// Ref: perl/lib/Hako/Turn.pm:3445
func getLandName(land, lv int) string {
	switch land {
	case hconst.LandSea:
		if lv == 1 {
			return "浅瀬"
		}
		return "海"
	case hconst.LandWaste:
		return "荒地"
	case hconst.LandPlains:
		return "平地"
	case hconst.LandTown:
		if lv < 30 {
			return "村"
		} else if lv < 100 {
			return "町"
		}
		return "都市"
	case hconst.LandForest:
		return "森"
	case hconst.LandFarm:
		return "農場"
	case hconst.LandFactory:
		return "工場"
	case hconst.LandBase:
		return "ミサイル基地"
	case hconst.LandDefence:
		return "防衛施設"
	case hconst.LandMountain:
		return "山"
	case hconst.LandMonster:
		_, name, _ := monsterSpec(lv)
		return name
	case hconst.LandSbase:
		return "海底基地"
	case hconst.LandOil:
		return "海底油田"
	case hconst.LandMonument:
		return hconst.MonumentName[lv]
	case hconst.LandHaribote:
		return "ハリボテ"
	default:
		return "不明"
	}
}

// monsterSpec returns monster specifications from landValue
// Returns: kind, name, hp
// Ref: perl/lib/Hako/Main.pm:958
func monsterSpec(lv int) (int, string, int) {
	// 種類
	kind := lv / 10

	// 名前
	name := hconst.MonsterName[kind]

	// 体力
	hp := lv - (kind * 10)

	return kind, name, hp
}

// wideDamage applies wide area damage (meteor/missile impact)
// Ref: perl/lib/Hako/Turn.pm:2444
func wideDamage(id, name string, land [][]int, landValue [][]int, x, y int) {
	for i := 0; i < 19; i++ {
		sx := x + hconst.Ax[i]
		sy := y + hconst.Ay[i]

		// 行による位置調整
		if (sy%2) == 0 && (y%2) == 1 {
			sx--
		}

		// 範囲外判定
		if sx < 0 || sx >= hconst.IslandSize || sy < 0 || sy >= hconst.IslandSize {
			continue
		}

		landKind := land[sx][sy]
		lv := landValue[sx][sy]
		landName := getLandName(landKind, lv)
		point := fmt.Sprintf("(%d, %d)", sx, sy)

		// 範囲による分岐
		if i < 7 {
			// 中心、および1ヘックス
			if landKind == hconst.LandSea {
				landValue[sx][sy] = 0
				continue
			} else if landKind == hconst.LandSbase || landKind == hconst.LandOil {
				logWideDamageSea2(id, name, landName, point)
				land[sx][sy] = hconst.LandSea
				landValue[sx][sy] = 0
			} else {
				if landKind == hconst.LandMonster {
					logWideDamageMonsterSea(id, name, landName, point)
				} else {
					logWideDamageSea(id, name, landName, point)
				}
				land[sx][sy] = hconst.LandSea
				if i == 0 {
					// 海
					landValue[sx][sy] = 0
				} else {
					// 浅瀬
					landValue[sx][sy] = 1
				}
			}
		} else {
			// 2ヘックス
			if landKind == hconst.LandSea || landKind == hconst.LandOil ||
				landKind == hconst.LandWaste || landKind == hconst.LandMountain ||
				landKind == hconst.LandSbase {
				continue
			} else if landKind == hconst.LandMonster {
				logWideDamageMonster(id, name, landName, point)
				land[sx][sy] = hconst.LandWaste
				landValue[sx][sy] = 0
			} else {
				logWideDamageWaste(id, name, landName, point)
				land[sx][sy] = hconst.LandWaste
				landValue[sx][sy] = 0
			}
		}
	}
}
