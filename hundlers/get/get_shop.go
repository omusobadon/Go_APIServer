// 店舗情報のGET
package get

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// レスポンスに変換する構造体
type GetShopResponseBody struct {
	Message string         `json:"message"`
	Length  int            `json:"length"`
	Shop    []db.ShopModel `json:"shop"`
}

var get_shop_cnt int // GetShopの呼び出しカウント

func GetShop(w http.ResponseWriter, r *http.Request) {
	get_shop_cnt++
	var (
		shop    []db.ShopModel
		status  int
		message string
	)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetShopResponseBody{
			Message: message,
			Length:  len(shop),
			Shop:    shop,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[Get Shop.%d][%d] %s\n", get_shop_cnt, status, message)
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
	shop, err := client.Shop.FindMany().Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("店舗テーブル取得エラー : ", err)
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
