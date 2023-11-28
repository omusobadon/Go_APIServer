// 各テーブル用の構造体まとめ
package main

import (
	"time"
)

// 注文情報テーブル
// State : [1]:予約受付, [2]:予約終了, [3]:予約キャンセル
type Order struct {
	ID       int       `json:"id"`
	Customer int       `json:"customer"`
	Product  int       `json:"product"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Num      int       `json:"num"`
	Time     time.Time `json:"time"`
	State    int       `json:"state"`
	Note     string    `json:"note"`
}

// 在庫テーブル
type Stock struct {
	ID       *int       `json:"id"`
	Category *string    `json:"category"`
	Name     *string    `json:"name"`
	Interval *time.Time `json:"interval"`
	Value    *int       `json:"value"`
	Num      *int       `json:"num"`
	Note     *string    `json:"note"`
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
