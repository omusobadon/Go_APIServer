// 商品情報のGET
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
type GetProductResponseBody struct {
	Message string
	Length  int
	Product []db.ProductModel
}

var get_product_cnt int // GetProductの呼び出しカウント

func GetProduct(w http.ResponseWriter, r *http.Request) {
	var product []db.ProductModel
	var status int
	var message string
	get_product_cnt++

	fmt.Printf("* Get Product No.%d *\n", get_product_cnt)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetProductResponseBody{
			Message: message,
			Length:  len(product),
			Product: product,
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

		fmt.Printf("* Get Product No.%d End *\n", get_product_cnt)
	}()

	// リクエストパラメータの処理
	group_id, err := strconv.Atoi(r.FormValue("group"))
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("不正なパラメータ : ", err)
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

	// groupパラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if group_id == 0 {
		product, err = client.Product.FindMany().Exec(ctx)

	} else {
		product, err = client.Product.FindMany(
			db.Product.GroupID.Equals(group_id),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(product) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
