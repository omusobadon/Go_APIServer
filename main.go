package main

import "fmt"

func main() {
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
