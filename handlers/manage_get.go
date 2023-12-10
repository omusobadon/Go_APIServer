package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type MGetResponseBody struct {
	Message  string             `json:"message"`
	Order    []db.OrderModel    `json:"order"`
	Stock    []db.StockModel    `json:"stock"`
	Product  []db.ProductModel  `json:"product"`
	Fee      []db.FeeModel      `json:"fee"`
	Payment  []db.PaymentModel  `json:"payment"`
	Customer []db.CustomerModel `json:"customer"`
}

var mget_cnt int // ManageGETのカウント用

func ManageGet(w http.ResponseWriter, r *http.Request) {
	var order []db.OrderModel
	var stock []db.StockModel
	var product []db.ProductModel
	var fee []db.FeeModel
	var payment []db.PaymentModel
	var customer []db.CustomerModel
	var status int
	var message string
	mget_cnt++

	fmt.Printf("# Manage Get No.%d #\n", mget_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		res := MGetResponseBody{
			Message:  message,
			Order:    order,
			Stock:    stock,
			Product:  product,
			Fee:      fee,
			Payment:  payment,
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

		// 処理結果メッセージの表示（サーバ側）
		if status == 0 || message == "" {
			fmt.Println("ステータスコードまたはメッセージがありません")
		} else {
			fmt.Printf("[%d] %s\n", status, message)
		}

		fmt.Printf("# Manage Get No.%d End #\n", mget_cnt)
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

	// Orderテーブルの内容を一括取得
	order, err := client.Order.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("注文テーブル取得エラー : ", err)
	}

	// Stockテーブルの内容を一括取得
	stock, err = client.Stock.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("在庫テーブル取得エラー : ", err)
		return
	}

	// Productテーブルの内容を一括取得
	product, err = client.Product.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品テーブル取得エラー : ", err)
		return
	}

	// Feeテーブルの内容を一括取得
	fee, err = client.Fee.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("料金テーブル取得エラー : ", err)
		return
	}

	// Paymentテーブルの内容を一括取得
	payment, err = client.Payment.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("決済処理テーブル取得エラー : ", err)
		return
	}

	// Customerテーブルの内容を一括取得
	customer, err = client.Customer.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("顧客情報テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
