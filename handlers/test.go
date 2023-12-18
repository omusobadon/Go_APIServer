package handlers

import (
	"Go_APIServer/db"
	"fmt"
	"net/http"
)

var cnt int

func Test(w http.ResponseWriter, r *http.Request) {
	cnt++

	fmt.Printf("! Test No.%d !\n", cnt)

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

	// ctx := context.Background()

	// test, err := client.Stock.FindMany().With(db.Stock.Product.Fetch().With(db.Product.Group.Fetch().With(db.ProductGroup.Shop.Fetch()))).With(db.Stock.Time.Fetch()).Exec(ctx)
	// if err != nil {
	// 	fmt.Println("エラー :", err)
	// }

	// fmt.Println(test)

	// // レスポンスをJSON形式で返す
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// if err := json.NewEncoder(w).Encode(test); err != nil {
	// 	http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
	// }

	fmt.Printf("! Test No.%d End !\n", cnt)
}
