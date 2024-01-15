package delete

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteStockRequest struct {
	ID *int `json:"id"`
}

type DeleteStockResponse struct {
	Message    string             `json:"message"`
	Request    DeleteStockRequest `json:"request"`
	Registered db.StockModel      `json:"registered"`
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     DeleteStockRequest
		stock   db.StockModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(DeleteStockResponse)

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
			res := new(DeleteStockResponse)

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

		fmt.Printf("[Delete Stock][%d] %s\n", status, message)

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

	// 顧客情報の削除
	deleted, err := client.Stock.FindUnique(
		db.Stock.ID.EqualsIfPresent(req.ID),
	).Delete().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Stockテーブル削除エラー : ", err)
		return
	}

	stock = *deleted

	status = http.StatusOK
	message = "正常終了"
}
