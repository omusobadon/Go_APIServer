// 顧客情報のPOST
package post

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// リクエストを変換する構造体
type PostCustomerRequest struct {
	Name    string `json:"name"`
	Mail    string `json:"mail"`
	Phone   string `json:"phone"`
	Passwd  string `json:"password"`
	Address string `json:"address"`
	Payment string `json:"payment_info"`
}

// レスポンスに変換する構造体
type PostCustomerResponse struct {
	Message    string              `json:"message"`
	Request    PostCustomerRequest `json:"request"`
	Registered db.CustomerModel    `json:"registered"`
}

var post_customer_cnt int

func PostCustomer(w http.ResponseWriter, r *http.Request) {
	post_customer_cnt++
	var (
		status   int    = http.StatusNotImplemented
		message  string = "メッセージがありません"
		req      PostCustomerRequest
		customer db.CustomerModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// レスポンスボディ
		res := new(PostCustomerResponse)

		// レスポンスボディの作成
		res.Message = message
		res.Request = req
		res.Registered = customer

		// レスポンス構造体をJSONに変換して送信
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[Post Customer.%d][%d] %s\n", post_customer_cnt, status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	if req.Name == "" || req.Mail == "" || req.Phone == "" {
		status = http.StatusBadRequest
		message = "必要な情報がありません"
		return
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

	// 顧客情報の挿入
	c, err := client.Customer.CreateOne(
		db.Customer.Name.Set(req.Name),
		db.Customer.Mail.Set(req.Mail),
		db.Customer.Phone.Set(req.Phone),
		db.Customer.Password.Set(req.Passwd),
		db.Customer.Address.Set(req.Address),
		db.Customer.PaymentInfo.Set(req.Payment),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("顧客テーブル挿入エラー : ", err)
		return
	}

	customer = *c

	status = http.StatusOK
	message = "正常終了"
}
