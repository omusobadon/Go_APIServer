package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
	"time"
)

const (
	// time_free_enable bool = true
	// Seat_enable      bool = false

	// timezone =

	// 商品・在庫テーブルが空の場合、自動生成するかどうか
	auto_insert bool = true
)

func APIServer() error {

	// サーバ定義
	var server = &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// テーブルに情報がない場合に自動インサート（テスト用）
	if auto_insert {
		if err := AutoInsert(); err != nil {
			fmt.Println(err)
			fmt.Println("処理を続行します")
		}
	}

	// 各テーブルのチェック
	// 未実装

	fmt.Println("Server started!")

	// 各ハンドラの呼び出し
	// GET
	http.HandleFunc("/get_shop", handlers.GetShop)
	http.HandleFunc("/get_group", handlers.GetGroup)
	http.HandleFunc("/get_product", handlers.GetProduct)
	http.HandleFunc("/get_price", handlers.GetPrice)
	http.HandleFunc("/get_seat", handlers.GetSeat)
	http.HandleFunc("/get_stock", handlers.GetStock)

	// POST
	http.HandleFunc("/post_order", handlers.PostOrder)

	// http.HandleFunc("/order_change", handlers.OrderChange)
	// http.HandleFunc("/manage_get", handlers.ManageGet)
	// http.HandleFunc("/manage_post", handlers.ManagePost)
	http.HandleFunc("/test", handlers.Test)

	// サーバの起動
	server.ListenAndServe()

	return nil
}
