// 在庫情報のGET
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type StockGetResponseBody struct {
	Message string          `json:"message"`
	Stock   []db.StockModel `json:"stock"`
}

var stock_get_cnt int // orderGETのカウント用

func StockGet(w http.ResponseWriter, r *http.Request) {
	var stock []db.StockModel
	var status int
	var message string
	stock_get_cnt++

	fmt.Printf("* Stock Get No.%d *\n", stock_get_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		res := StockGetResponseBody{
			Message: message,
			Stock:   stock,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		// 処理結果メッセージの表示（サーバ側）
		if status == 0 || message == "" {
			fmt.Println("ステータスコードまたはメッセージがありません")
		} else {
			fmt.Printf("[%d] %s\n", status, message)
		}

		fmt.Printf("* Stock Get No.%d End *\n", stock_get_cnt)
	}()

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("クライアント接続エラー : ", err)
		return
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー : ", err))
		}
	}()

	ctx := context.Background()

	// Stockテーブルの内容を一括取得
	stock, err := client.Stock.FindMany().With(
		db.Stock.Price.Fetch().With(
			db.Price.Product.Fetch().With(
				db.Product.Group.Fetch(),
			),
		),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("在庫テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
