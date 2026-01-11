// Package types defines data types for the Hakoniwa game.
// This file contains Island and related type definitions
//
// Ref: docs/migration/DATA_FORMAT.md
package types

import "github.com/neguse/hakoniwa/internal/hako/hconst"

// Island represents an island's data
// Ref: perl/lib/Hako/Main.pm, docs/migration/DATA_FORMAT.md
type Island struct {
	// Basic info (from hakojima.dat)
	Name     string // [0] 島の名前
	ID       string // [1] 島のID
	Prize    int    // [2] 受賞フラグ
	Absent   int    // [3] 連続資金繰り数
	Comment  string // [4] コメント
	Password string // [5] 暗号化パスワード
	Money    int    // [6] 資金（100億円単位）
	Food     int    // [7] 食料（100トン単位）
	Pop      int    // [8] 人口（100人単位）
	Area     int    // [9] 広さ（100平方km単位）
	Farm     int    // [10] 農場数
	Factory  int    // [11] 工場数
	Mountain int    // [12] 採掘場数
	Score    int    // スコア（ソート用）

	// Turn processing fields
	OldPop      int  // ターン開始前の人口
	Dead        bool // 死滅フラグ
	BigMissile  int  // 巨大ミサイル着弾予定数
	MonsterSend int  // 怪獣派遣カウンタ
	Propaganda  int  // 誘致活動フラグ
	Prepare2    int  // 地ならしカウンタ

	// Land data (from island.{ID})
	Land      [][]int // 地形タイプ (11x11 or 12x12)
	LandValue [][]int // 地形の値 (11x11 or 12x12)

	// Commands (from island.{ID})
	Commands []Command // コマンドキュー

	// Local BBS (from island.{ID})
	Lbbs []LbbsEntry // ローカル掲示板
}

// Command represents a development command
type Command struct {
	Kind   int    // コマンド種類
	Target string // ターゲットID
	X      int    // X座標
	Y      int    // Y座標
	Arg    int    // 引数
}

// LbbsEntry represents a local BBS entry
type LbbsEntry struct {
	Name    string // 投稿者名
	Message string // メッセージ
}

// NewIsland creates a new island with default values
func NewIsland(name, id, password string) *Island {
	size := hconst.IslandSize
	island := &Island{
		Name:     name,
		ID:       id,
		Password: password,
		Money:    hconst.InitialMoney,
		Food:     hconst.InitialFood,
		Area:     size * size,
		Land:     make([][]int, size),
		LandValue: make([][]int, size),
		Commands: make([]Command, 0, hconst.CommandMax),
		Lbbs:     make([]LbbsEntry, 0, hconst.LbbsMax),
	}

	// Initialize land
	for i := 0; i < size; i++ {
		island.Land[i] = make([]int, size)
		island.LandValue[i] = make([]int, size)
		for j := 0; j < size; j++ {
			island.Land[i][j] = hconst.LandSea
			island.LandValue[i][j] = 0
		}
	}

	return island
}
