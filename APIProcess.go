package main

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"net/http"
)

type HTTPError struct {
	Status  int
	Message string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("[%d] : %s", err.Status, err.Message)
}

// 注文処理 Orderテーブルを受け取って注文処理を行い、ステータスコードとメッセージを返す
func (order *Order) Process() error {

	// 注文時刻を取得
	order.Time = GetTime()

	fmt.Println("注文情報 :", order)

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		message = fmt.Sprint("クライアント接続エラー :", err)
		status = http.StatusBadRequest
		return
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// 注文情報の商品idと一致する在庫情報を取得
	stock, err := client.Stock.FindUnique(
		db.Stock.ID.Equals(order.Product),
	).Exec(ctx)
	if err != nil {
		message = fmt.Sprint("在庫テーブル取得エラー : ", err)
		status = http.StatusBadRequest
		return
	}

	// 在庫が注文数を上回っていたら注文処理を行う
	if stock.InnerStock.Num >= order.Num {
		// 在庫テーブルに注文情報を反映
		_, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(order.Product),
		).Update(
			db.Stock.Num.Set(stock.InnerStock.Num - order.Num),
		).Exec(ctx)
		if err != nil {
			message = fmt.Sprint("在庫テーブルアップデートエラー :", err)
			status = http.StatusBadRequest
			return
		}

		// 注文テーブルに注文情報をインサート
		if err := order.Insert(client); err != nil {
			message = fmt.Sprint("注文テーブルインサートエラー :", err)
			status = http.StatusBadRequest

			// 注文を登録できなかった場合に在庫の数量を戻す
			_, err := client.Stock.FindMany(
				db.Stock.ID.Equals(order.Product),
			).Update(
				db.Stock.Num.Set(stock.InnerStock.Num + order.Num),
			).Exec(ctx)
			if err != nil {
				message = fmt.Sprint("在庫整合性エラー :", err)
				status = http.StatusInternalServerError
				return
			}
			return
		}

		// 正常終了のとき
		// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
		order_info, err := client.Order.FindFirst(
			db.Order.Customer.Equals(order.Customer),
			db.Order.Time.Equals(order.Time),
		).Exec(ctx)
		if err != nil {
			message = fmt.Sprint("注文情報取得エラー :", err)
			status = http.StatusInternalServerError
			return
		}

		order.ID = order_info.ID
		fmt.Println("注文受付 :", order)

		message = "正常終了"
		status = http.StatusOK

	} else {
		// 在庫不足のとき
		message = "在庫不足"
		status = http.StatusBadRequest
	}
}
