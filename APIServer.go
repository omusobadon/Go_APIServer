package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
)

// 商品・在庫テーブルが空の場合、自動生成するかどうか
const auto_insert = true

func APIServer() error {

	if auto_insert {
		// 商品・在庫テーブルが空の場合は自動生成するAutoInsert
		if err := handlers.AutoInsert(); err != nil {
			fmt.Println("自動インサートエラー : 処理を続行します")
		}
	}

	fmt.Println("Server started!")

	// 各ハンドラの呼び出し
	http.HandleFunc("/get", handlers.OrderGet)
	http.HandleFunc("/post", handlers.OrderPost)
	http.HandleFunc("/change", handlers.OrderChange)
	http.HandleFunc("/manage_get", handlers.ManageGet)
	http.HandleFunc("/manage_post", handlers.ManagePost)
	http.HandleFunc("/test", handlers.Test)

	// サーバの起動(TCPアドレス, http.Handler)
	http.ListenAndServe(":8080", nil)

	return nil
}
