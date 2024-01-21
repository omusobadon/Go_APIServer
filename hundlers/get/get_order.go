// 商品グループ情報のGET
package get

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// レスポンスに変換する構造体
type GetOrderResponseBody struct {
	Message string          `json:"message"`
	Length  int             `json:"length"`
	Order   []db.OrderModel `json:"order"`
}

var get_order_cnt int // ShopGetの呼び出しカウント

func GetOrder(w http.ResponseWriter, r *http.Request) {
	get_order_cnt++
	var (
		order   []db.OrderModel
		status  int
		message string
		err     error
	)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetOrderResponseBody{
			Message: message,
			Length:  len(order),
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

		fmt.Printf("[Get Order.%d][%d] %s\n", get_order_cnt, status, message)
	}()

	// リクエストパラメータの取得
	customer_str := r.FormValue("customer_id")

	// パラメータが空でない場合はIntに変換
	var customer_id int
	if customer_str != "" {
		customer_id, err = strconv.Atoi(customer_str)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("不正なパラメータ : ", err)
			return
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

	// パラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if customer_id == 0 {
		order, err = client.Order.FindMany().With(
			db.Order.PaymentState.Fetch(),
			db.Order.ReservationCancel.Fetch(),
			db.Order.ReservationEnd.Fetch(),
			db.Order.OrderDetail.Fetch(),
		).Exec(ctx)

	} else {
		order, err = client.Order.FindMany(
			db.Order.CustomerID.Equals(customer_id),
		).With(
			db.Order.PaymentState.Fetch(),
			db.Order.ReservationCancel.Fetch(),
			db.Order.ReservationEnd.Fetch(),
			db.Order.OrderDetail.Fetch(),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(order) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
