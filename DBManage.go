package main

import (
	"context"
	"fmt"
	"time"

	"Go_APIServer/db"
)

// DB管理情報
// Type 0:更新, 1:インサート, 2:削除
// Table 編集テーブル名
// Info 更新内容
type ManageInfo struct {
	Type  int            `jdon:"type"`
	Table string         `json:"table"`
	Info  map[string]any `json:"info"`
}

type ManageRes struct {
	Status int        `json:"status"`
	Time   time.Time  `json:"time"`
	Info   ManageInfo `json:"info"`
}

// 更新処理を行い、結果を返す
func DBManage(info ManageInfo) ManageRes {
	var res ManageRes

	// 更新時刻を取得
	res.Time = GetTime()

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Println("クライアント接続エラー :", err)
		res.Status = 30
		return res
	}

	// 関数の終了時に実行
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// Stockテーブルの処理 -------------------------------------------------
	if info.Table == "stock" {
		fmt.Print("[Stock] ")

		// 更新処理
		if info.Type == 0 {
			fmt.Println("更新 :", info.Info)

			var ok bool
			var stock Stock
			for index, value := range info.Info {
				fmt.Printf("%s : %T\n", index, value)
			}

			id, ok := info.Info["id"].(float64)
			name, ok := info.Info["name"].(string)
			num, ok := info.Info["num"].(float64)
			if !ok {
				fmt.Println("型アサーションエラー")
				res.Status = 30
				return res
			}

			stock.ID = int(id)
			stock.Name = name
			stock.Num = int(num)

			// Stockテーブルを更新
			_, err := client.Stock.FindUnique(
				db.Stock.ProductID.Equals(stock.ID),
			).Update(
				db.Stock.ProductName.Set(stock.Name),
				db.Stock.StockNum.Set(stock.Num),
			).Exec(ctx)
			if err != nil {
				fmt.Println("在庫テーブルアップデートエラー :", err)
				res.Status = 30
				return res
			}
		}
	}

	res.Status = 10
	return res
}
