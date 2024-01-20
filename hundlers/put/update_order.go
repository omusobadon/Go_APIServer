package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateOrderRequest struct {
	ID       *int       `json:"id"`
	Customer *int       `json:"customer_id"`
	Start    *time.Time `json:"start_at"`
	End      *time.Time `json:"end_at"`
	Accept   *bool      `json:"is_accepted"`
	Pend     *bool      `json:"is_pending"`
	Remark   *string    `json:"remark"`
}

type UpdateOrderResponseSuccess struct {
	Message    string             `json:"message"`
	Request    UpdateOrderRequest `json:"request"`
	Registered db.OrderModel      `json:"registered"`
}

type UpdateOrderResponseFailure struct {
	Message string             `json:"message"`
	Request UpdateOrderRequest `json:"request"`
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var (
		status     int    = http.StatusNotImplemented
		message    string = "メッセージがありません"
		req        UpdateOrderRequest
		registered db.OrderModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(UpdateOrderResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = registered

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(UpdateOrderResponseFailure)

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

		fmt.Printf("[Update Order][%d] %s\n", status, message)

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
	created, err := client.Order.FindUnique(
		db.Order.ID.EqualsIfPresent(req.ID),
	).Update(
		db.Order.CustomerID.SetIfPresent(req.Customer),
		db.Order.StartAt.SetIfPresent(req.Start),
		db.Order.EndAt.SetIfPresent(req.End),
		db.Order.IsAccepted.SetIfPresent(req.Accept),
		db.Order.IsPending.SetIfPresent(req.Pend),
		db.Order.Remark.SetIfPresent(req.Remark),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Orderテーブル挿入エラー : ", err)
		return
	}

	registered = *created

	status = http.StatusOK
	message = "正常終了"
}
