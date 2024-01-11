// 在庫情報のGET
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// レスポンスに変換する構造体
type GetStockResponseBody struct {
	Message string          `json:"message"`
	Length  int             `json:"length"`
	Stock   []db.StockModel `json:"stock"`
}

var get_stock_cnt int // orderGETのカウント用

func GetStock(w http.ResponseWriter, r *http.Request) {
	get_stock_cnt++
	var (
		stock   []db.StockModel
		status  int
		message string
		err     error
	)

	fmt.Printf("* Get Stock No.%d *\n", get_stock_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		res := GetStockResponseBody{
			Message: message,
			Length:  len(stock),
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

		fmt.Printf("[Get Stock.%d][%d] %s\n", get_stock_cnt, status, message)
	}()

	// リクエストパラメータの取得
	price_str := r.FormValue("price_id")

	// パラメータが空でない場合はIntに変換
	var price_id int
	if price_str != "" {
		price_id, err = strconv.Atoi(price_str)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("不正なパラメータ : ", err)
			return
		}
	}

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

	// priceパラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if price_id == 0 {
		stock, err = client.Stock.FindMany().Exec(ctx)

	} else {
		stock, err = client.Stock.FindMany(
			db.Stock.PriceID.Equals(price_id),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(stock) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
