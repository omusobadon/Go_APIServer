// 管理情報のGET
package srcbk

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type GetManageResponseBody struct {
	Message string                 `json:"message"`
	Shop    []db.ShopModel         `json:"shop"`
	Group   []db.ProductGroupModel `json:"group_group"`
	Product []db.ProductModel      `json:"product"`
	Price   []db.PriceModel        `json:"price"`
	Seat    []db.SeatModel         `json:"seat"`
	Stock   []db.StockModel        `json:"stock"`

	Customer []db.CustomerModel          `json:"customer"`
	Order    []db.OrderModel             `json:"order"`
	Detail   []db.OrderDetailModel       `json:"order_detail"`
	Payment  []db.PaymentStateModel      `json:"payment_state"`
	Cancel   []db.ReservationCancelModel `json:"reservation_cancel"`
	End      []db.ReservationEndModel    `json:"reservation_end"`
}

var get_manage_cnt int // GetManageの呼び出しカウント

func GetManagebk(w http.ResponseWriter, r *http.Request) {
	get_manage_cnt++
	res := new(GetManageResponseBody)
	var (
		status  int
		message string
		err     error
	)

	fmt.Printf("* Get Manage No.%d *\n", get_manage_cnt)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res.Message = message

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

		fmt.Printf("* Get Manage No.%d End *\n", get_manage_cnt)
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
	res.Shop, err = client.Shop.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("店舗テーブル取得エラー : ", err)
		return
	}

	// ProductGroup
	res.Group, err = client.ProductGroup.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("グループテーブル取得エラー : ", err)
		return
	}

	// Product
	res.Product, err = client.Product.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品テーブル取得エラー : ", err)
		return
	}

	// Price
	res.Price, err = client.Price.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("価格テーブル取得エラー : ", err)
		return
	}

	// Seat
	res.Seat, err = client.Seat.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("座席テーブル取得エラー : ", err)
		return
	}

	// Stock
	res.Stock, err = client.Stock.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("在庫テーブル取得エラー : ", err)
		return
	}

	// Customer
	res.Customer, err = client.Customer.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("顧客テーブル取得エラー : ", err)
		return
	}

	// Order
	res.Order, err = client.Order.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("注文テーブル取得エラー : ", err)
		return
	}

	// OrderDetail
	res.Detail, err = client.OrderDetail.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("注文詳細テーブル取得エラー : ", err)
		return
	}

	// PaymentState
	res.Payment, err = client.PaymentState.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("決済状態テーブル取得エラー : ", err)
		return
	}

	// ReservationCancel
	res.Cancel, err = client.ReservationCancel.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("予約キャンセルテーブル取得エラー : ", err)
		return
	}

	// ReservationEnd
	res.End, err = client.ReservationEnd.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("予約終了テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
