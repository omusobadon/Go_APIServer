package main

import (
	"Go_APIServer/hundlers/delete"
	"Go_APIServer/hundlers/get"
	"Go_APIServer/hundlers/post"
	"Go_APIServer/hundlers/put"
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
	http.HandleFunc("/get_order", CORSMiddleware(get.GetOrder))
	http.HandleFunc("/get_manage", CORSMiddleware(get.GetManage))

	// POST
	http.HandleFunc("/post_order", CORSMiddleware(post.PostOrder))

	// PUT
	http.HandleFunc("/create_shop", CORSMiddleware(put.CreateShop))
	http.HandleFunc("/create_group", CORSMiddleware(put.CreateGroup))
	http.HandleFunc("/create_product", CORSMiddleware(put.CreateProduct))
	http.HandleFunc("/create_price", CORSMiddleware(put.CreatePrice))
	http.HandleFunc("/create_seat", CORSMiddleware(put.CreateSeat))
	http.HandleFunc("/create_stock", CORSMiddleware(put.CreateStock))
	http.HandleFunc("/create_customer", CORSMiddleware(put.CreateCustomer))

	http.HandleFunc("/update_shop", CORSMiddleware(put.UpdateShop))
	http.HandleFunc("/update_group", CORSMiddleware(put.UpdateGroup))
	http.HandleFunc("/update_product", CORSMiddleware(put.UpdateProduct))
	http.HandleFunc("/update_price", CORSMiddleware(put.UpdatePrice))
	http.HandleFunc("/update_seat", CORSMiddleware(put.UpdateSeat))
	http.HandleFunc("/update_stock", CORSMiddleware(put.UpdateStock))
	http.HandleFunc("/update_customer", CORSMiddleware(put.UpdateCustomer))
	http.HandleFunc("/update_order", CORSMiddleware(put.UpdateOrder))
	http.HandleFunc("/update_order_detail", CORSMiddleware(put.UpdateOrderDetail))

	// DELETE
	http.HandleFunc("/delete_stock", CORSMiddleware(delete.DeleteStock))
	http.HandleFunc("/delete_customer", CORSMiddleware(delete.DeleteCustomer))
	http.HandleFunc("/delete_detail", CORSMiddleware(delete.DeleteDetail))

	// サーバの起動
	fmt.Println("[Start APIServer]")
	server.ListenAndServe()

	return nil
}
