package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateStockRequest struct {
	Price  *int       `json:"price_id"`
	Name   *string    `json:"name"`
	Qty    *int       `json:"qty"`
	Start  *time.Time `json:"start_at"`
	End    *time.Time `json:"end_at"`
	Enable *bool      `json:"is_enable"`
}

type CreateStockResponseSuccess struct {
	Message    string             `json:"message"`
	Request    CreateStockRequest `json:"request"`
	Registered db.StockModel      `json:"registered"`
}

type CreateStockResponseFailure struct {
	Message string             `json:"message"`
	Request CreateStockRequest `json:"request"`
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     CreateStockRequest
		stock   db.StockModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(CreateStockResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = stock

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(CreateStockResponseFailure)

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

		fmt.Printf("[Create Stock][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.Price == nil {
		status = http.StatusBadRequest
		message = "price_id is null"
		return
	}
	if req.Name == nil {
		status = http.StatusBadRequest
		message = "name is null"
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

	// 商品グループ情報の挿入
	created, err := client.Stock.CreateOne(
		db.Stock.Name.SetIfPresent(req.Name),
		db.Stock.Price.Link(
			db.Price.ID.EqualsIfPresent(req.Price),
		),
		db.Stock.Qty.SetIfPresent(req.Qty),
		db.Stock.StartAt.SetIfPresent(req.Start),
		db.Stock.EndAt.SetIfPresent(req.End),
		db.Stock.IsEnable.SetIfPresent(req.Enable),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Priceテーブル挿入エラー : ", err)
		return
	}

	stock = *created

	seat, err := client.Stock.FindUnique(
		db.Stock.ID.Equals(stock.ID),
	).With(
		db.Stock.Price.Fetch().With(
			db.Price.Product.Fetch().With(
				db.Product.Seat.Fetch(),
			),
		),
	).Exec(ctx)
	if err != nil {
		status = http.StatusInternalServerError
		message = fmt.Sprint("Seatテーブル取得エラー : ", err)
		return
	}

	fmt.Println(seat)

	status = http.StatusOK
	message = "正常終了"
}
