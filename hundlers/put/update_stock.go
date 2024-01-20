package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateStockRequest struct {
	ID     *int       `json:"id"`
	Price  *int       `json:"price_id"`
	Name   *string    `json:"name"`
	Qty    *int       `json:"qty"`
	Start  *time.Time `json:"start_at"`
	End    *time.Time `json:"end_at"`
	Enable *bool      `json:"is_enable"`
}

type UpdateStockResponseSuccess struct {
	Message    string             `json:"message"`
	Request    UpdateStockRequest `json:"request"`
	Registered db.StockModel      `json:"registered"`
}

type UpdateStockResponseFailure struct {
	Message string             `json:"message"`
	Request UpdateStockRequest `json:"request"`
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	var (
		status     int    = http.StatusNotImplemented
		message    string = "メッセージがありません"
		req        UpdateStockRequest
		registered db.StockModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(UpdateStockResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = registered

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(UpdateStockResponseFailure)

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

		fmt.Printf("[Update Stock][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.ID == nil {
		status = http.StatusBadRequest
		message = "id is null"
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

	// 顧客情報の挿入
	created, err := client.Stock.FindUnique(
		db.Stock.ID.EqualsIfPresent(req.ID),
	).Update(
		db.Stock.PriceID.SetIfPresent(req.Price),
		db.Stock.Name.SetIfPresent(req.Name),
		db.Stock.Qty.SetIfPresent(req.Qty),
		db.Stock.StartAt.SetIfPresent(req.Start),
		db.Stock.EndAt.SetIfPresent(req.End),
		db.Stock.IsEnable.SetIfPresent(req.Enable),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Stockテーブル挿入エラー : ", err)
		return
	}

	registered = *created

	status = http.StatusOK
	message = "正常終了"
}
