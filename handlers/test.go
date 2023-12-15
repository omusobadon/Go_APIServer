package handlers

import (
	"Go_APIServer/db"
	"fmt"
	"net/http"
	"time"
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

	fmt.Println(time.Now())

	fmt.Printf("! Test No.%d End !\n", cnt)
}
