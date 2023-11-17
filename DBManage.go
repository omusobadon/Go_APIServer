package main

import (
	"fmt"

	"Go_APIServer/db"
)

// 更新処理を行い、結果を返す
func (info *EditInfo) TableEdit() EditResponse {
	res := EditResponse{Info: info}

	// 更新時刻を取得
	info.Time = GetTime()

	// DBクライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Println("クライアント接続エラー :", err)
		res.Status = 30
		return res
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	// Stockテーブルの処理 -------------------------------------------------
	if info.Table == "stock" {
		var stock Stock

		fmt.Print("[Stock] ")

		// Type [1]:Update, [2]:Insert, [3]:Delete
		if info.Type == 1 {
			fmt.Println("Update :", stock)

			_, err := client.Stock.FindUnique(
				db.Stock.ProductID.Equals(stock.ID),
			).Update(
				db.Stock.ProductName.Set(stock.Name),
				db.Stock.StockNum.Set(stock.Num),
			).Exec(ctx)
			if err != nil {
				fmt.Println("StockテーブルUpdateエラー :", err)
				res.Status = 30
				return res
			}

		} else if info.Type == 2 {
			fmt.Println("Insert :", stock)

			// StockテーブルをInsert
			_, err := client.Stock.CreateOne(
				db.Stock.ProductName.Set(stock.Name),
				db.Stock.StockNum.Set(stock.Num),
			).Exec(ctx)
			if err != nil {
				fmt.Println("StockテーブルInsertエラー :", err)
				res.Status = 30
				return res
			}

		} else if info.Type == 3 {
			fmt.Println("Delete :", stock)

			// StockテーブルをDelete
			_, err := client.Stock.FindUnique(
				db.Stock.ProductID.Equals(stock.ID),
			).Delete().Exec(ctx)
			if err != nil {
				fmt.Println("StockテーブルDeleteエラー :", err)
				res.Status = 30
				return res
			}

		} else {
			fmt.Println("エラー : Typeが見つかりません")
			res.Status = 30
			return res
		}
	}

	// 処理が正常終了したらEditInfoテーブルに登録
	_, err := client.Manage.CreateOne(
		db.Manage.Table.Set(m.Table),
		db.Manage.Type.Set(m.Type),
		db.Manage.Info.Set(m.Info),
		db.Manage.Time.Set(m.Time),
	).Exec(ctx)
	if err != nil {
		fmt.Println("ManageテーブルInsertエラー :", err)
		res.Status = 30
		return res
	}

	return res
}
