// 各テーブル用の構造体まとめ
package main

import (
	"time"
)

// 注文情報テーブル
type Order struct {
	ID       int       `json:"id"`
	Customer int       `json:"customer"`
	Product  int       `json:"product"`
	Num      int       `json:"num"`
	Time     time.Time `json:"time"`
}

// 在庫テーブル
type Stock struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
}

// POSTされるDB管理情報
// Type 1:Update, 2:Insert, 3:Delete
// Table テーブル名
// Info 更新内容
type EditInfo struct {
	ID    int            `json:"id"`
	Table string         `json:"table"`
	Type  int            `json:"type"`
	Info  map[string]any `json:"info"`
	Time  time.Time      `json:"time"`
}
