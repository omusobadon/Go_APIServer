package main

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"net/http"
)

// レスポンス用のHTTPステータスコードとメッセージ
type Response struct {
	Status  int
	Message string
}

// 注文処理 Orderテーブルを受け取って注文処理を行い、ステータスコードとメッセージを返す
func (order *Order) Process() *Response {

	// 注文時刻を取得
	order.Time = GetTime()

	fmt.Println("注文情報 :", order)

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("クライアント接続エラー :", err),
		}
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー :", err))
		}
	}()

	ctx := context.Background()

	// 注文情報の商品idと一致する在庫情報を取得
	stock, err := client.Stock.FindUnique(
		db.Stock.ID.Equals(order.Product),
	).Exec(ctx)
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprint("在庫テーブル取得エラー : ", err),
		}
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
			return &Response{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprint("在庫テーブルアップデートエラー :", err),
			}
		}

		// 注文テーブルに注文情報をインサート
		if err := order.Insert(client); err != nil {

			// 注文を登録できなかった場合に在庫の数量を戻す
			_, err := client.Stock.FindMany(
				db.Stock.ID.Equals(order.Product),
			).Update(
				db.Stock.Num.Set(stock.InnerStock.Num + order.Num),
			).Exec(ctx)
			if err != nil {
				panic(fmt.Sprint("在庫整合性エラー :", err))
			}

			return &Response{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprint("注文テーブルインサートエラー :", err),
			}
		}

		// 正常終了のとき
		// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
		order_info, err := client.Order.FindFirst(
			db.Order.Customer.Equals(order.Customer),
			db.Order.Time.Equals(order.Time),
		).Exec(ctx)
		if err != nil {
			return &Response{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprint("注文情報取得エラー :", err),
			}
		}

		order.ID = order_info.ID
		fmt.Println("注文受付 :", order)

		return &Response{
			Status:  http.StatusOK,
			Message: "正常終了",
		}

	} else {
		// 在庫不足のとき
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "在庫不足",
		}
	}
}
