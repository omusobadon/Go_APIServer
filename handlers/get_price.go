// 値段情報のGET
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
type GetPriceResponseBody struct {
	Message string          `json:"message"`
	Length  int             `json:"length"`
	Price   []db.PriceModel `json:"price"`
}

var get_price_cnt int // GetPriceの呼び出しカウント

func GetPrice(w http.ResponseWriter, r *http.Request) {
	get_price_cnt++
	var (
		price   []db.PriceModel
		status  int
		message string
		err     error
	)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetPriceResponseBody{
			Message: message,
			Length:  len(price),
			Price:   price,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[Get Price.%d][%d] %s\n", get_price_cnt, status, message)
	}()

	// リクエストパラメータの取得
	product_str := r.FormValue("product_id")

	// パラメータが空でない場合はIntに変換
	var product_id int
	if product_str != "" {
		product_id, err = strconv.Atoi(product_str)
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

	// productパラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if product_id == 0 {
		price, err = client.Price.FindMany().Exec(ctx)

	} else {
		price, err = client.Price.FindMany(
			db.Price.ProductID.Equals(product_id),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(price) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
