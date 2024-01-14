package main

import (
	"Go_APIServer/get"
	"Go_APIServer/post"
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

	// fmt.Println("Server started!")

	// 各ハンドラの呼び出し
	// GET
	http.HandleFunc("/get_shop", CORSMiddleware(get.GetShop))
	http.HandleFunc("/get_group", CORSMiddleware(get.GetGroup))
	http.HandleFunc("/get_product", CORSMiddleware(get.GetProduct))
	http.HandleFunc("/get_price", CORSMiddleware(get.GetPrice))
	http.HandleFunc("/get_seat", CORSMiddleware(get.GetSeat))
	http.HandleFunc("/get_stock", CORSMiddleware(get.GetStock))

	// 管理用GET
	http.HandleFunc("/get_customer", CORSMiddleware(get.GetCustomer))
	http.HandleFunc("/get_manage", CORSMiddleware(get.GetManage))

	// POST
	http.HandleFunc("/post_order", CORSMiddleware(post.PostOrder))
	http.HandleFunc("/post_customer", CORSMiddleware(post.PostCustomer))

	// http.HandleFunc("/order_change", post.OrderChange)
	// http.HandleFunc("/manage_get", post.ManageGet)
	// http.HandleFunc("/manage_post", post.ManagePost)
	http.HandleFunc("/test", post.Test)

	// サーバの起動
	fmt.Println("[Server start]")
	server.ListenAndServe()

	return nil
}
