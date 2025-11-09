package core

import (
	"testing"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// TestMonsterSpec tests the MonsterSpec function
func TestMonsterSpec(t *testing.T) {
	tests := []struct {
		name       string
		lv         int
		wantKind   int
		wantName   string
		wantHP     int
	}{
		{
			name:     "メカいのら HP 0",
			lv:       0,
			wantKind: 0,
			wantName: "メカいのら",
			wantHP:   0,
		},
		{
			name:     "メカいのら HP 5",
			lv:       5,
			wantKind: 0,
			wantName: "メカいのら",
			wantHP:   5,
		},
		{
			name:     "いのら HP 3",
			lv:       13,
			wantKind: 1,
			wantName: "いのら",
			wantHP:   3,
		},
		{
			name:     "サンジラ HP 7",
			lv:       27,
			wantKind: 2,
			wantName: "サンジラ",
			wantHP:   7,
		},
		{
			name:     "レッドいのら HP 9",
			lv:       39,
			wantKind: 3,
			wantName: "レッドいのら",
			wantHP:   9,
		},
		{
			name:     "ダークいのら HP 2",
			lv:       42,
			wantKind: 4,
			wantName: "ダークいのら",
			wantHP:   2,
		},
		{
			name:     "いのらゴースト HP 0",
			lv:       50,
			wantKind: 5,
			wantName: "いのらゴースト",
			wantHP:   0,
		},
		{
			name:     "クジラ HP 8",
			lv:       68,
			wantKind: 6,
			wantName: "クジラ",
			wantHP:   8,
		},
		{
			name:     "キングいのら HP 5",
			lv:       75,
			wantKind: 7,
			wantName: "キングいのら",
			wantHP:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, name, hp := MonsterSpec(tt.lv)
			if kind != tt.wantKind {
				t.Errorf("MonsterSpec(%d) kind = %d, want %d", tt.lv, kind, tt.wantKind)
			}
			if name != tt.wantName {
				t.Errorf("MonsterSpec(%d) name = %s, want %s", tt.lv, name, tt.wantName)
			}
			if hp != tt.wantHP {
				t.Errorf("MonsterSpec(%d) hp = %d, want %d", tt.lv, hp, tt.wantHP)
			}
		})
	}
}

// TestExpToLevel tests the ExpToLevel function
func TestExpToLevel(t *testing.T) {
	tests := []struct {
		name     string
		landKind int
		exp      int
		wantLv   int
	}{
		// Missile Base tests (LandBase)
		{
			name:     "ミサイル基地 Lv1 (exp=0)",
			landKind: hconst.LandBase,
			exp:      0,
			wantLv:   1,
		},
		{
			name:     "ミサイル基地 Lv1 (exp=19)",
			landKind: hconst.LandBase,
			exp:      19,
			wantLv:   1,
		},
		{
			name:     "ミサイル基地 Lv2 (exp=20)",
			landKind: hconst.LandBase,
			exp:      20,
			wantLv:   2,
		},
		{
			name:     "ミサイル基地 Lv2 (exp=59)",
			landKind: hconst.LandBase,
			exp:      59,
			wantLv:   2,
		},
		{
			name:     "ミサイル基地 Lv3 (exp=60)",
			landKind: hconst.LandBase,
			exp:      60,
			wantLv:   3,
		},
		{
			name:     "ミサイル基地 Lv3 (exp=119)",
			landKind: hconst.LandBase,
			exp:      119,
			wantLv:   3,
		},
		{
			name:     "ミサイル基地 Lv4 (exp=120)",
			landKind: hconst.LandBase,
			exp:      120,
			wantLv:   4,
		},
		{
			name:     "ミサイル基地 Lv4 (exp=199)",
			landKind: hconst.LandBase,
			exp:      199,
			wantLv:   4,
		},
		{
			name:     "ミサイル基地 Lv5 (exp=200)",
			landKind: hconst.LandBase,
			exp:      200,
			wantLv:   5,
		},
		{
			name:     "ミサイル基地 Lv5 (exp=1000)",
			landKind: hconst.LandBase,
			exp:      1000,
			wantLv:   5,
		},
		// Sea Base tests (LandSbase)
		{
			name:     "海底基地 Lv1 (exp=0)",
			landKind: hconst.LandSbase,
			exp:      0,
			wantLv:   1,
		},
		{
			name:     "海底基地 Lv1 (exp=49)",
			landKind: hconst.LandSbase,
			exp:      49,
			wantLv:   1,
		},
		{
			name:     "海底基地 Lv2 (exp=50)",
			landKind: hconst.LandSbase,
			exp:      50,
			wantLv:   2,
		},
		{
			name:     "海底基地 Lv2 (exp=199)",
			landKind: hconst.LandSbase,
			exp:      199,
			wantLv:   2,
		},
		{
			name:     "海底基地 Lv3 (exp=200)",
			landKind: hconst.LandSbase,
			exp:      200,
			wantLv:   3,
		},
		{
			name:     "海底基地 Lv3 (exp=1000)",
			landKind: hconst.LandSbase,
			exp:      1000,
			wantLv:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := ExpToLevel(tt.landKind, tt.exp)
			if lv != tt.wantLv {
				t.Errorf("ExpToLevel(%d, %d) = %d, want %d", tt.landKind, tt.exp, lv, tt.wantLv)
			}
		})
	}
}

// TestAboutMoney tests the AboutMoney function
func TestAboutMoney(t *testing.T) {
	tests := []struct {
		name  string
		money int
		want  string
	}{
		{
			name:  "100億円未満",
			money: 100,
			want:  "推定500億円未満",
		},
		{
			name:  "499億円も500億円未満",
			money: 499,
			want:  "推定500億円未満",
		},
		{
			name:  "500億円",
			money: 500,
			want:  "推定1000億円",
		},
		{
			name:  "999億円",
			money: 999,
			want:  "推定1000億円",
		},
		{
			name:  "1000億円",
			money: 1000,
			want:  "推定1000億円",
		},
		{
			name:  "1499億円",
			money: 1499,
			want:  "推定1000億円",
		},
		{
			name:  "1500億円",
			money: 1500,
			want:  "推定2000億円",
		},
		{
			name:  "5000億円",
			money: 5000,
			want:  "推定5000億円",
		},
		{
			name:  "9999億円",
			money: 9999,
			want:  "推定10000億円",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AboutMoney(tt.money)
			if got != tt.want {
				t.Errorf("AboutMoney(%d) = %s, want %s", tt.money, got, tt.want)
			}
		})
	}
}

// TestHtmlEscape tests the HtmlEscape function
func TestHtmlEscape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "特殊文字なし",
			input: "Hello World",
			want:  "Hello World",
		},
		{
			name:  "アンパサンド",
			input: "A & B",
			want:  "A &amp; B",
		},
		{
			name:  "タグ",
			input: "<script>alert('XSS')</script>",
			want:  "&lt;script&gt;alert('XSS')&lt;/script&gt;",
		},
		{
			name:  "ダブルクォート",
			input: `Say "Hello"`,
			want:  `Say &quot;Hello&quot;`,
		},
		{
			name:  "全て混在",
			input: `<a href="page.html?a=1&b=2">Link</a>`,
			want:  `&lt;a href=&quot;page.html?a=1&amp;b=2&quot;&gt;Link&lt;/a&gt;`,
		},
		{
			name:  "日本語とタグ",
			input: "<div>こんにちは</div>",
			want:  "&lt;div&gt;こんにちは&lt;/div&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HtmlEscape(tt.input)
			if got != tt.want {
				t.Errorf("HtmlEscape(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestCheckPassword tests the CheckPassword function
func TestCheckPassword(t *testing.T) {
	// Set master password for testing
	oldMaster := hconst.MasterPassword
	hconst.MasterPassword = "master123"
	defer func() { hconst.MasterPassword = oldMaster }()

	// Note: hconst.CryptOn is a const and cannot be changed in tests
	// In Phase 1, cryptCompat returns password as-is, so testing works the same

	tests := []struct {
		name   string
		stored string
		input  string
		want   bool
	}{
		{
			name:   "正しいパスワード",
			stored: "password123",
			input:  "password123",
			want:   true,
		},
		{
			name:   "間違ったパスワード",
			stored: "password123",
			input:  "wrongpassword",
			want:   false,
		},
		{
			name:   "空のパスワード",
			stored: "password123",
			input:  "",
			want:   false,
		},
		{
			name:   "マスターパスワード",
			stored: "password123",
			input:  "master123",
			want:   true,
		},
		{
			name:   "大文字小文字区別",
			stored: "Password",
			input:  "password",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPassword(tt.stored, tt.input)
			if got != tt.want {
				t.Errorf("CheckPassword(%q, %q) = %v, want %v", tt.stored, tt.input, got, tt.want)
			}
		})
	}
}

// TestMakeRandomPointArray tests the MakeRandomPointArray function
func TestMakeRandomPointArray(t *testing.T) {
	// Initialize variable package
	variable.Rpx = nil
	variable.Rpy = nil

	MakeRandomPointArray()

	// Check that arrays are created
	if variable.Rpx == nil {
		t.Fatal("Rpx is nil after MakeRandomPointArray")
	}
	if variable.Rpy == nil {
		t.Fatal("Rpy is nil after MakeRandomPointArray")
	}

	// Check length
	if len(variable.Rpx) != hconst.PointNumber {
		t.Errorf("len(Rpx) = %d, want %d", len(variable.Rpx), hconst.PointNumber)
	}
	if len(variable.Rpy) != hconst.PointNumber {
		t.Errorf("len(Rpy) = %d, want %d", len(variable.Rpy), hconst.PointNumber)
	}

	// Check that all coordinates appear exactly once
	coordMap := make(map[[2]int]int)
	for i := 0; i < hconst.PointNumber; i++ {
		x := variable.Rpx[i]
		y := variable.Rpy[i]

		// Check range
		if x < 0 || x >= hconst.IslandSize {
			t.Errorf("Rpx[%d] = %d, out of range [0, %d)", i, x, hconst.IslandSize)
		}
		if y < 0 || y >= hconst.IslandSize {
			t.Errorf("Rpy[%d] = %d, out of range [0, %d)", i, y, hconst.IslandSize)
		}

		// Count occurrences
		coord := [2]int{x, y}
		coordMap[coord]++
	}

	// Check that all coordinates appear exactly once
	if len(coordMap) != hconst.PointNumber {
		t.Errorf("unique coordinates = %d, want %d", len(coordMap), hconst.PointNumber)
	}

	for coord, count := range coordMap {
		if count != 1 {
			t.Errorf("coordinate %v appears %d times, want 1", coord, count)
		}
	}

	// Check that the result is shuffled (not in sequential order)
	// This is a probabilistic test - it might fail rarely
	sequential := true
	expectedX := 0
	expectedY := 0
	for i := 0; i < hconst.PointNumber; i++ {
		if variable.Rpx[i] != expectedX || variable.Rpy[i] != expectedY {
			sequential = false
			break
		}
		expectedX++
		if expectedX >= hconst.IslandSize {
			expectedX = 0
			expectedY++
		}
	}

	if sequential {
		t.Log("Warning: Coordinates are in sequential order (might be random, but unlikely)")
	}
}

// TestMin tests the min function
func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a < b",
			a:    5,
			b:    10,
			want: 5,
		},
		{
			name: "a > b",
			a:    10,
			b:    5,
			want: 5,
		},
		{
			name: "a == b",
			a:    7,
			b:    7,
			want: 7,
		},
		{
			name: "負の数",
			a:    -5,
			b:    3,
			want: -5,
		},
		{
			name: "0とゼロ",
			a:    0,
			b:    0,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := min(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestCutColumn tests the cutColumn function
func TestCutColumn(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		column int
		want   string
	}{
		{
			name:   "短い文字列",
			input:  "Hello",
			column: 10,
			want:   "Hello",
		},
		{
			name:   "ちょうど同じ長さ",
			input:  "12345",
			column: 5,
			want:   "12345",
		},
		{
			name:   "切り詰め",
			input:  "1234567890",
			column: 5,
			want:   "12345",
		},
		{
			name:   "日本語（UTF-8）",
			input:  "こんにちは世界",
			column: 4,
			want:   "こんにち",
		},
		{
			name:   "空文字列",
			input:  "",
			column: 5,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cutColumn(tt.input, tt.column)
			if got != tt.want {
				t.Errorf("cutColumn(%q, %d) = %q, want %q", tt.input, tt.column, got, tt.want)
			}
		})
	}
}
