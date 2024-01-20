package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateGroupRequest struct {
	ID       *int    `json:"id"`
	Shop     *int    `json:"shop_id"`
	Name     *string `json:"name"`
	Start    *int    `json:"start_before"`
	Invalid  *int    `json:"invalid_duration"`
	Unit     *int    `json:"unit_time"`
	Max      *int    `json:"max_time"`
	Interval *int    `json:"interval"`
}

type UpdateGroupResponseSuccess struct {
	Message    string               `json:"message"`
	Request    UpdateGroupRequest   `json:"request"`
	Registered db.ProductGroupModel `json:"registered"`
}

type UpdateGroupResponseFailure struct {
	Message string             `json:"message"`
	Request UpdateGroupRequest `json:"request"`
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var (
		status     int    = http.StatusNotImplemented
		message    string = "メッセージがありません"
		req        UpdateGroupRequest
		registered db.ProductGroupModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(UpdateGroupResponseSuccess)

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
			res := new(UpdateGroupResponseFailure)

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

		fmt.Printf("[Update Group][%d] %s\n", status, message)

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
	created, err := client.ProductGroup.FindUnique(
		db.ProductGroup.ID.EqualsIfPresent(req.ID),
	).Update(
		db.ProductGroup.ShopID.SetIfPresent(req.Shop),
		db.ProductGroup.Name.SetIfPresent(req.Name),
		db.ProductGroup.StartBefore.SetIfPresent(req.Start),
		db.ProductGroup.InvalidDuration.SetIfPresent(req.Invalid),
		db.ProductGroup.UnitTime.SetIfPresent(req.Unit),
		db.ProductGroup.MaxTime.SetIfPresent(req.Max),
		db.ProductGroup.Interval.SetIfPresent(req.Interval),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("ProductGroupテーブル挿入エラー : ", err)
		return
	}

	registered = *created

	status = http.StatusOK
	message = "正常終了"
}
