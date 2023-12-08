package main

import (
	"Go_APIServer/handlers"
	"fmt"
	"net/http"
)

func APIServer() error {
	fmt.Println("Server started.")

	// 各ハンドラの呼び出し
	http.HandleFunc("/get", handlers.OrderGet)
	http.HandleFunc("/post", handlers.OrderPost)
	http.HandleFunc("/change", handlers.OrderChange)
	http.HandleFunc("/manage_get", handlers.ManageGet)
	http.HandleFunc("/manage_post", handlers.ManagePost)

	// サーバの起動(TCPアドレス, http.Handler)
	http.ListenAndServe(":8080", nil)

	return nil
}
