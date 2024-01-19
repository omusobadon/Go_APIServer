package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreatePriceRequest struct {
	Product *int    `json:"product_id"`
	Name    *string `json:"name"`
	Value   *int    `json:"value"`
	Tax     *int    `json:"tax"`
	Remark  *string `json:"remark"`
}

type CreatePriceResponseSuccess struct {
	Message    string             `json:"message"`
	Request    CreatePriceRequest `json:"request"`
	Registered db.PriceModel      `json:"registered"`
}

type CreatePriceResponseFailure struct {
	Message string             `json:"message"`
	Request CreatePriceRequest `json:"request"`
}

func CreatePrice(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     CreatePriceRequest
		price   db.PriceModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(CreatePriceResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = price

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(CreatePriceResponseFailure)

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

		fmt.Printf("[Create Price][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.Product == nil {
		status = http.StatusBadRequest
		message = "product_id is null"
		return
	}
	if req.Name == nil {
		status = http.StatusBadRequest
		message = "name is null"
		return
	}
	if req.Value == nil {
		status = http.StatusBadRequest
		message = "value is null"
		return
	}
	if req.Tax == nil {
		status = http.StatusBadRequest
		message = "tax is null"
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
	created, err := client.Price.CreateOne(
		db.Price.Name.SetIfPresent(req.Name),
		db.Price.Value.SetIfPresent(req.Value),
		db.Price.Tax.SetIfPresent(req.Tax),
		db.Price.Product.Link(
			db.Product.ID.EqualsIfPresent(req.Product),
		),
		db.Price.Remark.SetIfPresent(req.Remark),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Priceテーブル挿入エラー : ", err)
		return
	}

	price = *created

	status = http.StatusOK
	message = "正常終了"
}
