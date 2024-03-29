package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateSeatRequest struct {
	Product *int    `json:"product_id"`
	Row     *string `json:"row"`
	Column  *string `json:"column"`
	Enable  *bool   `json:"is_enable"`
	Remark  *string `json:"remark"`
}

type CreateSeatResponseSuccess struct {
	Message    string                    `json:"message"`
	Request    CreateSeatRequest         `json:"request"`
	Registered db.SeatModel              `json:"registered"`
	Generated  []db.SeatReservationModel `json:"seat_reservation"`
}

type CreateSeatResponseFailure struct {
	Message string            `json:"message"`
	Request CreateSeatRequest `json:"request"`
}

func CreateSeat(w http.ResponseWriter, r *http.Request) {
	var (
		status    int    = http.StatusNotImplemented
		message   string = "メッセージがありません"
		req       CreateSeatRequest
		seat      db.SeatModel
		generated []db.SeatReservationModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(CreateSeatResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = seat
			res.Generated = generated

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(CreateSeatResponseFailure)

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

		fmt.Printf("[Create Seat][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.Product == nil {
		status = http.StatusBadRequest
		message = "product_id is null"
		return
	}
	if req.Row == nil {
		status = http.StatusBadRequest
		message = "row is null"
		return
	}
	if req.Column == nil {
		status = http.StatusBadRequest
		message = "column is null"
		return
	}
	if req.Enable == nil {
		status = http.StatusBadRequest
		message = "is_enable is null"
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

	// 座席情報の挿入
	created, err := client.Seat.CreateOne(
		db.Seat.Row.SetIfPresent(req.Row),
		db.Seat.Column.SetIfPresent(req.Column),
		db.Seat.Product.Link(
			db.Product.ID.EqualsIfPresent(req.Product),
		),
		db.Seat.IsEnable.SetIfPresent(req.Enable),
		db.Seat.Remark.SetIfPresent(req.Remark),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Seatテーブル挿入エラー : ", err)
		return
	}

	seat = *created

	// 作成したSeatと紐づくStockを取得
	product, _ := client.Product.FindUnique(
		db.Product.ID.Equals(seat.ProductID),
	).With(
		db.Product.Price.Fetch().With(
			db.Price.Stock.Fetch(),
		),
	).Exec(ctx)

	for _, p := range product.RelationsProduct.Price {
		for _, s := range p.RelationsPrice.Stock {

			// SeatReservationにデータがない場合は生成
			g, err := client.SeatReservation.UpsertOne(
				db.SeatReservation.StockIDSeatID(
					db.SeatReservation.StockID.Equals(s.ID),
					db.SeatReservation.SeatID.Equals(seat.ID),
				),
			).Create(
				db.SeatReservation.Stock.Link(
					db.Stock.ID.Equals(s.ID),
				),
				db.SeatReservation.Seat.Link(
					db.Seat.ID.Equals(seat.ID),
				),
				db.SeatReservation.IsReserved.Set(false),
			).Update().Exec(ctx)
			if err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("seatreservationアップサートエラー : ", err)
				return
			}

			generated = append(generated, *g)
		}
	}

	status = http.StatusOK
	message = "正常終了"
}
