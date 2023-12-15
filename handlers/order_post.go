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

// リクエストを変換する構造体
type PostRequestBody struct {
	Customer int       `json:"customer"`
	Product  int       `json:"product"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Num      int       `json:"num"`
	Note     string    `json:"note"`
}

// レスポンスに変換する構造体
type PostResponseBody struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

var post_cnt int // orderPOSTのカウント用

func OrderPost(w http.ResponseWriter, r *http.Request) {
	var req PostRequestBody
	var order Order
	var status int
	var message string
	post_cnt++

	fmt.Printf("*** Post No.%d ***\n", post_cnt)

	// リクエスト処理後のレスポンス処理
	defer func() {
		// レスポンスボディの作成
		res := PostResponseBody{
			Message: message,
			Order:   order,
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

		fmt.Printf("*** Post No.%d End ***\n", post_cnt)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	fmt.Printf("注文情報 : %+v\n", req)

	// リクエストをOrderテーブルにコピー
	order.Customer = req.Customer
	order.Product = req.Product
	order.Start = req.Start
	order.End = req.End
	order.Num = req.Num
	order.Note = req.Note

	// 注文時刻を取得
	order.Time = GetTime()

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

	// 注文情報の商品idと一致する在庫情報を取得
	stock, err := client.Stock.FindUnique(
		db.Stock.ID.Equals(order.Product),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("在庫テーブル取得エラー : ", err)
		return
	}

	// 予約開始時刻が終了時刻より後でないかチェック
	if i := order.End.Sub(order.Start); i <= 0 {
		status = http.StatusBadRequest
		message = "予約時刻が不正"
		return
	}

	// 在庫が注文数を上回っていたら注文処理を行う
	if stock.InnerStock.Num >= order.Num {
		// 在庫テーブルに注文情報を反映
		_, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(order.Product),
		).Update(
			db.Stock.Num.Set(stock.InnerStock.Num - order.Num),
		).Exec(ctx)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("在庫テーブルアップデートエラー : ", err)
			return
		}

		// Order.Stateを注文受付状態[1]に変更
		order.State = 1

		// 注文テーブルに注文情報をインサート
		if err := order.Insert(client); err != nil {

			// 注文を登録できなかった場合に在庫の数量を戻す
			_, err := client.Stock.FindMany(
				db.Stock.ID.Equals(order.Product),
			).Update(
				db.Stock.Num.Set(stock.InnerStock.Num + order.Num),
			).Exec(ctx)
			if err != nil {
				panic(fmt.Sprint("在庫整合性エラー : ", err))
			}

			status = http.StatusBadRequest
			message = fmt.Sprint("注文テーブルインサートエラー : ", err)
			return
		}

		// 正常終了のとき
		// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
		order_info, err := client.Order.FindFirst(
			db.Order.CustomerID.Equals(order.Customer),
			db.Order.Time.Equals(order.Time),
		).Exec(ctx)
		if err != nil {
			status = http.StatusInternalServerError
			message = fmt.Sprint("注文ID取得エラー : ", err)
			return
		}

		order.ID = order_info.ID
		fmt.Printf("注文終了 : %+v\n", order)

		status = http.StatusOK
		message = "正常終了"

	} else {
		// 在庫不足のとき
		status = http.StatusBadRequest
		message = "在庫不足"
	}
}
