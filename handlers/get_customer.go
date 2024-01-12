// 商品グループ情報のGET
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type GetCustomerResponseBody struct {
	Message  string             `json:"message"`
	Length   int                `json:"length"`
	Customer []db.CustomerModel `json:"customer"`
}

var get_customer_cnt int // ShopGetの呼び出しカウント

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	get_customer_cnt++
	var (
		customer []db.CustomerModel
		status   int
		message  string
		err      error
	)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetCustomerResponseBody{
			Message:  message,
			Length:   len(customer),
			Customer: customer,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[Get Customer.%d][%d] %s\n", get_customer_cnt, status, message)
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

	customer, err = client.Customer.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Customerテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	// if len(customer) == 0 {
	// 	status = http.StatusBadRequest
	// 	message = "顧客情報がありません"
	// 	return
	// }

	status = http.StatusOK
	message = "正常終了"
}
