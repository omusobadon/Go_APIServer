// 各テーブル用の構造体まとめ
package handlers

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
	ID       *int       `json:"id"`       // 商品ID
	Category *string    `json:"category"` // 商品カテゴリ
	Name     *string    `json:"name"`     // 商品名
	Start    *time.Time `json:"start"`    // 開始時刻
	End      *time.Time `json:"end"`      // 終了時刻

	// 終了時刻から次の開始時刻までのインターバル
	// 有効な値が入力されている場合は、インターバル後に自動的にDefaultNumでリセット
	Interval   *string `json:"interval"`
	Value      *int    `json:"value"`       // 価格
	DefaultNum *int    `json:"default_num"` // デフォルトの在庫数
	Num        *int    `json:"num"`         // 現在の在庫数
	State      *bool   `json:"state"`       // 無効の場合は注文を受け付けない
	Note       *string `json:"note"`        // 備考
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
