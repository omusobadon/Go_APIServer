// 注文情報のPOST
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
// ユーザによる時刻指定を有効にするか
// 有効にするとStockの開始・終了時刻は無効となる
// time_free_enable bool = true
// seat_enable      bool = false // 座席指定を有効にするか
// payment_enable   bool = false // 決済処理をするか
// user_notification_enable bool = false // 予約時間の終了をにユーザへ通知するか

// 予約時間終了後にユーザに終了処理をさせるか
// ユーザが終了処理をするまで在庫は戻らない
// user_end_enable bool = false // 未実装
)

// リクエストを変換する構造体
type PostOrderRequest struct {
	Customer int                      `json:"customer_id"`
	Start    time.Time                `json:"start_at"`
	End      time.Time                `json:"end_at"`
	People   int                      `json:"number_people"`
	Remark   string                   `json:"remark"`
	Detail   []PostOrderRequestDetail `json:"detail"`
}

type PostOrderRequestDetail struct {
	Stock int `json:"stock_id"`
	Seat  int `json:"seat_id"`
	Qty   int `json:"qty"`
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

var post_order_cnt int // PostOrderのカウント用

func PostOrder(w http.ResponseWriter, r *http.Request) {
	post_order_cnt++
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     PostOrderRequest
		order   db.OrderModel
		detail  []db.OrderDetailModel
	)

	fmt.Printf("*** Post Order No.%d ***\n", post_order_cnt)

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
			res.Order = order
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

		fmt.Printf("[%d] %s\n", status, message)
		fmt.Printf("*** Post Order No.%d End ***\n", post_order_cnt)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
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

	// 顧客情報が存在するか確認
	_, err := client.Customer.FindUnique(
		db.Customer.ID.Equals(req.Customer),
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

	// 予約開始時刻が現在よりも前のとき
	if req.Start.Before(now) {
		status = http.StatusBadRequest
		message = "予約時間が過ぎています"
		return
	}

	// 予約開始時刻が終了時刻より後でないかチェック
	if req.Start.After(req.End) {
		status = http.StatusBadRequest
		message = "予約開始時刻が終了時刻よりも後です"
		return
	}

	// 複数の注文の処理
	for _, v := range req.Detail {

		// 注文情報の在庫IDと一致する情報を取得
		stock, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(v.Stock)).With(
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
			message = fmt.Sprint("無効な在庫 : ", err)
			return
		}

		// Productテーブルの予約開始可能期間(start_before)から予約開始可能時刻を計算
		start_dur := time.Duration(stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.StartBefore)
		start := now.Add(start_dur * time.Hour)

		// 予約開始時刻が予約開始可能時刻より前のとき
		if req.Start.Before(start) {
			status = http.StatusBadRequest
			message = "予約開始時刻までが短すぎます"
			return
		}

		// 予約可能期間の終了時刻を計算
		end_dur := time.Duration(stock.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Group.AvailableDuration)
		end := start.Add(end_dur * time.Hour)

		// 予約開始時刻が予約可能期間より後のとき
		if req.Start.After(end) {
			status = http.StatusBadRequest
			message = "まだ予約できません"
			return
		}

		// // 在庫が注文数よりも少ない場合はエラー
		var qty int
		if qty, _ = stock.Qty(); qty < v.Qty {
			status = http.StatusBadRequest
			message = "在庫不足"
			return
		}

		// 在庫テーブルに注文情報を反映
		_, err = client.Stock.FindUnique(
			db.Stock.ID.Equals(v.Stock),
		).Update(
			db.Stock.Qty.Set(qty - v.Qty),
		).Exec(ctx)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("在庫テーブルアップデートエラー : ", err)
			return
		}

	}

	// Orderテーブルに注文情報をインサート
	ord, err := client.Order.CreateOne(
		db.Order.IsAccepted.Set(true),
		db.Order.Customer.Link(
			db.Customer.ID.Equals(req.Customer),
		),
		db.Order.StartAt.Set(req.Start),
		db.Order.EndAt.Set(req.End),
		db.Order.NumberPeople.Set(req.People),
		db.Order.Remark.Set(req.Remark),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("注文テーブルインサートエラー : ", err)
		return
	}

	order = *ord

	// OrderDetailテーブルに注文情報をインサート
	for _, v := range req.Detail {
		d, err := client.OrderDetail.CreateOne(
			db.OrderDetail.Order.Link(
				db.Order.ID.Equals(order.ID),
			),
			db.OrderDetail.Stock.Link(
				db.Stock.ID.Equals(v.Stock),
			),
			db.OrderDetail.Qty.Set(v.Qty),
		).Exec(ctx)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("注文詳細テーブルインサートエラー : ", err)
			return
		}

		detail = append(detail, *d)
	}

	fmt.Println(order)
	fmt.Println(detail)

	status = http.StatusOK
	message = "正常終了"
}
