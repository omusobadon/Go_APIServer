// 座席情報のGET
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// レスポンスに変換する構造体
type GetSeatResponseBody struct {
	Message string         `json:"message"`
	Length  int            `json:"length"`
	Seat    []db.SeatModel `json:"seat"`
}

var get_seat_cnt int // ShopGetの呼び出しカウント

func GetSeat(w http.ResponseWriter, r *http.Request) {
	var seat []db.SeatModel
	var status int
	var message string
	get_seat_cnt++

	fmt.Printf("* Get Seat No.%d *\n", get_seat_cnt)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetSeatResponseBody{
			Message: message,
			Seat:    seat,
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

		fmt.Printf("* Get Seat No.%d End *\n", get_seat_cnt)
	}()

	// リクエストパラメータの処理
	product_id, err := strconv.Atoi(r.FormValue("product"))
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("不正なパラメータ : ", err)
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

	// productパラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if product_id == 0 {
		seat, err = client.Seat.FindMany().Exec(ctx)

	} else {
		seat, err = client.Seat.FindMany(
			db.Seat.ProductID.Equals(product_id),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(seat) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}