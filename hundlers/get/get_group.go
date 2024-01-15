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
type GetGroupResponseBody struct {
	Message string                 `json:"message"`
	Length  int                    `json:"length"`
	Group   []db.ProductGroupModel `json:"group"`
}

var get_group_cnt int // ShopGetの呼び出しカウント

func GetGroup(w http.ResponseWriter, r *http.Request) {
	get_group_cnt++
	var (
		group   []db.ProductGroupModel
		status  int
		message string
		err     error
	)

	// リクエスト処理後のレスポンス作成
	defer func() {
		// レスポンスボディの作成
		res := GetGroupResponseBody{
			Message: message,
			Length:  len(group),
			Group:   group,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		fmt.Printf("[Get Group.%d][%d] %s\n", get_group_cnt, status, message)
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

	// shopパラメータが"0"のときテーブルの内容を一括取得
	// "0"以外のときはパラメータで指定した情報を取得
	if shop_id == 0 {
		group, err = client.ProductGroup.FindMany().Exec(ctx)

	} else {
		group, err = client.ProductGroup.FindMany(
			db.ProductGroup.ShopID.Equals(shop_id),
		).Exec(ctx)
	}
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("商品グループテーブル取得エラー : ", err)
		return
	}

	// 取得した情報がないとき
	if len(group) == 0 {
		status = http.StatusBadRequest
		message = "商品グループ情報がありません"
		return
	}

	status = http.StatusOK
	message = "正常終了"
}
