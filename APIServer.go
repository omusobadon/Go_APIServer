package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
)

// 商品・在庫テーブルが空の場合、自動生成するかどうか
const auto_insert = true

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

	if auto_insert {
		// 商品・在庫テーブルが空の場合は自動生成するAutoInsert
		if err := handlers.AutoInsert(); err != nil {
			fmt.Println("自動インサートエラー :", err)
			fmt.Println("処理を続行します")
		}
	}

	fmt.Println("Server started!")

	// CORS対応
	http.HandleFunc("/get", CORSMiddleware(handlers.OrderGet))
	http.HandleFunc("/post", CORSMiddleware(handlers.OrderPost))
	http.HandleFunc("/change", CORSMiddleware(handlers.OrderChange))
	http.HandleFunc("/manage_get", CORSMiddleware(handlers.ManageGet))
	http.HandleFunc("/manage_post", CORSMiddleware(handlers.ManagePost))
	http.HandleFunc("/test", CORSMiddleware(handlers.Test))

	// サーバの起動(TCPアドレス, http.Handler)
	http.ListenAndServe(":8080", nil)

	return nil
}
