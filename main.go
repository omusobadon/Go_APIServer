package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// config.ymlデコード用構造体
type Options struct {
	Time_free_enable  bool
	Seat_enable       bool
	Payment_enable    bool
	User_end_enable   bool
	User_notification bool
	Timezone          string
	Default_delay     int
}

// 商品・在庫テーブルが空の場合、自動生成するかどうか（テスト用）
const auto_insert bool = true

var OPTIONS Options

func main() {

	// config.ymlの読み込み
	content, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	// ymlのデコード
	yaml.Unmarshal(content, &OPTIONS)

	// optionの確認
	fmt.Printf("[OPTIONS]: %+v\n", OPTIONS)

	// テーブルに情報がない場合に自動インサート（テスト用）
	if auto_insert {
		if err := AutoInsert(); err != nil {
			fmt.Println(err)
			fmt.Println("処理を続行します")
		}
	}

	// 各テーブルのチェック
	// 未実装

	// スケジューラの起動
	go scheduler()

	// APIServerの起動
	if err := APIServer(); err != nil {
		panic(err)
	}
}
