// 注文情報のPOST
package post

import (
	"Go_APIServer/db"
	"Go_APIServer/ini"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// リクエストを変換する構造体
type PostOrderRequest struct {
	Customer *int                      `json:"customer_id"`
	Start    *time.Time                `json:"start_at"`
	End      *time.Time                `json:"end_at"`
	Remark   *string                   `json:"remark"`
	Detail   *[]PostOrderRequestDetail `json:"detail"`
}

type PostOrderRequestDetail struct {
	Stock  *int `json:"stock_id"`
	Seat   *int `json:"seat_id"`
	People *int `json:"number_people"`
	Qty    *int `json:"qty"`
}

// レスポンスに変換する構造体（処理成功）
type PostOrderResponseSuccess struct {
	Message string                `json:"message"`
	Request PostOrderRequest      `json:"request"`
	Order   db.OrderModel         `json:"order"`
	Detail  []db.OrderDetailModel `json:"order_detail"`
}

// レスポンスに変換する構造体（処理失敗）
type PostOrderResponseFailure struct {
	Message string           `json:"message"`
	Request PostOrderRequest `json:"request"`
}

var options ini.LoadedOptions = ini.Options
var post_order_cnt int // PostOrderのカウント用

func PostOrder(w http.ResponseWriter, r *http.Request) {
	// body := http.MaxBytesReader(w, r, 10)
	// fmt.Println(body)

	post_order_cnt++
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     PostOrderRequest
		order   *db.OrderModel
		detail  []db.OrderDetailModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(PostOrderResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Order = *order
			res.Detail = detail

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(PostOrderResponseFailure)

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

		fmt.Printf("[Post Order.%d][%d] %s\n", post_order_cnt, status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// 値があるかチェック
	if req.Customer == nil {
		status = http.StatusBadRequest
		message = "customer_id is null"
		return
	}
	if options.Time_free_enable {
		if req.Start == nil {
			status = http.StatusBadRequest
			message = "start_at is null"
			return
		}
		if req.End == nil {
			status = http.StatusBadRequest
			message = "end_at is null"
			return
		}
	}
	if req.Detail == nil {
		status = http.StatusBadRequest
		message = "detail is null"
		return
	}

	for _, d := range *req.Detail {
		if d.Stock == nil {
			status = http.StatusBadRequest
			message = "stock_id is null"
			return
		}

		if options.Seat_enable {
			if d.Seat == nil {
				status = http.StatusBadRequest
				message = "seat is null"
				return
			}

		} else {
			if d.Qty == nil {
				status = http.StatusBadRequest
				message = "qty is null"
				return
			}
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

	// 顧客情報が存在するか確認
	_, err := client.Customer.FindUnique(
		db.Customer.ID.Equals(*req.Customer),
	).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		status = http.StatusBadRequest
		message = "顧客情報がありません"
		return
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("顧客テーブル取得エラー : ", err)
		return
	}

	// 現在時刻の取得
	now := time.Now()

	// ユーザによる時刻指定が有効のとき
	// 時刻指定が正しいかチェック
	if options.Time_free_enable {

		// 予約開始時刻が現在よりも前のときはエラー
		if req.Start.Before(now) {
			status = http.StatusBadRequest
			message = "予約時間が過ぎています"
			return
		}

		// 予約開始時刻が終了時刻より後の場合はエラー
		if req.Start.After(*req.End) {
			status = http.StatusBadRequest
			message = "予約開始時刻が終了時刻よりも後です"
			return
		}
	}

	// トランザクションの開始
	// Stock用
	var transaction_stock []db.RawStockModel
	err = client.Prisma.QueryRaw("BEGIN").Exec(ctx, &transaction_stock)
	if err != nil {
		status = http.StatusInternalServerError
		message = fmt.Sprint("トランザクション開始エラー : ", err)
		return
	}

	// ReservationSeat用
	var transaction_seat []db.RawSeatReservationModel
	err = client.Prisma.QueryRaw("BEGIN").Exec(ctx, &transaction_seat)
	if err != nil {
		status = http.StatusInternalServerError
		message = fmt.Sprint("トランザクション開始エラー : ", err)
		return
	}

	// 注文処理
	for _, v := range *req.Detail {

		// 注文情報の在庫IDと一致する情報を取得
		stock, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(*v.Stock)).With(
			db.Stock.Price.Fetch().With(
				db.Price.Product.Fetch().With(
					db.Product.Group.Fetch(),
				),
			),
		).Exec(ctx)
		if errors.Is(err, db.ErrNotFound) {
			status = http.StatusBadRequest
			message = "在庫情報がありません"
			return
		}
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("在庫テーブル取得エラー : ", err)
			return
		}

		// 在庫状態が有効かチェック
		if !stock.IsEnable {
			status = http.StatusBadRequest
			message = "在庫状態が無効"
			return
		}

		// Productテーブルのstart_beforeとinvalid_durationを取得
		start_before_dur := time.Duration(stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.StartBefore)
		invalid_duration := time.Duration(stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.InvalidDuration)

		// ユーザによる時刻指定が有効の場合は、リクエストstart_atから予約可能か判断
		// 無効の場合は、Stockのstart_atから判断
		if options.Time_free_enable {

			// 予約可能期間を計算
			start := req.Start.Add(start_before_dur * time.Hour * -1)
			end := req.Start.Add(invalid_duration * time.Hour * -1)

			// fmt.Printf("start: %v, end: %v\n", start, end)

			// 現在時刻が予約開始可能時刻より前のとき
			if now.Before(start) {
				status = http.StatusBadRequest
				message = "まだ予約できません"
				return
			}

			// 現在時刻が予約可能期間より後のとき
			if now.After(end) {
				status = http.StatusBadRequest
				message = "予約受付時刻を過ぎました"
				return
			}

			// 終了時間 - 開始時間
			time_sub := req.End.Sub(*req.Start)

			// unit_time(時間単位)通りかチェック
			unit_time, ok := stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.UnitTime()
			if ok {
				if unit_time != 0 {
					div := int(time_sub.Minutes()) % unit_time

					if div != 0 {
						status = http.StatusBadRequest
						message = "指定された時間単位ではありません"
						return
					}
				}
			}

			// max_time(最大予約時間)以内かチェック
			max_time, ok := stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.MaxTime()
			if ok {
				if max_time != 0 {
					if time_sub > (time.Duration(max_time) * time.Hour) {
						status = http.StatusBadRequest
						message = "最大予約時間を超過しています"
						return
					}
				}
			}

		} else {

			// 予約可能期間を計算
			s, _ := stock.StartAt()
			start := s.Add(start_before_dur * time.Hour * -1)

			e, _ := stock.StartAt()
			end := e.Add(invalid_duration * time.Hour * -1)

			// fmt.Printf("start: %v, end: %v\n", start, end)

			// 現在時刻が予約開始可能時刻より前のとき
			if now.Before(start) {
				status = http.StatusBadRequest
				message = "まだ予約できません"
				return
			}

			// 現在時刻が予約可能期間より後のとき
			if now.After(end) {
				status = http.StatusBadRequest
				message = "予約受付時刻を過ぎました"
				return
			}
		}

		// Productテーブルのmax_people(最大人数)以内かチェック
		// max_people = 0 の場合は最大人数無指定としてチェックしない
		max_people, ok := stock.RelationsStock.Price.RelationsPrice.Product.MaxPeople()
		if ok {
			if max_people != 0 {
				if *v.People > max_people {
					status = http.StatusBadRequest
					message = "人数超過です"
					return
				}
			}
		}

		// 座席が有効のときは座席情報を参照
		// 無効のときは在庫の数量を参照
		if options.Seat_enable {

			// SeatReservationの取得
			seat_resev, err := client.SeatReservation.FindUnique(
				db.SeatReservation.StockIDSeatID(
					db.SeatReservation.StockID.Equals(*v.Stock),
					db.SeatReservation.SeatID.Equals(*v.Seat),
				),
			).With(
				db.SeatReservation.Seat.Fetch(),
			).Exec(ctx)
			if errors.Is(err, db.ErrNotFound) {
				status = http.StatusBadRequest
				message = "座席情報がありません"
				return
			}
			if err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("SeatReservation取得エラー : ", err)
				return
			}

			// 座席が無効の場合はエラー
			if !seat_resev.RelationsSeatReservation.Seat.IsEnable {
				status = http.StatusBadRequest
				message = "無効な座席"
				return
			}

			// 予約済みでない場合は予約
			if seat_resev.IsReserved {
				status = http.StatusBadRequest
				message = "予約済みです"
				return

			} else {
				_, err := client.SeatReservation.FindUnique(
					db.SeatReservation.StockIDSeatID(
						db.SeatReservation.StockID.Equals(*v.Stock),
						db.SeatReservation.SeatID.Equals(*v.Seat),
					),
				).Update(
					db.SeatReservation.IsReserved.Set(true),
				).Exec(ctx)
				if err != nil {
					status = http.StatusBadRequest
					message = fmt.Sprint("Seat挿入エラー : ", err)
					return
				}
			}

		} else {

			// // 在庫が注文数よりも少ない場合はエラー
			qty, _ := stock.Qty()
			if qty < *v.Qty {
				status = http.StatusBadRequest
				message = "在庫不足"
				return
			}

			// 在庫テーブルに注文情報を反映
			_, err = client.Stock.FindUnique(
				db.Stock.ID.Equals(*v.Stock),
			).Update(
				db.Stock.Qty.Set(qty - *v.Qty),
			).Exec(ctx)
			if err != nil {
				status = http.StatusBadRequest
				message = fmt.Sprint("在庫テーブルアップデートエラー : ", err)
				return
			}
		}
	}

	// トランザクションの終了
	_ = client.Prisma.QueryRaw("COMMIT").Exec(ctx, transaction_stock)
	_ = client.Prisma.QueryRaw("COMMIT").Exec(ctx, transaction_seat)

	// Orderテーブルに注文情報をインサート
	if options.Time_free_enable {
		order, err = client.Order.CreateOne(
			db.Order.IsAccepted.Set(true),
			db.Order.IsPending.Set(false),
			db.Order.Customer.Link(
				db.Customer.ID.Equals(*req.Customer),
			),
			db.Order.StartAt.Set(*req.Start),
			db.Order.EndAt.Set(*req.End),
			db.Order.Remark.SetIfPresent(req.Remark),
		).Exec(ctx)

	} else {
		order, err = client.Order.CreateOne(
			db.Order.IsAccepted.Set(true),
			db.Order.IsPending.Set(false),
			db.Order.Customer.Link(
				db.Customer.ID.Equals(*req.Customer),
			),
			db.Order.Remark.SetIfPresent(req.Remark),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("注文テーブルインサートエラー : ", err)
		return
	}

	// OrderDetailテーブルに注文情報をインサート
	for _, v := range *req.Detail {

		var d *db.OrderDetailModel
		if options.Seat_enable {
			d, err = client.OrderDetail.CreateOne(
				db.OrderDetail.Order.Link(
					db.Order.ID.Equals(order.ID),
				),
				db.OrderDetail.Stock.Link(
					db.Stock.ID.Equals(*v.Stock),
				),
				db.OrderDetail.Seat.Link(
					db.Seat.ID.Equals(*v.Seat),
				),
				db.OrderDetail.NumberPeople.SetIfPresent(v.People),
			).Exec(ctx)

		} else {
			d, err = client.OrderDetail.CreateOne(
				db.OrderDetail.Order.Link(
					db.Order.ID.Equals(order.ID),
				),
				db.OrderDetail.Stock.Link(
					db.Stock.ID.Equals(*v.Stock),
				),
				db.OrderDetail.NumberPeople.SetIfPresent(v.People),
				db.OrderDetail.Qty.Set(*v.Qty),
			).Exec(ctx)
		}
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("注文詳細テーブルインサートエラー : ", err)
			return
		}

		detail = append(detail, *d)
	}

	status = http.StatusOK
	message = "正常終了"
}
