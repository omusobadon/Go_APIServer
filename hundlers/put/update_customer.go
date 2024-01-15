package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateCustomerRequest struct {
	ID       *int    `json:"id"`
	Name     *string `json:"name"`
	Mail     *string `json:"mail"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
	Address  *string `json:"address"`
	Payment  *string `json:"payment_info"`
}

type UpdateCustomerResponseSuccess struct {
	Message    string                `json:"message"`
	Request    UpdateCustomerRequest `json:"request"`
	Registered db.CustomerModel      `json:"registered"`
}

type UpdateCustomerResponseFailure struct {
	Message string                `json:"message"`
	Request UpdateCustomerRequest `json:"request"`
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		status   int    = http.StatusNotImplemented
		message  string = "メッセージがありません"
		req      UpdateCustomerRequest
		customer db.CustomerModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(UpdateCustomerResponseSuccess)

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

		} else {
			res := new(UpdateCustomerResponseFailure)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}
		}

		fmt.Printf("[Update Customer][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.ID == nil {
		status = http.StatusBadRequest
		message = "id is null"
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
	created, err := client.Customer.FindUnique(
		db.Customer.ID.EqualsIfPresent(req.ID),
	).Update(
		db.Customer.Name.SetIfPresent(req.Name),
		db.Customer.Mail.SetIfPresent(req.Mail),
		db.Customer.Phone.SetIfPresent(req.Phone),
		db.Customer.Password.SetIfPresent(req.Password),
		db.Customer.Address.SetIfPresent(req.Address),
		db.Customer.PaymentInfo.SetIfPresent(req.Payment),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Customerテーブル挿入エラー : ", err)
		return
	}

	customer = *created

	status = http.StatusOK
	message = "正常終了"
}
