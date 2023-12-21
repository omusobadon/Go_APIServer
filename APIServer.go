package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
)

const (
	// timezone =

	// 商品・在庫テーブルが空の場合、自動生成するかどうか
	auto_insert bool = true
)

func APIServer() error {

	if auto_insert {
		// 商品・在庫テーブルが空の場合は自動生成するAutoInsert
		if err := handlers.AutoInsert(); err != nil {
			fmt.Println(err)
			fmt.Println("処理を続行します")
		}
	}

	// 各テーブルのチェック

	fmt.Println("Server started!")

	// 各ハンドラの呼び出し
	http.HandleFunc("/shop_get", handlers.ShopGet)
	http.HandleFunc("/stock_get", handlers.StockGet)
	http.HandleFunc("/seat_get", handlers.SeatGet)
	// http.HandleFunc("/order_post", handlers.OrderPost)
	// http.HandleFunc("/order_change", handlers.OrderChange)
	// http.HandleFunc("/manage_get", handlers.ManageGet)
	// http.HandleFunc("/manage_post", handlers.ManagePost)
	http.HandleFunc("/test", handlers.Test)

	// サーバの起動(TCPアドレス, http.Handler)
	http.ListenAndServe(":8080", nil)

	return nil
}
