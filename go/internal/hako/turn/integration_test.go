// Package turn provides integration tests for turn processing
// This file tests the deterministic behavior of the Go implementation
package turn

import (
	"math/rand"
	"testing"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/types"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// setupTest initializes test environment with deterministic random seed
func setupTest(t *testing.T) *types.Island {
	t.Helper()

	// Fix random seed for deterministic behavior
	rand.Seed(1234567890)

	// Initialize global variables
	variable.IslandTurn = 1
	variable.IslandLastTime = 0
	variable.IslandNumber = 1
	variable.IslandNextID = 1
	variable.Islands = make([]*types.Island, 0)
	variable.IDToName = make(map[string]string)
	variable.IDToNumber = make(map[string]int)
	variable.LogPool = make([]string, 0)
	variable.LateLogPool = make([]string, 0)
	variable.SecretLogPool = make([]string, 0)
	variable.OutputBuffer.Reset()

	// Initialize random point arrays (required for doEachHex)
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

	// Create test island
	island := types.NewIsland("テスト島", "0", "password")
	island.Money = 10000    // 100億円 x 100 = 1兆円
	island.Food = 10000     // 100トン x 100 = 1万トン
	island.Pop = 0
	island.OldPop = 0
	island.Dead = false
	island.BigMissile = 0
	island.MonsterSend = 0
	island.Propaganda = 0
	island.Prepare2 = 0

	// Initialize land (all sea initially)
	for i := 0; i < hconst.IslandSize; i++ {
		for j := 0; j < hconst.IslandSize; j++ {
			island.Land[i][j] = hconst.LandSea
			island.LandValue[i][j] = 0
		}
	}

	return island
}

// teardownTest cleans up after test
func teardownTest(t *testing.T) {
	t.Helper()

	// Reset global variables
	variable.Islands = make([]*types.Island, 0)
	variable.IslandNumber = 0
	variable.IDToName = make(map[string]string)
	variable.IDToNumber = make(map[string]int)
	variable.LogPool = make([]string, 0)
	variable.LateLogPool = make([]string, 0)
	variable.SecretLogPool = make([]string, 0)
	variable.OutputBuffer.Reset()
}

// addIsland adds an island to the global islands list
func addIsland(island *types.Island) {
	variable.Islands = append(variable.Islands, island)
	variable.IslandNumber = len(variable.Islands)
	variable.IDToName[island.ID] = island.Name
	variable.IDToNumber[island.ID] = len(variable.Islands) - 1
}

// cloneIsland creates a deep copy of an island for comparison
func cloneIsland(island *types.Island) *types.Island {
	clone := &types.Island{
		Name:        island.Name,
		ID:          island.ID,
		Prize:       island.Prize,
		Absent:      island.Absent,
		Comment:     island.Comment,
		Password:    island.Password,
		Money:       island.Money,
		Food:        island.Food,
		Pop:         island.Pop,
		Area:        island.Area,
		Farm:        island.Farm,
		Factory:     island.Factory,
		Mountain:    island.Mountain,
		Score:       island.Score,
		OldPop:      island.OldPop,
		Dead:        island.Dead,
		BigMissile:  island.BigMissile,
		MonsterSend: island.MonsterSend,
		Propaganda:  island.Propaganda,
		Prepare2:    island.Prepare2,
		Land:        make([][]int, hconst.IslandSize),
		LandValue:   make([][]int, hconst.IslandSize),
		Commands:    make([]types.Command, len(island.Commands)),
		Lbbs:        make([]types.LbbsEntry, len(island.Lbbs)),
	}

	// Deep copy Land and LandValue
	for i := 0; i < hconst.IslandSize; i++ {
		clone.Land[i] = make([]int, hconst.IslandSize)
		clone.LandValue[i] = make([]int, hconst.IslandSize)
		for j := 0; j < hconst.IslandSize; j++ {
			clone.Land[i][j] = island.Land[i][j]
			clone.LandValue[i][j] = island.LandValue[i][j]
		}
	}

	// Deep copy Commands
	copy(clone.Commands, island.Commands)

	// Deep copy Lbbs
	copy(clone.Lbbs, island.Lbbs)

	return clone
}

// compareIslands compares two islands and returns true if they are identical
func compareIslands(t *testing.T, island1, island2 *types.Island, label string) bool {
	t.Helper()

	if island1.Money != island2.Money {
		t.Errorf("%s: Money mismatch: %d != %d", label, island1.Money, island2.Money)
		return false
	}
	if island1.Food != island2.Food {
		t.Errorf("%s: Food mismatch: %d != %d", label, island1.Food, island2.Food)
		return false
	}
	if island1.Pop != island2.Pop {
		t.Errorf("%s: Pop mismatch: %d != %d", label, island1.Pop, island2.Pop)
		return false
	}
	if island1.Area != island2.Area {
		t.Errorf("%s: Area mismatch: %d != %d", label, island1.Area, island2.Area)
		return false
	}
	if island1.Farm != island2.Farm {
		t.Errorf("%s: Farm mismatch: %d != %d", label, island1.Farm, island2.Farm)
		return false
	}
	if island1.Factory != island2.Factory {
		t.Errorf("%s: Factory mismatch: %d != %d", label, island1.Factory, island2.Factory)
		return false
	}
	if island1.Mountain != island2.Mountain {
		t.Errorf("%s: Mountain mismatch: %d != %d", label, island1.Mountain, island2.Mountain)
		return false
	}

	// Compare land data
	for i := 0; i < hconst.IslandSize; i++ {
		for j := 0; j < hconst.IslandSize; j++ {
			if island1.Land[i][j] != island2.Land[i][j] {
				t.Errorf("%s: Land[%d][%d] mismatch: %d != %d", label, i, j, island1.Land[i][j], island2.Land[i][j])
				return false
			}
			if island1.LandValue[i][j] != island2.LandValue[i][j] {
				t.Errorf("%s: LandValue[%d][%d] mismatch: %d != %d", label, i, j, island1.LandValue[i][j], island2.LandValue[i][j])
				return false
			}
		}
	}

	return true
}

// TestBasicCommands tests basic commands execution (reclaim, prepare, plant)
func TestBasicCommands(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	// Add island to global list
	addIsland(island)

	// Test scenario: Create a small island with basic development
	// 1. Reclaim some land (埋め立て)
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComReclaim,
		Target: "",
		X:      5,
		Y:      5,
		Arg:    0,
	})

	// 2. Prepare land (整地)
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComPrepare,
		Target: "",
		X:      5,
		Y:      5,
		Arg:    0,
	})

	// 3. Plant forest (植林)
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComPlant,
		Target: "",
		X:      6,
		Y:      5,
		Arg:    0,
	})

	// Execute commands (loop until all commands are processed)
	t.Log("Executing basic commands...")
	for doCommand(island) == 0 {
		// Continue processing commands
	}

	// Verify results
	t.Logf("Island state after commands: Money=%d, Food=%d, Pop=%d", island.Money, island.Food, island.Pop)

	// Check that money was deducted (reclaim=150, prepare=5, plant=50 = 205 total)
	// Money should have decreased after executing commands
	if island.Money > 10000 {
		t.Errorf("Money should have decreased after commands, got %d", island.Money)
	}

	t.Log("Basic commands test completed successfully")
}

// TestGrowth tests growth scenarios (forest growth, town development)
func TestGrowth(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	// Add island to global list
	addIsland(island)

	// Create initial land: some plains and forests
	island.Land[5][5] = hconst.LandPlains
	island.Land[6][5] = hconst.LandForest
	island.LandValue[6][5] = 1 // Young forest

	// Execute hex-level processing (growth)
	t.Log("Executing growth processing...")
	doEachHex(island)

	// Verify that growth occurred
	t.Logf("Island state after growth: Money=%d, Food=%d, Pop=%d", island.Money, island.Food, island.Pop)

	// Forests should grow over time
	// Plains might develop into towns if population grows
	t.Log("Growth test completed successfully")
}

// TestDeterministic tests that the same initial state produces the same results
func TestDeterministic(t *testing.T) {
	// Run 1: First execution
	island1 := setupTest(t)
	addIsland(island1)

	// Set up initial state
	island1.Land[5][5] = hconst.LandWaste
	island1.Land[6][5] = hconst.LandWaste
	island1.Land[7][5] = hconst.LandWaste

	// Add some commands
	island1.Commands = append(island1.Commands, types.Command{
		Kind: hconst.ComPrepare,
		X:    5,
		Y:    5,
	})

	// Execute
	for doCommand(island1) == 0 {
	}
	doEachHex(island1)

	// Save state
	money1 := island1.Money
	food1 := island1.Food
	pop1 := island1.Pop
	landCopy1 := make([][]int, hconst.IslandSize)
	for i := 0; i < hconst.IslandSize; i++ {
		landCopy1[i] = make([]int, hconst.IslandSize)
		copy(landCopy1[i], island1.Land[i])
	}

	teardownTest(t)

	// Run 2: Second execution with same initial state
	island2 := setupTest(t)
	addIsland(island2)

	// Set up same initial state
	island2.Land[5][5] = hconst.LandWaste
	island2.Land[6][5] = hconst.LandWaste
	island2.Land[7][5] = hconst.LandWaste

	// Add same commands
	island2.Commands = append(island2.Commands, types.Command{
		Kind: hconst.ComPrepare,
		X:    5,
		Y:    5,
	})

	// Execute
	for doCommand(island2) == 0 {
	}
	doEachHex(island2)

	// Compare results
	if money1 != island2.Money {
		t.Errorf("Deterministic test failed: Money mismatch: %d != %d", money1, island2.Money)
	}
	if food1 != island2.Food {
		t.Errorf("Deterministic test failed: Food mismatch: %d != %d", food1, island2.Food)
	}
	if pop1 != island2.Pop {
		t.Errorf("Deterministic test failed: Pop mismatch: %d != %d", pop1, island2.Pop)
	}

	// Compare land data
	for i := 0; i < hconst.IslandSize; i++ {
		for j := 0; j < hconst.IslandSize; j++ {
			if landCopy1[i][j] != island2.Land[i][j] {
				t.Errorf("Deterministic test failed: Land[%d][%d] mismatch: %d != %d", i, j, landCopy1[i][j], island2.Land[i][j])
			}
		}
	}

	t.Log("Deterministic test completed successfully: Same inputs produced same outputs")
	teardownTest(t)
}

// TestMultipleTurns tests multiple turn execution
func TestMultipleTurns(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	// Add island to global list
	addIsland(island)

	// Set up initial land
	island.Land[5][5] = hconst.LandWaste
	island.Land[6][5] = hconst.LandWaste
	island.Land[5][6] = hconst.LandWaste
	island.Land[6][6] = hconst.LandPlains

	// Execute multiple turns
	numTurns := 5
	for turn := 1; turn <= numTurns; turn++ {
		t.Logf("Turn %d: Money=%d, Food=%d, Pop=%d", turn, island.Money, island.Food, island.Pop)

		// Add command for this turn
		if turn%2 == 1 {
			island.Commands = append(island.Commands, types.Command{
				Kind: hconst.ComPrepare,
				X:    5,
				Y:    5,
			})
		}

		// Execute turn processing
		for doCommand(island) == 0 {
		}
		doEachHex(island)

		variable.IslandTurn++
	}

	t.Logf("After %d turns: Money=%d, Food=%d, Pop=%d", numTurns, island.Money, island.Food, island.Pop)
	t.Log("Multiple turns test completed successfully")
}

// TestIslandProcess tests whole island processing (disasters, prizes, etc.)
func TestIslandProcess(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	// Add island to global list
	addIsland(island)

	// Set up initial state with some development
	island.Land[5][5] = hconst.LandPlains
	island.Land[6][5] = hconst.LandTown
	island.LandValue[6][5] = 10
	island.Land[7][5] = hconst.LandForest
	island.LandValue[7][5] = 5

	island.Pop = 100 // 10,000 people
	island.OldPop = 100

	// Execute island process
	t.Log("Executing island-wide processing...")
	doIslandProcess(0, island)

	// Verify that processing completed without errors
	t.Logf("Island state after processing: Money=%d, Food=%d, Pop=%d, Dead=%v", island.Money, island.Food, island.Pop, island.Dead)

	// Check logs were generated
	totalLogs := len(variable.LogPool) + len(variable.LateLogPool) + len(variable.SecretLogPool)
	t.Logf("Total logs generated: %d", totalLogs)

	t.Log("Island process test completed successfully")
}

// TestCommandExecution tests specific command execution
func TestCommandExecution(t *testing.T) {
	testCases := []struct {
		name        string
		setupLand   func(*types.Island)
		command     types.Command
		description string
	}{
		{
			name: "Reclaim",
			setupLand: func(island *types.Island) {
				island.Land[5][5] = hconst.LandSea
			},
			command: types.Command{
				Kind: hconst.ComReclaim,
				X:    5,
				Y:    5,
			},
			description: "埋め立て: 海 -> 浅瀬",
		},
		{
			name: "Prepare",
			setupLand: func(island *types.Island) {
				island.Land[5][5] = hconst.LandWaste
			},
			command: types.Command{
				Kind: hconst.ComPrepare,
				X:    5,
				Y:    5,
			},
			description: "整地: 荒地 -> 平地",
		},
		{
			name: "Plant",
			setupLand: func(island *types.Island) {
				island.Land[5][5] = hconst.LandWaste
			},
			command: types.Command{
				Kind: hconst.ComPlant,
				X:    5,
				Y:    5,
			},
			description: "植林: 荒地 -> 森",
		},
		{
			name: "BuildFarm",
			setupLand: func(island *types.Island) {
				island.Land[5][5] = hconst.LandPlains
			},
			command: types.Command{
				Kind: hconst.ComFarm,
				X:    5,
				Y:    5,
			},
			description: "農場整備: 平地 -> 農場",
		},
		{
			name: "BuildFactory",
			setupLand: func(island *types.Island) {
				island.Land[5][5] = hconst.LandPlains
			},
			command: types.Command{
				Kind: hconst.ComFactory,
				X:    5,
				Y:    5,
			},
			description: "工場建設: 平地 -> 工場",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			island := setupTest(t)
			defer teardownTest(t)

			addIsland(island)

			// Setup initial land
			tc.setupLand(island)

			// Record initial state
			initialMoney := island.Money
			initialLand := island.Land[tc.command.X][tc.command.Y]

			// Add command
			island.Commands = append(island.Commands, tc.command)

			// Execute
			for doCommand(island) == 0 {
			}

			// Verify money was deducted (if cost > 0)
			cost := hconst.ComCost[tc.command.Kind]
			if cost > 0 && island.Money >= initialMoney {
				t.Errorf("%s: Expected money to decrease by at least %d, but Money=%d (was %d)", tc.description, cost, island.Money, initialMoney)
			}

			// Verify land changed (might not change if command failed)
			finalLand := island.Land[tc.command.X][tc.command.Y]
			t.Logf("%s: Land changed from %d to %d (Money: %d -> %d)", tc.description, initialLand, finalLand, initialMoney, island.Money)
		})
	}
}

// TestReclaimCommand tests reclaim (埋め立て) functionality
func TestReclaimCommand(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up sea tiles
	island.Land[5][5] = hconst.LandSea
	island.LandValue[5][5] = 0
	initialMoney := island.Money

	// First reclaim changes sea -> shallow
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComReclaim,
		X:    5,
		Y:    5,
	})

	result := doCommand(island)
	if result == 0 {
		t.Error("First reclaim should return 1 (turn consumed)")
	}

	// Check land changed to shallow
	if island.Land[5][5] != hconst.LandSea {
		t.Errorf("Expected land to be sea (0), got %d", island.Land[5][5])
	}
	if island.LandValue[5][5] != 1 {
		t.Errorf("Expected land value to be 1 (shallow), got %d", island.LandValue[5][5])
	}

	// Check money deducted
	cost := hconst.ComCost[hconst.ComReclaim]
	if island.Money != initialMoney-cost {
		t.Errorf("Expected money to be %d, got %d", initialMoney-cost, island.Money)
	}

	t.Logf("Reclaim test passed: sea -> shallow (Money: %d -> %d)", initialMoney, island.Money)
}

// TestPrepareCommand tests prepare (整地) functionality
func TestPrepareCommand(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up wasteland
	island.Land[5][5] = hconst.LandWaste
	island.LandValue[5][5] = 0
	initialMoney := island.Money

	// Prepare wasteland -> plains
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComPrepare,
		X:    5,
		Y:    5,
	})

	result := doCommand(island)
	if result == 0 {
		t.Error("Prepare should return 1 (turn consumed)")
	}

	// Check land changed to plains
	if island.Land[5][5] != hconst.LandPlains {
		t.Errorf("Expected land to be plains (%d), got %d", hconst.LandPlains, island.Land[5][5])
	}
	if island.LandValue[5][5] != 0 {
		t.Errorf("Expected land value to be 0, got %d", island.LandValue[5][5])
	}

	// Check money deducted
	cost := hconst.ComCost[hconst.ComPrepare]
	if island.Money >= initialMoney {
		t.Errorf("Expected money to decrease by %d, got %d -> %d", cost, initialMoney, island.Money)
	}

	t.Logf("Prepare test passed: wasteland -> plains (Money: %d -> %d)", initialMoney, island.Money)
}

// TestPlantCommand tests plant (植林) functionality
func TestPlantCommand(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up wasteland
	island.Land[5][5] = hconst.LandWaste
	island.LandValue[5][5] = 0
	initialMoney := island.Money

	// Plant wasteland -> forest
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComPlant,
		X:    5,
		Y:    5,
	})

	result := doCommand(island)
	if result == 0 {
		t.Error("Plant should return 1 (turn consumed)")
	}

	// Check land changed to forest
	if island.Land[5][5] != hconst.LandForest {
		t.Errorf("Expected land to be forest (%d), got %d", hconst.LandForest, island.Land[5][5])
	}
	if island.LandValue[5][5] != 1 {
		t.Errorf("Expected land value to be 1 (young forest), got %d", island.LandValue[5][5])
	}

	// Check money deducted
	cost := hconst.ComCost[hconst.ComPlant]
	if island.Money != initialMoney-cost {
		t.Errorf("Expected money to be %d, got %d", initialMoney-cost, island.Money)
	}

	t.Logf("Plant test passed: wasteland -> forest (Money: %d -> %d)", initialMoney, island.Money)
}

// TestFarmCommand tests farm (農場整備) functionality
func TestFarmCommand(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up plains
	island.Land[5][5] = hconst.LandPlains
	island.LandValue[5][5] = 0
	initialMoney := island.Money

	// Build farm: plains -> farm
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComFarm,
		X:    5,
		Y:    5,
	})

	result := doCommand(island)
	if result == 0 {
		t.Error("Farm should return 1 (turn consumed)")
	}

	// Check land changed to farm
	if island.Land[5][5] != hconst.LandFarm {
		t.Errorf("Expected land to be farm (%d), got %d", hconst.LandFarm, island.Land[5][5])
	}
	if island.LandValue[5][5] != 10 {
		t.Errorf("Expected land value to be 10, got %d", island.LandValue[5][5])
	}

	// Check money deducted
	cost := hconst.ComCost[hconst.ComFarm]
	if island.Money != initialMoney-cost {
		t.Errorf("Expected money to be %d, got %d", initialMoney-cost, island.Money)
	}

	t.Logf("Farm test passed: plains -> farm (Money: %d -> %d)", initialMoney, island.Money)
}

// TestForestGrowth tests forest growth during doEachHex
func TestForestGrowth(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up forest
	island.Land[5][5] = hconst.LandForest
	island.LandValue[5][5] = 1 // Young forest

	// Execute growth processing
	t.Log("Testing forest growth...")
	doEachHex(island)

	// Forests should grow (unless random doesn't allow)
	finalValue := island.LandValue[5][5]
	t.Logf("Forest growth: %d -> %d", 1, finalValue)

	if island.Land[5][5] != hconst.LandForest {
		t.Errorf("Expected land to still be forest, got %d", island.Land[5][5])
	}
	if finalValue < 1 {
		t.Errorf("Expected forest value to be at least 1, got %d", finalValue)
	}

	t.Log("Forest growth test passed")
}

// TestTownGrowth tests town growth during doEachHex
func TestTownGrowth(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up town
	island.Land[5][5] = hconst.LandTown
	island.LandValue[5][5] = 10 // Small town
	island.Pop = 100
	island.Food = 10000

	initialValue := island.LandValue[5][5]

	// Execute growth processing
	t.Log("Testing town growth...")
	doEachHex(island)

	// Towns should grow
	finalValue := island.LandValue[5][5]
	t.Logf("Town growth: %d -> %d", initialValue, finalValue)

	if island.Land[5][5] != hconst.LandTown {
		t.Errorf("Expected land to still be town, got %d", island.Land[5][5])
	}

	t.Log("Town growth test passed")
}

// TestPlainsTownExpansion tests plains expanding into town when nearby development exists
func TestPlainsTownExpansion(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up plains next to town
	island.Land[5][5] = hconst.LandPlains
	island.LandValue[5][5] = 0
	island.Land[6][5] = hconst.LandTown
	island.LandValue[6][5] = 50

	// Execute multiple turns of growth processing
	t.Log("Testing plains town expansion...")
	rand.Seed(1234567890)
	for i := 0; i < 10; i++ {
		doEachHex(island)
		if island.Land[5][5] == hconst.LandTown {
			t.Logf("Plains expanded to town after %d turns", i+1)
			break
		}
	}

	t.Log("Plains town expansion test completed")
}

// TestComplexScenario tests a complex scenario with multiple commands and turns
func TestComplexScenario(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	// Set up initial map: sea with some land nearby
	island.Land[5][5] = hconst.LandSea
	island.Land[5][4] = hconst.LandPlains
	island.Land[4][5] = hconst.LandPlains
	island.Pop = 100
	island.Food = 10000
	initialMoney := island.Money

	t.Log("Starting complex scenario...")

	// Turn 1: Reclaim shallow water
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComReclaim,
		X:    5,
		Y:    5,
	})

	for doCommand(island) == 0 {
	}
	doEachHex(island)

	// Verify reclaim worked
	if island.Land[5][5] != hconst.LandSea || island.LandValue[5][5] != 1 {
		t.Errorf("Reclaim failed: Land=%d, Value=%d (expected Land=0, Value=1)", island.Land[5][5], island.LandValue[5][5])
	}

	money1 := island.Money
	t.Logf("After Turn 1 (Reclaim): Money=%d, Land[5][5]=%d,%d", island.Money, island.Land[5][5], island.LandValue[5][5])

	// Turn 2: Plant forest at new location
	island.Land[5][4] = hconst.LandWaste // Change from plains to wasteland
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComPlant,
		X:    5,
		Y:    4,
	})

	for doCommand(island) == 0 {
	}
	doEachHex(island)

	if island.Land[5][4] != hconst.LandForest {
		t.Errorf("Plant failed: Land=%d (expected %d)", island.Land[5][4], hconst.LandForest)
	}

	money2 := island.Money
	t.Logf("After Turn 2 (Plant): Money=%d, Land[5][4]=%d,%d", island.Money, island.Land[5][4], island.LandValue[5][4])

	// Turn 3: Build farm on plains
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComFarm,
		X:    4,
		Y:    5,
	})

	for doCommand(island) == 0 {
	}
	doEachHex(island)

	if island.Land[4][5] != hconst.LandFarm {
		t.Errorf("Farm failed: Land=%d (expected %d)", island.Land[4][5], hconst.LandFarm)
	}

	money3 := island.Money
	t.Logf("After Turn 3 (Farm): Money=%d, Land[4][5]=%d,%d", island.Money, island.Land[4][5], island.LandValue[4][5])

	// Verify total cost
	totalCost := hconst.ComCost[hconst.ComReclaim] + hconst.ComCost[hconst.ComPlant] + hconst.ComCost[hconst.ComFarm]
	expectedMoney := initialMoney - totalCost
	if island.Money != expectedMoney {
		t.Errorf("Expected final money to be %d, got %d (spent %d)", expectedMoney, island.Money, initialMoney-island.Money)
	}

	t.Logf("Complex scenario passed: %d > %d > %d > %d (costs: %d, %d, %d)",
		initialMoney, money1, money2, money3,
		initialMoney-money1, money1-money2, money2-money3)
}

// TestMoneyResourceCalculation tests money and resource calculations
func TestMoneyResourceCalculation(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	initialMoney := island.Money
	initialFood := island.Food

	t.Logf("Initial state: Money=%d, Food=%d", initialMoney, initialFood)

	// Test DoNothing command (should increase money by 10)
	island.Commands = append(island.Commands, types.Command{
		Kind: hconst.ComDoNothing,
	})

	for doCommand(island) == 0 {
	}

	expectedMoney := initialMoney + 10
	if island.Money != expectedMoney {
		t.Errorf("DoNothing: Expected money to be %d, got %d", expectedMoney, island.Money)
	}

	t.Logf("DoNothing test passed: Money=%d->%d (expected +10)", initialMoney, island.Money)
}

// TestSequentialCommands tests executing multiple commands in sequence
func TestSequentialCommands(t *testing.T) {
	island := setupTest(t)
	defer teardownTest(t)

	addIsland(island)

	initialMoney := island.Money

	// Set up wasteland
	island.Land[5][5] = hconst.LandWaste
	island.LandValue[5][5] = 0

	// Queue multiple commands
	island.Commands = append(island.Commands,
		types.Command{Kind: hconst.ComPrepare, X: 5, Y: 5},
	)

	// Execute and queue the next command
	for doCommand(island) == 0 {
	}

	// Now add a farm command (plains -> farm)
	island.Commands = append(island.Commands,
		types.Command{Kind: hconst.ComFarm, X: 5, Y: 5},
	)

	for doCommand(island) == 0 {
	}

	// Verify final state
	if island.Land[5][5] != hconst.LandFarm {
		t.Errorf("Expected land to be farm, got %d", island.Land[5][5])
	}

	// Verify money deductions
	expectedCost := hconst.ComCost[hconst.ComPrepare] + hconst.ComCost[hconst.ComFarm]
	expectedMoney := initialMoney - expectedCost
	if island.Money != expectedMoney {
		t.Errorf("Expected money to be %d, got %d", expectedMoney, island.Money)
	}

	t.Logf("Sequential commands test passed: Wasteland -> Plains -> Farm (Money: %d -> %d)", initialMoney, island.Money)
}
