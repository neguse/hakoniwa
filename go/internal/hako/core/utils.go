// Package core provides core utilities for the Hakoniwa game.
// This file contains utility functions translated from Perl lib/Hako/Main.pm
//
// Ref: perl/lib/Hako/Main.pm
package core

import (
	"crypto/des"
	"fmt"
	"math/rand"
	"strings"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// min returns the minimum of two integers
// Ref: perl/lib/Hako/Main.pm:858
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// encode encodes a password
// Ref: perl/lib/Hako/Main.pm:863
func encode(password string) string {
	if hconst.CryptOn {
		// Perl's crypt() with salt "h2"
		// Phase 1: simplified implementation (will be improved in Phase 2)
		return cryptCompat(password, "h2")
	}
	return password
}

// cryptCompat implements a simplified DES-based crypt compatible with Perl's crypt()
// Phase 1: Basic implementation
func cryptCompat(password, salt string) string {
	// Phase 1: Simplified - just return a hash
	// TODO: Implement proper DES crypt for full Perl compatibility
	if len(password) == 0 {
		return ""
	}
	// Placeholder: In Phase 1, we'll use a simple encoding
	// This needs to be compatible with Perl's crypt()
	return password // FIXME: Implement proper crypt
}

// CheckPassword checks if a password matches (exported for external use)
// Ref: perl/lib/Hako/Main.pm:873
func CheckPassword(stored, input string) bool {
	// null check
	if input == "" {
		return false
	}

	// Master password check
	if hconst.MasterPassword == input {
		return true
	}

	// Normal check
	if stored == encode(input) {
		return true
	}

	return false
}

// AboutMoney rounds money to 1000億単位 (exported for external use)
// Ref: perl/lib/Hako/Main.pm:895
func AboutMoney(m int) string {
	if m < 500 {
		return fmt.Sprintf("推定500%s未満", hconst.UnitMoney)
	}
	m = (m + 500) / 1000
	return fmt.Sprintf("推定%d000%s", m, hconst.UnitMoney)
}

// HtmlEscape escapes HTML special characters (exported for external use)
// Ref: perl/lib/Hako/Main.pm:907
func HtmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// cutColumn cuts string to specified column length
// Ref: perl/lib/Hako/Main.pm:917
func cutColumn(s string, c int) string {
	// Phase 1: Simple implementation (works for ASCII and UTF-8)
	runes := []rune(s)
	if len(runes) <= c {
		return s
	}
	return string(runes[:c])
}

// nameToNumber finds island number by name (not ID, but index in Islands array)
// Ref: perl/lib/Hako/Main.pm:942
func nameToNumber(name string) int {
	// Search all islands
	for i := 0; i < variable.IslandNumber; i++ {
		island := variable.Islands[i]
		if island.Name == name {
			return i
		}
	}
	// Not found
	return -1
}

// MonsterSpec returns monster information from level (exported for external use)
// Ref: perl/lib/Hako/Main.pm:958
func MonsterSpec(lv int) (kind int, name string, hp int) {
	// Kind
	kind = lv / 10

	// Name
	if kind < len(hconst.MonsterName) {
		name = hconst.MonsterName[kind]
	}

	// HP
	hp = lv - (kind * 10)

	return kind, name, hp
}

// ExpToLevel calculates level from experience points (exported for external use)
// Ref: perl/lib/Hako/Main.pm:975
func ExpToLevel(landKind int, exp int) int {
	if landKind == hconst.LandBase {
		// Missile base
		for i := hconst.MaxBaseLevel; i > 1; i-- {
			if exp >= hconst.BaseLevelUp[i-2] {
				return i
			}
		}
		return 1
	} else {
		// Sea base
		for i := hconst.MaxSBaseLevel; i > 1; i-- {
			if exp >= hconst.SBaseLevelUp[i-2] {
				return i
			}
		}
		return 1
	}
}

// makeRandomPointArray creates shuffled coordinate arrays
// Sets (@Hrpx, @Hrpy) so that numbers from (0,0) to (size-1, size-1) appear exactly once
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

// random returns a random integer in [0, n)
// Ref: perl/lib/Hako/Main.pm:1022
func random(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.Intn(n)
}

// Suppress unused import warning for Phase 1
var _ = des.BlockSize
