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
type GetManageTestResponse struct {
	Message   string                 `json:"message"`
	Groups    []db.ProductGroupModel `json:"groups"`
	Customers []db.CustomerModel     `json:"customers"`
}

var get_manage_test_cnt int // PostOrderのカウント用

func GetManage(w http.ResponseWriter, r *http.Request) {
	get_manage_test_cnt++
	var (
		err       error
		status    int    = http.StatusNotImplemented
		message   string = "メッセージがありません"
		groups    []db.ProductGroupModel
		customers []db.CustomerModel
	)

	// 処理終了後のレスポンス処理
	defer func() {
		res := new(GetManageTestResponse)

		// レスポンスボディの作成
		res.Message = message
		res.Groups = groups
		res.Customers = customers

		// レスポンスヘッダーの作成
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// レスポンス構造体をJSONに変換して送信
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[GetManage.%d][%d] %s\n", get_manage_test_cnt, status, message)

	}()

	// リクエストパラメータの取得
	shop_str := r.FormValue("shop_id")

	// パラメータが空でない場合はIntに変換
	var shop_id int
	if shop_str != "" {
		shop_id, err = strconv.Atoi(shop_str)
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

	// GET
	groups, err = client.ProductGroup.FindMany(
		db.ProductGroup.ID.Equals(shop_id),
	).With(
		db.ProductGroup.Product.Fetch().With(
			db.Product.Price.Fetch().With(
				db.Price.Stock.Fetch().With(
					db.Stock.OrderDetail.Fetch(),
				),
			),
		).With(
			db.Product.Seat.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("テーブル取得エラー : ", err)
		return
	}

	customers, err = client.Customer.FindMany().With(
		db.Customer.Order.Fetch().With(
			db.Order.PaymentState.Fetch(),
		).With(
			db.Order.ReservationCancel.Fetch(),
		).With(
			db.Order.ReservationEnd.Fetch(),
		).With(
			db.Order.OrderDetail.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
