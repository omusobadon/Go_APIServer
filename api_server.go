package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
	"time"
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

	fmt.Println("Server started!")

	// 各ハンドラの呼び出し
	// GET
	http.HandleFunc("/get_shop", CORSMiddleware(handlers.GetShop))
	http.HandleFunc("/get_group", CORSMiddleware(handlers.GetGroup))
	http.HandleFunc("/get_product", CORSMiddleware(handlers.GetProduct))
	http.HandleFunc("/get_price", CORSMiddleware(handlers.GetPrice))
	http.HandleFunc("/get_seat", CORSMiddleware(handlers.GetSeat))
	http.HandleFunc("/get_stock", CORSMiddleware(handlers.GetStock))

	// 管理用GET
	http.HandleFunc("/get_customer", CORSMiddleware(handlers.GetCustomer))
	http.HandleFunc("/get_manage", CORSMiddleware(handlers.GetManage))

	// POST
	http.HandleFunc("/post_order", CORSMiddleware(handlers.PostOrder))
	http.HandleFunc("/post_customer", CORSMiddleware(handlers.PostCustomer))

	// http.HandleFunc("/order_change", handlers.OrderChange)
	// http.HandleFunc("/manage_get", handlers.ManageGet)
	// http.HandleFunc("/manage_post", handlers.ManagePost)
	http.HandleFunc("/test", handlers.Test)

	// サーバの起動
	server.ListenAndServe()

	return nil
}
