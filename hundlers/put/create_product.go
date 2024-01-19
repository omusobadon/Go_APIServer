package put

import (
	"Go_APIServer/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateProductRequest struct {
	Group  *int    `json:"group_id"`
	Name   *string `json:"name"`
	Max    *int    `json:"max_people"`
	Qty    *int    `json:"qty"`
	Remark *string `json:"remark"`
	Img    *string `json:"img_data"`
}

type CreateProductResponseSuccess struct {
	Message    string               `json:"message"`
	Request    CreateProductRequest `json:"request"`
	Registered db.ProductModel      `json:"registered"`
}

type CreateProductResponseFailure struct {
	Message string               `json:"message"`
	Request CreateProductRequest `json:"request"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var (
		status  int    = http.StatusNotImplemented
		message string = "メッセージがありません"
		req     CreateProductRequest
		product db.ProductModel
	)

	// 処理終了後のレスポンス処理
	defer func() {

		// レスポンスヘッダー
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		// 注文成功時
		if status == http.StatusOK {
			res := new(CreateProductResponseSuccess)

			// レスポンスボディの作成
			res.Message = message
			res.Request = req
			res.Registered = product

			// レスポンス構造体をJSONに変換して送信
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー : ", err)
			}

		} else {
			res := new(CreateProductResponseFailure)

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

		fmt.Printf("[Create Product][%d] %s\n", status, message)

	}()

	// 注文情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("POSTデコードエラー : ", err)
		return
	}

	// リクエストの中身が存在するか確認
	if req.Group == nil {
		status = http.StatusBadRequest
		message = "group_id is null"
		return
	}
	if req.Name == nil {
		status = http.StatusBadRequest
		message = "name is null"
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
	created, err := client.Product.CreateOne(
		db.Product.Name.SetIfPresent(req.Name),
		db.Product.Group.Link(
			db.ProductGroup.ID.EqualsIfPresent(req.Group),
		),
		db.Product.MaxPeople.SetIfPresent(req.Max),
		db.Product.Qty.SetIfPresent(req.Qty),
		db.Product.Remark.SetIfPresent(req.Remark),
		db.Product.ImgData.SetIfPresent(req.Img),
	).Exec(ctx)
	if err != nil {
		status = http.StatusBadRequest
		message = fmt.Sprint("Productテーブル挿入エラー : ", err)
		return
	}

	product = *created

	status = http.StatusOK
	message = "正常終了"
}
