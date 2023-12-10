package handlers

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		// status = http.StatusBadRequest
		// message = fmt.Sprint("クライアント接続エラー : ", err)
		return
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー : ", err))
		}
	}()

	ctx := context.Background()

	product, err := client.Product.FindMany(
		db.Product.Stock.Some(
			db.Stock.ID.Equals(1),
		),
	).Exec(ctx)
	if err != nil {
		// status = http.StatusBadRequest
		// message = fmt.Sprint("商品テーブル取得エラー : ", err)
		return
	}

	p := product[0]

	fmt.Println(p.Stock())
}
