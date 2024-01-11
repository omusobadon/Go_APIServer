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
	// user_end_enable bool = false
	// user_notification bool = false

	// timezone =

	// 商品・在庫テーブルが空の場合、自動生成するかどうか
	auto_insert bool = true
)

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Header,Authorization,Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	http.HandleFunc("/get_shop", CORSMiddleware(handlers.GetShop))
	http.HandleFunc("/get_group", CORSMiddleware(handlers.GetGroup))
	http.HandleFunc("/get_product", CORSMiddleware(handlers.GetProduct))
	http.HandleFunc("/get_price", CORSMiddleware(handlers.GetPrice))
	http.HandleFunc("/get_seat", CORSMiddleware(handlers.GetSeat))
	http.HandleFunc("/get_stock", CORSMiddleware(handlers.GetStock))
	http.HandleFunc("/get_manage", CORSMiddleware(handlers.GetManage))

	// POST
	http.HandleFunc("/post_order", CORSMiddleware(handlers.PostOrder))

	// http.HandleFunc("/order_change", handlers.OrderChange)
	// http.HandleFunc("/manage_get", handlers.ManageGet)
	// http.HandleFunc("/manage_post", handlers.ManagePost)
	http.HandleFunc("/test", handlers.Test)

	// サーバの起動
	server.ListenAndServe()

	return nil
}
