package handlers

import (
	"Go_APIServer/db"
	"encoding/json"
	"fmt"
	"net/http"
)

var cnt int

func Test(w http.ResponseWriter, r *http.Request) {
	var status int = 200
	var message string = "test"
	var test string
	cnt++

	fmt.Printf("! Test No.%d !\n", cnt)

	// リクエスト処理後のレスポンス作成
	defer func() {

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(test); err != nil {
			http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
			status = http.StatusInternalServerError
			message = fmt.Sprint("レスポンスの作成エラー : ", err)
		}

		// 処理結果メッセージの表示（サーバ側）
		if status == 0 || message == "" {
			fmt.Println("ステータスコードまたはメッセージがありません")
		} else {
			fmt.Printf("[%d] %s\n", status, message)
		}

		fmt.Printf("! Test No.%d End !\n", cnt)
	}()

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

	test = r.FormValue("test")
}
