// Package core provides file I/O functions for the Hakoniwa game.
// This file contains file read/write functions translated from Perl lib/Hako/Main.pm
//
// Ref: perl/lib/Hako/Main.pm:197-441
package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/types"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// ReadIslandsFile reads the main island data file
// Ref: perl/lib/Hako/Main.pm:197
func ReadIslandsFile(id string) bool {
	// Open data file
	filename := fmt.Sprintf("%s/hakojima.dat", hconst.DirName)
	file, err := os.Open(filename)
	if err != nil {
		// Try temp file
		tmpName := fmt.Sprintf("%s/hakojima.tmp", hconst.DirName)
		os.Rename(tmpName, filename)
		file, err = os.Open(filename)
		if err != nil {
			return false
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read turn number
	if !scanner.Scan() {
		return false
	}
	variable.IslandTurn, _ = strconv.Atoi(scanner.Text())
	if variable.IslandTurn == 0 {
		return false
	}

	// Read last time
	if !scanner.Scan() {
		return false
	}
	variable.IslandLastTime, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	if variable.IslandLastTime == 0 {
		return false
	}

	// Read island count
	if !scanner.Scan() {
		return false
	}
	variable.IslandNumber, _ = strconv.Atoi(scanner.Text())

	// Read next ID
	if !scanner.Scan() {
		return false
	}
	variable.IslandNextID, _ = strconv.Atoi(scanner.Text())

	// Check if turn processing needed
	now := time.Now().Unix()
	num := 0
	if (hconst.Debug && variable.MainMode == "Hdebugturn") ||
		(now-variable.IslandLastTime >= hconst.UnitTime) {
		variable.MainMode = "turn"
		num = -1 // Read all islands
	} else if id != "" && id != "0" {
		numID, _ := strconv.Atoi(id)
		num = numID
	}

	// Read islands
	variable.Islands = make([]*types.Island, variable.IslandNumber)
	variable.IDToNumber = make(map[string]int)
	variable.IDToName = make(map[string]string)

	for i := 0; i < variable.IslandNumber; i++ {
		island := readIsland(scanner, num)
		if island == nil {
			return false
		}
		variable.Islands[i] = island
		variable.IDToNumber[island.ID] = i
	}

	return true
}

// readIsland reads one island's data
// Ref: perl/lib/Hako/Main.pm:246
func readIsland(scanner *bufio.Scanner, num int) *types.Island {
	island := &types.Island{}

	// Name and score
	if !scanner.Scan() {
		return nil
	}
	line := scanner.Text()
	parts := strings.Split(line, ",")
	island.Name = parts[0]
	if len(parts) > 1 {
		// Score exists but not used in current version
	}

	// ID
	if !scanner.Scan() {
		return nil
	}
	island.ID = scanner.Text()

	// Prize
	if !scanner.Scan() {
		return nil
	}
	island.Prize, _ = strconv.Atoi(scanner.Text())

	// Absent
	if !scanner.Scan() {
		return nil
	}
	island.Absent, _ = strconv.Atoi(scanner.Text())

	// Comment
	if !scanner.Scan() {
		return nil
	}
	island.Comment = scanner.Text()

	// Password
	if !scanner.Scan() {
		return nil
	}
	island.Password = scanner.Text()

	// Money
	if !scanner.Scan() {
		return nil
	}
	island.Money, _ = strconv.Atoi(scanner.Text())

	// Food
	if !scanner.Scan() {
		return nil
	}
	island.Food, _ = strconv.Atoi(scanner.Text())

	// Pop
	if !scanner.Scan() {
		return nil
	}
	island.Pop, _ = strconv.Atoi(scanner.Text())

	// Area
	if !scanner.Scan() {
		return nil
	}
	island.Area, _ = strconv.Atoi(scanner.Text())

	// Farm
	if !scanner.Scan() {
		return nil
	}
	island.Farm, _ = strconv.Atoi(scanner.Text())

	// Factory
	if !scanner.Scan() {
		return nil
	}
	island.Factory, _ = strconv.Atoi(scanner.Text())

	// Mountain
	if !scanner.Scan() {
		return nil
	}
	island.Mountain, _ = strconv.Atoi(scanner.Text())

	// Save to IDToName
	variable.IDToName[island.ID] = island.Name

	// Read land data if needed
	idNum, _ := strconv.Atoi(island.ID)
	if num == -1 || num == idNum {
		if !readIslandData(island) {
			// If failed, create empty data
			size := hconst.IslandSize
			island.Land = make([][]int, size)
			island.LandValue = make([][]int, size)
			for i := 0; i < size; i++ {
				island.Land[i] = make([]int, size)
				island.LandValue[i] = make([]int, size)
			}
			island.Commands = make([]types.Command, 0)
			island.Lbbs = make([]types.LbbsEntry, 0)
		}
	}

	return island
}

// readIslandData reads island-specific data from island.{ID} file
// Ref: perl/lib/Hako/Main.pm:283
func readIslandData(island *types.Island) bool {
	filename := fmt.Sprintf("%s/island.%s", hconst.DirName, island.ID)
	file, err := os.Open(filename)
	if err != nil {
		// Try temp file
		tmpName := fmt.Sprintf("%s/islandtmp.%s", hconst.DirName, island.ID)
		os.Rename(tmpName, filename)
		file, err = os.Open(filename)
		if err != nil {
			return false
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	size := hconst.IslandSize

	// Initialize land arrays
	island.Land = make([][]int, size)
	island.LandValue = make([][]int, size)
	for i := 0; i < size; i++ {
		island.Land[i] = make([]int, size)
		island.LandValue[i] = make([]int, size)
	}

	// Read land data (hex format: 1 char land type + 2 chars value)
	for y := 0; y < size; y++ {
		if !scanner.Scan() {
			return false
		}
		line := scanner.Text()
		for x := 0; x < size; x++ {
			if len(line) < 3 {
				return false
			}
			landType, _ := strconv.ParseInt(line[0:1], 16, 32)
			landValue, _ := strconv.ParseInt(line[1:3], 16, 32)
			island.Land[x][y] = int(landType)
			island.LandValue[x][y] = int(landValue)
			line = line[3:]
		}
	}

	// Read commands
	island.Commands = make([]types.Command, 0, hconst.CommandMax)
	for i := 0; i < hconst.CommandMax; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) >= 5 {
			cmd := types.Command{}
			cmd.Kind, _ = strconv.Atoi(parts[0])
			cmd.Target = parts[1]
			cmd.X, _ = strconv.Atoi(parts[2])
			cmd.Y, _ = strconv.Atoi(parts[3])
			cmd.Arg, _ = strconv.Atoi(parts[4])
			island.Commands = append(island.Commands, cmd)
		}
	}

	// Read local BBS
	island.Lbbs = make([]types.LbbsEntry, 0, hconst.LbbsMax)
	for i := 0; i < hconst.LbbsMax; i++ {
		if !scanner.Scan() {
			break
		}
		entry := types.LbbsEntry{
			Message: scanner.Text(),
		}
		island.Lbbs = append(island.Lbbs, entry)
	}

	return true
}

// WriteIslandsFile writes the main island data file
// Ref: perl/lib/Hako/Main.pm:348
func WriteIslandsFile() bool {
	// Write to temp file first
	filename := fmt.Sprintf("%s/hakojima.tmp", hconst.DirName)
	file, err := os.Create(filename)
	if err != nil {
		return false
	}

	writer := bufio.NewWriter(file)

	// Write header
	fmt.Fprintf(writer, "%d\n", variable.IslandTurn)
	fmt.Fprintf(writer, "%d\n", variable.IslandLastTime)
	fmt.Fprintf(writer, "%d\n", variable.IslandNumber)
	fmt.Fprintf(writer, "%d\n", variable.IslandNextID)

	// Write islands
	for i := 0; i < variable.IslandNumber; i++ {
		island := variable.Islands[i]
		fmt.Fprintf(writer, "%s\n", island.Name)
		fmt.Fprintf(writer, "%s\n", island.ID)
		fmt.Fprintf(writer, "%d\n", island.Prize)
		fmt.Fprintf(writer, "%d\n", island.Absent)
		fmt.Fprintf(writer, "%s\n", island.Comment)
		fmt.Fprintf(writer, "%s\n", island.Password)
		fmt.Fprintf(writer, "%d\n", island.Money)
		fmt.Fprintf(writer, "%d\n", island.Food)
		fmt.Fprintf(writer, "%d\n", island.Pop)
		fmt.Fprintf(writer, "%d\n", island.Area)
		fmt.Fprintf(writer, "%d\n", island.Farm)
		fmt.Fprintf(writer, "%d\n", island.Factory)
		fmt.Fprintf(writer, "%d\n", island.Mountain)
	}

	writer.Flush()
	file.Close()

	// Rename to actual file
	actualName := fmt.Sprintf("%s/hakojima.dat", hconst.DirName)
	return os.Rename(filename, actualName) == nil
}

// WriteIsland writes individual island data
// Ref: perl/lib/Hako/Main.pm:376
func WriteIsland(island *types.Island) bool {
	// Write to temp file first
	filename := fmt.Sprintf("%s/islandtmp.%s", hconst.DirName, island.ID)
	file, err := os.Create(filename)
	if err != nil {
		return false
	}

	writer := bufio.NewWriter(file)
	size := hconst.IslandSize

	// Write land data in hex format
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			fmt.Fprintf(writer, "%X%02X", island.Land[x][y], island.LandValue[x][y])
		}
		fmt.Fprintf(writer, "\n")
	}

	// Write commands
	for i := 0; i < hconst.CommandMax; i++ {
		if i < len(island.Commands) {
			cmd := island.Commands[i]
			fmt.Fprintf(writer, "%d,%s,%d,%d,%d\n",
				cmd.Kind, cmd.Target, cmd.X, cmd.Y, cmd.Arg)
		} else {
			fmt.Fprintf(writer, "0,0,0,0,0\n")
		}
	}

	// Write local BBS
	for i := 0; i < hconst.LbbsMax; i++ {
		if i < len(island.Lbbs) {
			fmt.Fprintf(writer, "%s\n", island.Lbbs[i].Message)
		} else {
			fmt.Fprintf(writer, "\n")
		}
	}

	writer.Flush()
	file.Close()

	// Rename to actual file
	actualName := fmt.Sprintf("%s/island.%s", hconst.DirName, island.ID)
	return os.Rename(filename, actualName) == nil
}

// ReadIsland reads individual island data (exported for external use)
// Ref: perl/lib/Hako/Main.pm:246
func ReadIsland(island *types.Island) bool {
	return readIslandData(island)
}
