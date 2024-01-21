package delete

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DeleteCustomerResponseSuccess struct {
	Message string           `json:"message"`
	Deleted db.CustomerModel `json:"deleted"`
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		err     error
		deleted db.CustomerModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 処理成功時
		if status == http.StatusOK {
			res := new(DeleteCustomerResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Deleted = deleted

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(DeleteResponseFailure)

			// レスポンスボディの作成
			res.Message = message

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}
		}

		fmt.Printf("[Delete Customer][%d] %s\n", status, message)

	}()

	// リクエストパラメータの取得
	id_str := r.FormValue("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		status = http.StatusBadRequest
		message = "不正なパラメータ"
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

	// Delete
	d, err := client.Customer.FindUnique(
		db.Customer.ID.Equals(id),
	).With(
		db.Customer.Order.Fetch().With(
			db.Order.PaymentState.Fetch(),
			db.Order.ReservationCancel.Fetch(),
			db.Order.ReservationEnd.Fetch(),
			db.Order.OrderDetail.Fetch(),
		),
	).Delete().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Customer削除エラー : ", err)
		return
	}

	deleted = *d

	status = http.StatusOK
	message = "正常終了"
}
