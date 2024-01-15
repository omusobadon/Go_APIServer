package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateShopRequest struct {
	Name    *string `json:"name"`
	Mail    *string `json:"mail"`
	Phone   *string `json:"phone"`
	Address *string `json:"address"`
}

type CreateShopResponseSuccess struct {
	Message    string            `json:"message"`
	Request    CreateShopRequest `json:"request"`
	Registered db.ShopModel      `json:"registered"`
}

type CreateShopResponseFailure struct {
	Message string            `json:"message"`
	Request CreateShopRequest `json:"request"`
}

func CreateShop(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     CreateShopRequest
		shop    db.ShopModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(CreateShopResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = shop

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(CreateShopResponseFailure)

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

		fmt.Printf("[Create Shop][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.Name == nil {
		status = http.StatusBadRequest
		message = "name is null"
		return
	}
	if req.Mail == nil {
		status = http.StatusBadRequest
		message = "mail is null"
		return
	}
	if req.Phone == nil {
		status = http.StatusBadRequest
		message = "phone is null"
		return
	}
	if req.Address == nil {
		status = http.StatusBadRequest
		message = "address is null"
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

	// 店舗情報の挿入
	created, err := client.Shop.CreateOne(
		db.Shop.Name.SetIfPresent(req.Name),
		db.Shop.Mail.SetIfPresent(req.Mail),
		db.Shop.Phone.SetIfPresent(req.Phone),
		db.Shop.Address.SetIfPresent(req.Address),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Shopテーブル挿入エラー : ", err)
		return
	}

	shop = *created

	status = http.StatusOK
	message = "正常終了"
}
