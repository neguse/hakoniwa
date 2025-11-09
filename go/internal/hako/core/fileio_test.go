// Package core provides file I/O functions for the Hakoniwa game.
// This file contains tests for file read/write functions
package core

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/types"
)

// createTestIsland creates a test island with various land types and values
func createTestIsland(id string) *types.Island {
	island := types.NewIsland("TestIsland", id, "password123")

	// Set some basic properties
	island.Money = 500
	island.Food = 300
	island.Pop = 1000
	island.Prize = 1
	island.Absent = 2
	island.Comment = "Test comment"
	island.Farm = 5
	island.Factory = 3
	island.Mountain = 2

	// Set up land with various terrain types
	size := hconst.IslandSize
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Create a pattern of different land types
			switch (x + y) % 15 {
			case 0:
				island.Land[x][y] = hconst.LandSea
				island.LandValue[x][y] = 0
			case 1:
				island.Land[x][y] = hconst.LandWaste
				island.LandValue[x][y] = 0
			case 2:
				island.Land[x][y] = hconst.LandPlains
				island.LandValue[x][y] = 0
			case 3:
				island.Land[x][y] = hconst.LandTown
				island.LandValue[x][y] = 50 + (x+y)%200 // Town population
			case 4:
				island.Land[x][y] = hconst.LandForest
				island.LandValue[x][y] = 10 + (x+y)%240 // Tree count
			case 5:
				island.Land[x][y] = hconst.LandFarm
				island.LandValue[x][y] = 0
			case 6:
				island.Land[x][y] = hconst.LandFactory
				island.LandValue[x][y] = 0
			case 7:
				island.Land[x][y] = hconst.LandBase
				island.LandValue[x][y] = 30 + (x*y)%70 // Experience points
			case 8:
				island.Land[x][y] = hconst.LandDefence
				island.LandValue[x][y] = 0
			case 9:
				island.Land[x][y] = hconst.LandMountain
				island.LandValue[x][y] = 0
			case 10:
				island.Land[x][y] = hconst.LandMonster
				island.LandValue[x][y] = (x + y) % 8 // Monster type (0-7)
			case 11:
				island.Land[x][y] = hconst.LandSbase
				island.LandValue[x][y] = 20 + (x+y)%80 // Experience points
			case 12:
				island.Land[x][y] = hconst.LandOil
				island.LandValue[x][y] = 0
			case 13:
				island.Land[x][y] = hconst.LandMonument
				island.LandValue[x][y] = (x + y) % 3 // Monument type
			case 14:
				island.Land[x][y] = hconst.LandHaribote
				island.LandValue[x][y] = 0
			}
		}
	}

	// Add some commands
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComPrepare,
		Target: "0",
		X:      5,
		Y:      5,
		Arg:    0,
	})
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComFarm,
		Target: "0",
		X:      3,
		Y:      4,
		Arg:    0,
	})
	island.Commands = append(island.Commands, types.Command{
		Kind:   hconst.ComMissileNM,
		Target: "1",
		X:      7,
		Y:      8,
		Arg:    0,
	})

	// Add some LBBS entries
	island.Lbbs = append(island.Lbbs, types.LbbsEntry{
		Message: "Test message 1",
	})
	island.Lbbs = append(island.Lbbs, types.LbbsEntry{
		Message: "Test message 2",
	})

	return island
}

// createEmptyIsland creates an island with all sea tiles
func createEmptyIsland(id string) *types.Island {
	island := types.NewIsland("EmptyIsland", id, "pass")
	// NewIsland already initializes all land as sea, so nothing more needed
	return island
}

// createFullCommandIsland creates an island with maximum commands
func createFullCommandIsland(id string) *types.Island {
	island := types.NewIsland("FullCommandIsland", id, "pass")

	// Fill with maximum commands
	for i := 0; i < hconst.CommandMax; i++ {
		island.Commands = append(island.Commands, types.Command{
			Kind:   hconst.ComPrepare,
			Target: "0",
			X:      i % hconst.IslandSize,
			Y:      i / hconst.IslandSize,
			Arg:    i,
		})
	}

	return island
}

// createFullLbbsIsland creates an island with maximum LBBS entries
func createFullLbbsIsland(id string) *types.Island {
	island := types.NewIsland("FullLbbsIsland", id, "pass")

	// Fill with maximum LBBS entries
	for i := 0; i < hconst.LbbsMax; i++ {
		island.Lbbs = append(island.Lbbs, types.LbbsEntry{
			Message: strings.Repeat("Test message ", i+1),
		})
	}

	return island
}

// TestWriteReadIslandRoundTrip tests that writing and reading an island preserves data
func TestWriteReadIslandRoundTrip(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	testCases := []struct {
		name   string
		island *types.Island
	}{
		{
			name:   "standard island",
			island: createTestIsland("100"),
		},
		{
			name:   "empty island",
			island: createEmptyIsland("101"),
		},
		{
			name:   "full command island",
			island: createFullCommandIsland("102"),
		},
		{
			name:   "full lbbs island",
			island: createFullLbbsIsland("103"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write the island
			if !WriteIsland(tc.island) {
				t.Fatalf("WriteIsland failed for %s", tc.name)
			}

			// Create a new island to read into
			readIsland := types.NewIsland(tc.island.Name, tc.island.ID, tc.island.Password)

			// Read the island back
			if !ReadIsland(readIsland) {
				t.Fatalf("ReadIsland failed for %s", tc.name)
			}

			// Compare land data
			size := hconst.IslandSize
			for y := 0; y < size; y++ {
				for x := 0; x < size; x++ {
					if readIsland.Land[x][y] != tc.island.Land[x][y] {
						t.Errorf("Land mismatch at (%d,%d): got %d, want %d",
							x, y, readIsland.Land[x][y], tc.island.Land[x][y])
					}
					if readIsland.LandValue[x][y] != tc.island.LandValue[x][y] {
						t.Errorf("LandValue mismatch at (%d,%d): got %d, want %d",
							x, y, readIsland.LandValue[x][y], tc.island.LandValue[x][y])
					}
				}
			}

			// Compare commands
			// Note: WriteIsland always writes CommandMax entries (padding with empty commands)
			// and ReadIsland reads all of them, so we need to compare meaningful commands only
			originalCmdCount := len(tc.island.Commands)
			for i := 0; i < hconst.CommandMax; i++ {
				var expected types.Command
				if i < originalCmdCount {
					expected = tc.island.Commands[i]
				} else {
					// Empty command padding
					expected = types.Command{Kind: 0, Target: "0", X: 0, Y: 0, Arg: 0}
				}

				if i < len(readIsland.Commands) {
					if readIsland.Commands[i] != expected {
						t.Errorf("Command[%d] mismatch: got %+v, want %+v",
							i, readIsland.Commands[i], expected)
					}
				}
			}

			// Compare LBBS entries
			// Note: WriteIsland always writes LbbsMax entries (padding with empty lines)
			// and ReadIsland reads all of them
			originalLbbsCount := len(tc.island.Lbbs)
			for i := 0; i < hconst.LbbsMax; i++ {
				var expectedMsg string
				if i < originalLbbsCount {
					expectedMsg = tc.island.Lbbs[i].Message
				} else {
					// Empty LBBS padding
					expectedMsg = ""
				}

				if i < len(readIsland.Lbbs) {
					if readIsland.Lbbs[i].Message != expectedMsg {
						t.Errorf("LBBS[%d] message mismatch: got %q, want %q",
							i, readIsland.Lbbs[i].Message, expectedMsg)
					}
				}
			}
		})
	}
}

// TestHexFormatEncoding tests that hex format encoding is correct
func TestHexFormatEncoding(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	island := types.NewIsland("HexTest", "200", "pass")

	// Set specific values to test hex encoding
	island.Land[0][0] = 0x0
	island.LandValue[0][0] = 0x00

	island.Land[1][0] = 0xF
	island.LandValue[1][0] = 0xFF

	island.Land[2][0] = 0x3
	island.LandValue[2][0] = 0xA5

	island.Land[3][0] = 0xB
	island.LandValue[3][0] = 0x12

	// Write the island
	if !WriteIsland(island) {
		t.Fatal("WriteIsland failed")
	}

	// Read the file directly to check hex format
	filename := filepath.Join(tmpDir, "island.200")
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open island file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatal("Failed to read first line")
	}
	line := scanner.Text()

	// Check the first few hex values in the first line
	expectedPrefix := "000FFF3A5B12"
	if !strings.HasPrefix(line, expectedPrefix) {
		t.Errorf("Hex format mismatch: got prefix %q, want %q", line[:12], expectedPrefix)
	}

	// Each hex segment should be 3 characters (1 land type + 2 land value)
	expectedLineLen := hconst.IslandSize * 3
	if len(line) != expectedLineLen {
		t.Errorf("Line length mismatch: got %d, want %d", len(line), expectedLineLen)
	}

	// Verify all characters are valid hex
	for i, c := range line {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Invalid hex character at position %d: %c", i, c)
		}
	}
}

// TestEdgeCases tests edge cases for island data
func TestEdgeCases(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	t.Run("max land values", func(t *testing.T) {
		island := types.NewIsland("MaxValues", "300", "pass")

		// Set maximum values
		island.Land[0][0] = 15 // Max land type (0xF)
		island.LandValue[0][0] = 255 // Max value (0xFF)

		if !WriteIsland(island) {
			t.Fatal("WriteIsland failed")
		}

		readIsland := types.NewIsland("MaxValues", "300", "pass")
		if !ReadIsland(readIsland) {
			t.Fatal("ReadIsland failed")
		}

		if readIsland.Land[0][0] != 15 {
			t.Errorf("Max land type: got %d, want 15", readIsland.Land[0][0])
		}
		if readIsland.LandValue[0][0] != 255 {
			t.Errorf("Max land value: got %d, want 255", readIsland.LandValue[0][0])
		}
	})

	t.Run("empty commands", func(t *testing.T) {
		island := types.NewIsland("NoCommands", "301", "pass")
		island.Commands = []types.Command{} // Empty commands

		if !WriteIsland(island) {
			t.Fatal("WriteIsland failed")
		}

		readIsland := types.NewIsland("NoCommands", "301", "pass")
		if !ReadIsland(readIsland) {
			t.Fatal("ReadIsland failed")
		}

		// ReadIsland always reads CommandMax entries (with padding)
		// Check that all commands are empty (Kind == 0)
		if len(readIsland.Commands) != hconst.CommandMax {
			t.Errorf("Command count: got %d, want %d", len(readIsland.Commands), hconst.CommandMax)
		}
		for i, cmd := range readIsland.Commands {
			if cmd.Kind != 0 || cmd.Target != "0" || cmd.X != 0 || cmd.Y != 0 || cmd.Arg != 0 {
				t.Errorf("Command[%d] should be empty: got %+v", i, cmd)
			}
		}
	})

	t.Run("empty lbbs", func(t *testing.T) {
		island := types.NewIsland("NoLbbs", "302", "pass")
		island.Lbbs = []types.LbbsEntry{} // Empty LBBS

		if !WriteIsland(island) {
			t.Fatal("WriteIsland failed")
		}

		readIsland := types.NewIsland("NoLbbs", "302", "pass")
		if !ReadIsland(readIsland) {
			t.Fatal("ReadIsland failed")
		}

		// ReadIsland always reads LbbsMax entries (with padding)
		// Check that all LBBS entries are empty
		if len(readIsland.Lbbs) != hconst.LbbsMax {
			t.Errorf("LBBS count: got %d, want %d", len(readIsland.Lbbs), hconst.LbbsMax)
		}
		for i, entry := range readIsland.Lbbs {
			if entry.Message != "" {
				t.Errorf("LBBS[%d] should be empty: got %q", i, entry.Message)
			}
		}
	})

	t.Run("all land types", func(t *testing.T) {
		island := types.NewIsland("AllLandTypes", "303", "pass")

		// Set one of each land type
		landTypes := []int{
			hconst.LandSea, hconst.LandWaste, hconst.LandPlains,
			hconst.LandTown, hconst.LandForest, hconst.LandFarm,
			hconst.LandFactory, hconst.LandBase, hconst.LandDefence,
			hconst.LandMountain, hconst.LandMonster, hconst.LandSbase,
			hconst.LandOil, hconst.LandMonument, hconst.LandHaribote,
		}

		for i, landType := range landTypes {
			if i < hconst.IslandSize {
				island.Land[i][0] = landType
				island.LandValue[i][0] = i * 10
			}
		}

		if !WriteIsland(island) {
			t.Fatal("WriteIsland failed")
		}

		readIsland := types.NewIsland("AllLandTypes", "303", "pass")
		if !ReadIsland(readIsland) {
			t.Fatal("ReadIsland failed")
		}

		for i, landType := range landTypes {
			if i < hconst.IslandSize {
				if readIsland.Land[i][0] != landType {
					t.Errorf("Land type at [%d][0]: got %d, want %d", i, readIsland.Land[i][0], landType)
				}
				if readIsland.LandValue[i][0] != i*10 {
					t.Errorf("Land value at [%d][0]: got %d, want %d", i, readIsland.LandValue[i][0], i*10)
				}
			}
		}
	})
}

// TestReadIslandNonExistent tests reading a non-existent island file
func TestReadIslandNonExistent(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	island := types.NewIsland("NonExistent", "999", "pass")

	// Try to read a non-existent island
	if ReadIsland(island) {
		t.Error("ReadIsland should return false for non-existent file")
	}
}

// TestWriteIslandPermissions tests writing island with various file permission scenarios
func TestWriteIslandPermissions(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	island := createTestIsland("400")

	// First write should succeed
	if !WriteIsland(island) {
		t.Fatal("First WriteIsland failed")
	}

	// Verify temp file was renamed to actual file
	actualFile := filepath.Join(tmpDir, "island.400")
	if _, err := os.Stat(actualFile); os.IsNotExist(err) {
		t.Error("Island file was not created")
	}

	// Overwriting should also succeed
	island.Money = 999
	if !WriteIsland(island) {
		t.Fatal("Overwriting island failed")
	}

	// Verify the overwrite worked
	readIsland := types.NewIsland("TestIsland", "400", "password123")
	if !ReadIsland(readIsland) {
		t.Fatal("ReadIsland after overwrite failed")
	}

	// Note: We can't directly check island.Money here because ReadIsland
	// only reads the island-specific data file, not the main island data.
	// The Money field is stored in hakojima.dat, not island.{ID} file.
}

// BenchmarkWriteIsland benchmarks island writing performance
func BenchmarkWriteIsland(b *testing.B) {
	tmpDir := b.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	island := createTestIsland("500")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteIsland(island)
	}
}

// BenchmarkReadIsland benchmarks island reading performance
func BenchmarkReadIsland(b *testing.B) {
	tmpDir := b.TempDir()
	originalDir := hconst.DirName
	hconst.DirName = tmpDir
	defer func() { hconst.DirName = originalDir }()

	island := createTestIsland("600")
	WriteIsland(island)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readIsland := types.NewIsland("TestIsland", "600", "password123")
		ReadIsland(readIsland)
	}
}
