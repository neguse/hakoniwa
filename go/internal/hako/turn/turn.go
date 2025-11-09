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
// Ref: perl/lib/Hako/Turn.pm:518
func income(island *types.Island) {
	// Food consumption
	if island.Pop > island.Food {
		// Starvation
		logStarve(island.ID, island.Name, island.Pop-island.Food)
		island.Pop -= (island.Pop - island.Food) * 2 / 3
		island.Food = 0
	} else {
		island.Food -= island.Pop
	}

	// Money income
	island.Money += island.Pop * 10

	// Food income (each farm produces 5 food units)
	island.Food += island.Farm * 5
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
		// Do nothing
		return 1

	case hconst.ComPrepare, hconst.ComPrepare2:
		// Land preparation
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

	default:
		// Other commands: Phase 1 simplified - skip
		return 1
	}
}

// doEachHex processes each hex on an island
// Ref: perl/lib/Hako/Turn.pm:1669
func doEachHex(island *types.Island) {
	// Phase 1: Simplified implementation
	// Process growth for forests and towns

	for i := 0; i < hconst.PointNumber; i++ {
		x := variable.Rpx[i]
		y := variable.Rpy[i]

		land := island.Land[x][y]
		lv := island.LandValue[x][y]

		switch land {
		case hconst.LandForest:
			// Forest growth
			if lv < 200 && rand.Intn(10) < 5 {
				island.LandValue[x][y]++
			}

		case hconst.LandTown:
			// Town growth
			if lv < 200 && rand.Intn(10) < 3 {
				island.LandValue[x][y]++
			}
		}
	}
}

// doIslandProcess processes whole island events
// Ref: perl/lib/Hako/Turn.pm:1956
func doIslandProcess(num int, island *types.Island) {
	// Phase 1: Simplified implementation
	// Full implementation would include:
	// - Monsters
	// - Earthquakes
	// - Typhoons
	// - Tsunamis
	// - Meteors
	// - Eruptions
	// etc.

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

func isValidCoord(x, y int) bool {
	return x >= 0 && x < hconst.IslandSize && y >= 0 && y < hconst.IslandSize
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
