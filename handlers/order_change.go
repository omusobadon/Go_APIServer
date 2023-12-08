// Order変更（予約の終了、キャンセル処理）
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var change_cnt int // orderChangeのカウント用

func OrderChange(w http.ResponseWriter, r *http.Request) {
	var order *Order
	var res *Response
	change_cnt++

	fmt.Printf("*** Change No.%d ***\n", change_cnt)

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("POSTデコードエラー :", err)
		return
	}

	fmt.Println("変更情報 :", order)

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("クライアント接続エラー :", err)
		return
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー :", err))
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Status)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Message = fmt.Sprint("レスポンスの作成エラー :", err)
		}

		// 処理結果メッセージの表示（サーバ側）
		if res.Status == 0 || res.Message == "" {
			fmt.Println("ステータスコードまたはメッセージがありません")
		} else {
			fmt.Printf("[%d] %s\n", res.Status, res.Message)
		}

		fmt.Printf("*** Change No.%d End ***\n", change_cnt)

	}()

	ctx := context.Background()

	// OrderIDから注文情報を取得
	order_info, err := client.Order.FindUnique(
		db.Order.ID.Equals(order.ID),
	).Exec(ctx)
	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("注文テーブル取得エラー :", err)
		return
	}

	// StockIDから在庫情報を取得
	stock_info, err := client.Stock.FindUnique(
		db.Stock.ID.Equals(order_info.InnerOrder.Product),
	).Exec(ctx)
	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("在庫テーブル取得エラー :", err)
		return
	}

	// 変更処理 [2]:予約完了, [3]:予約キャンセル
	switch order.State {
	case 2:
	case 3:
	default:
		res.Status = http.StatusBadRequest
		res.Message = "不正なステータス"
		return
	}

	// Orderテーブルのステータスを変更
	_, err = client.Order.FindUnique(
		db.Order.ID.Equals(order.ID),
	).Update(
		db.Order.State.Set(order.State),
	).Exec(ctx)
	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("注文テーブルアップデートエラー :", err)
		return
	}

	// 在庫を元に戻す
	id := order_info.InnerOrder.Product
	num := stock_info.InnerStock.Num + order_info.InnerOrder.Num

	_, err = client.Stock.FindUnique(
		db.Stock.ID.Equals(id),
	).Update(
		db.Stock.Num.Set(num),
	).Exec(ctx)
	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = fmt.Sprint("在庫テーブルアップデートエラー :", err)
	}

	// 正常終了時
	res.Status = http.StatusOK
	res.Message = "正常終了"

}
