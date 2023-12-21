// 座席情報のGET
package handlers

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type SeatGetResponseBody struct {
	Message string         `json:"message"`
	Seat    []db.SeatModel `json:"seat"`
}

var seat_get_cnt int // ShopGetの呼び出しカウント

func SeatGet(w http.ResponseWriter, r *http.Request) {
	var seat []db.SeatModel
	var status int
	var message string
	seat_get_cnt++

	fmt.Printf("* Seat Get No.%d *\n", seat_get_cnt)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := SeatGetResponseBody{
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

		fmt.Printf("* Seat Get No.%d End *\n", seat_get_cnt)
	}()

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

	// Shopテーブルの内容を一括取得
	seat, err := client.Seat.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("座席テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
