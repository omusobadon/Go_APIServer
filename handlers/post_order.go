// 注文情報のPOST
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	time_free_enable bool = true
	Seat_enable      bool = false
)

// リクエストを変換する構造体
type PostOrderRequestBody struct {
	Customer int       `json:"customer_id"`
	Stock    int       `json:"stock_id"`
	Seat     int       `json:"seat_id"`
	Qty      int       `json:"qty"`
	Start    time.Time `json:"start_at"`
	End      time.Time `json:"end_at"`
	People   int       `json:"number_people"`
	Remark   string    `json:"remark"`
}

// レスポンスに変換する構造体
type PostOrderResponseBody struct {
	Message string                `json:"message"`
	Order   db.OrderModel         `json:"order"`
	Detail  []db.OrderDetailModel `json:"detail"`
}

type TestRes struct {
	Message string        `json:"message"`
	Stock   db.StockModel `json:"stock"`
}

var post_order_cnt int // PostOrderのカウント用

func PostOrder(w http.ResponseWriter, r *http.Request) {
	post_order_cnt++
	var (
		req PostOrderRequestBody
		// order   Order
		status  int
		message string

		stock db.StockModel
	)

	fmt.Printf("*** Post Order No.%d ***\n", post_order_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		// res := PostOrderResponseBody{
		// 	Message: message,
		// 	Order:   order,
		// }

		res := TestRes{
			Message: message,
			Stock:   stock,
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

		fmt.Printf("*** Post Order No.%d End ***\n", post_order_cnt)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	fmt.Printf("注文情報 : %+v\n", req)

	// リクエストをOrderテーブルにコピー
	// order.Customer = req.Customer
	// order.Product = req.Product
	// order.Start = req.Start
	// order.End = req.End
	// order.Num = req.Num
	// order.Note = req.Note

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

	if time_free_enable {
		// 予約開始時刻が終了時刻より後でないかチェック
		if i := req.End.Sub(req.Start); i <= 0 {
			status = http.StatusBadRequest
			message = "予約時刻が不正"
			return
		}

		// 注文情報の在庫IDと一致する情報を取得
		stock, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(req.Stock)).With(
			db.Stock.Price.Fetch().With(
				db.Price.Product.Fetch().With(
					db.Product.Group.Fetch(),
				),
			),
		).Exec(ctx)
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

		fmt.Println("有効")

		// // 在庫が注文数を上回っているかチェック

		// if qty, _ := stock.Qty(); qty >= order.Num {

		// 	// 在庫テーブルに注文情報を反映
		// 	_, err := client.Stock.FindUnique(
		// 		db.Stock.ID.Equals(order.Product),
		// 	).Update(
		// 		db.Stock.Num.Set(stock.InnerStock.Num - order.Num),
		// 	).Exec(ctx)
		// 	if err != nil {
		// 		status = http.StatusBadRequest
		// 		message = fmt.Sprint("在庫テーブルアップデートエラー : ", err)
		// 		return
		// 	}

		// 	// Order.Stateを注文受付状態[1]に変更
		// 	order.State = 1

		// 	// 注文テーブルに注文情報をインサート
		// 	if err := order.Insert(client); err != nil {

		// 		// 注文を登録できなかった場合に在庫の数量を戻す
		// 		_, err := client.Stock.FindMany(
		// 			db.Stock.ID.Equals(order.Product),
		// 		).Update(
		// 			db.Stock.Num.Set(stock.InnerStock.Num + order.Num),
		// 		).Exec(ctx)
		// 		if err != nil {
		// 			panic(fmt.Sprint("在庫整合性エラー : ", err))
		// 		}

		// 		status = http.StatusBadRequest
		// 		message = fmt.Sprint("注文テーブルインサートエラー : ", err)
		// 		return
		// 	}

		// 	// 正常終了のとき
		// 	// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
		// 	order_info, err := client.Order.FindFirst(
		// 		db.Order.CustomerID.Equals(order.Customer),
		// 		db.Order.Time.Equals(order.Time),
		// 	).Exec(ctx)
		// 	if err != nil {
		// 		status = http.StatusInternalServerError
		// 		message = fmt.Sprint("注文ID取得エラー : ", err)
		// 		return
		// 	}

		// 	order.ID = order_info.ID
		// 	fmt.Printf("注文終了 : %+v\n", order)

		// 	status = http.StatusOK
		// 	message = "正常終了"

		// } else {
		// 	// 在庫不足のとき
		// 	status = http.StatusBadRequest
		// 	message = "在庫不足"
		// }

	}

	status = http.StatusOK
	message = "正常終了"

}
