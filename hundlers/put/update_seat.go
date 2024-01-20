package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateSeatRequest struct {
	ID      *int    `json:"id"`
	Product *int    `json:"product_id"`
	Row     *string `json:"row"`
	Column  *string `json:"column"`
	Enable  *bool   `json:"is_enable"`
	Remark  *string `json:"remark"`
}

type UpdateSeatResponseSuccess struct {
	Message    string            `json:"message"`
	Request    UpdateSeatRequest `json:"request"`
	Registered db.SeatModel      `json:"registered"`
}

type UpdateSeatResponseFailure struct {
	Message string            `json:"message"`
	Request UpdateSeatRequest `json:"request"`
}

func UpdateSeat(w http.ResponseWriter, r *http.Request) {
	var (
		status     int    = http.StatusNotImplemented
		message    string = "メッセージがありません"
		req        UpdateSeatRequest
		registered db.SeatModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(UpdateSeatResponseSuccess)

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
			res := new(UpdateSeatResponseFailure)

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

		fmt.Printf("[Update Seat][%d] %s\n", status, message)

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
	created, err := client.Seat.FindUnique(
		db.Seat.ID.EqualsIfPresent(req.ID),
	).Update(
		db.Seat.ProductID.SetIfPresent(req.Product),
		db.Seat.Row.SetIfPresent(req.Row),
		db.Seat.Column.SetIfPresent(req.Column),
		db.Seat.IsEnable.SetIfPresent(req.Enable),
		db.Seat.Remark.SetIfPresent(req.Remark),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Seatテーブル挿入エラー : ", err)
		return
	}

	registered = *created

	status = http.StatusOK
	message = "正常終了"
}
