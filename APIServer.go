package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"Go_APIServer/db"
)

// 注文処理後のレスポンス用
type ResponseBody struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

// 編集処理後のレスポンス用
type EditResponseBody struct {
	Message  string    `json:"message"`
	EditInfo *EditInfo `json:"edit_info"`
}

func APIServer() error {
	post_cnt := 0 // POSTのカウント用
	get_cnt := 0  // GETのカウント用
	mpost_cnt := 0

	fmt.Println("Server started.")

	// POST--------------------------------------------------------------------------------------
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		var order *Order
		var status int     // HTTPステータスコード
		var message string // 処理結果メッセージ
		post_cnt++

		fmt.Printf("*** Post No.%d ***\n", post_cnt)

		// 注文情報をデコード
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			message = fmt.Sprint("POSTデコードエラー", err)
			status = http.StatusBadRequest
			return
		}

		// 注文時刻を取得
		order.Time = GetTime()

		fmt.Println("注文情報 :", order)

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			message = fmt.Sprint("クライアント接続エラー :", err)
			status = http.StatusBadRequest
			return
		}
		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				panic(err)
			}

			// レスポンスボディの作成
			res := ResponseBody{
				Message: message,
				Order:   *order,
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				message = fmt.Sprint("レスポンスの作成エラー :", err)
				status = http.StatusInternalServerError
				return
			}

			// 処理結果メッセージの表示（サーバ側）
			if status == 0 || message == "" {
				fmt.Println("ステータスコードまたはメッセージがありません")
			} else {
				fmt.Printf("[%d] %s\n", status, message)
			}

			fmt.Printf("*** Post No.%d End ***\n", post_cnt)
		}()

		ctx := context.Background()

		// 注文情報の商品idと一致する在庫情報を取得
		stock, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(order.Product),
		).Exec(ctx)
		if err != nil {
			message = fmt.Sprint("在庫テーブル取得エラー : ", err)
			status = http.StatusBadRequest
			return
		}

		// 在庫が注文数を上回っていたら注文処理を行う
		if stock.InnerStock.Num >= order.Num {
			// 在庫テーブルに注文情報を反映
			_, err := client.Stock.FindUnique(
				db.Stock.ID.Equals(order.Product),
			).Update(
				db.Stock.Num.Set(stock.InnerStock.Num - order.Num),
			).Exec(ctx)
			if err != nil {
				message = fmt.Sprint("在庫テーブルアップデートエラー :", err)
				status = http.StatusBadRequest
				return
			}

			// 注文テーブルに注文情報をインサート
			if err := order.Insert(client); err != nil {
				message = fmt.Sprint("注文テーブルインサートエラー :", err)
				status = http.StatusBadRequest

				// 注文を登録できなかった場合に在庫の数量を戻す
				_, err := client.Stock.FindMany(
					db.Stock.ID.Equals(order.Product),
				).Update(
					db.Stock.Num.Set(stock.InnerStock.Num + order.Num),
				).Exec(ctx)
				if err != nil {
					message = fmt.Sprint("在庫整合性エラー :", err)
					status = http.StatusInternalServerError
					return
				}
				return
			}

			// 正常終了のとき
			// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
			order_info, err := client.Order.FindFirst(
				db.Order.Customer.Equals(order.Customer),
				db.Order.Time.Equals(order.Time),
			).Exec(ctx)
			if err != nil {
				message = fmt.Sprint("注文情報取得エラー :", err)
				status = http.StatusInternalServerError
				return
			}

			order.ID = order_info.ID
			fmt.Println("注文受付 :", order)

			message = "正常終了"
			status = http.StatusOK

		} else {
			// 在庫不足のとき
			message = "在庫不足"
			status = http.StatusBadRequest
		}
	})

	// GET-------------------------------------------------------------------------------------
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		get_cnt++
		fmt.Printf("* Get No.%d *\n", get_cnt)

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			fmt.Println("クライアント接続エラー :", err)
			return
		}

		// 関数終了時に実行
		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				panic(err)
			}

			fmt.Printf("* Get No.%d End *\n", get_cnt)
		}()

		ctx := context.Background()

		// Stockテーブルの内容を一括取得
		stock, err := client.Stock.FindMany().Exec(ctx)
		if err != nil {
			fmt.Println("在庫テーブル取得エラー :", err)
			return
		}

		// インデントしてJsonに変換
		stock_json, err := json.MarshalIndent(stock, "", "  ")
		// インデントなし
		// stock_json, err := json.Marshal(stock)
		if err != nil {
			fmt.Println("JSON変換エラー :", err)
			return
		}

		fmt.Fprintln(w, string(stock_json))
		fmt.Println("正常終了")
	})

	// GET テスト用 すべてのテーブルを一覧表示-------------------------------------------------------
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			fmt.Println(err)
		}
		defer func() {
			if err := client.Prisma.Disconnect(); err != nil {
				panic(err)
			}
		}()

		ctx := context.Background()

		// Stockテーブルの内容を一括取得
		stock, err := client.Stock.FindMany().Exec(ctx)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, "| 商品id | 商品名 |  個数  |\n")
		for i := 0; i < len(stock); i++ {
			s := stock[i].InnerStock
			fmt.Fprintf(w, "|%8d|%8s|%8d|\n", s.ID, s.Name, s.Num)
		}
		fmt.Fprintln(w)

		// Orderテーブルの内容を一括取得
		order, err := client.Order.FindMany().Exec(ctx)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprint(w, "|注文番号|  客id  | 商品id | 注文数 | 注文時間\n")
		for i := 0; i < len(order); i++ {
			o := order[i].InnerOrder
			fmt.Fprintf(w, "|%8d|%8d|%8d|%8d|%+v\n", o.ID, o.Customer, o.ID, o.Num, o.Time)
		}
	})

	// テーブル編集 -----------------------------------------------------------------------------------
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		var info EditInfo
		var status int
		var message string
		// res := EditResponseBody{EditInfo: &info}
		mpost_cnt++

		fmt.Printf("### Manage Post No.%d ###\n", mpost_cnt)

		// 更新情報をデコード
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			message = fmt.Sprint("POSTデコードエラー :", err)
			status = http.StatusBadRequest
			return
		}

		// 更新時刻を取得
		info.Time = GetTime()

		fmt.Println("EditInfo :", info)

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			message = fmt.Sprint("クライアント接続エラー :", err)
			status = http.StatusBadRequest
			return
		}
		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				panic(err)
			}

			// レスポンスボディの作成
			res := EditResponseBody{
				Message:  message,
				EditInfo: &info,
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				message = fmt.Sprint("レスポンスの作成エラー :", err)
				status = http.StatusInternalServerError
			}

			// 処理結果メッセージの表示（サーバ側）
			if status == 0 || message == "" {
				fmt.Println("ステータスコードまたはメッセージがありません")
			} else {
				fmt.Printf("[%d] %s", status, message)
			}

			fmt.Printf("### Manage Post No.%d END ###\n", mpost_cnt)
		}()

		// mapを各テーブル用の構造体に変換するため、一度jsonに変換
		info_json, err := json.Marshal(info.Info)
		if err != nil {
			message = fmt.Sprint("infoエンコードエラー :", err)
			status = http.StatusBadRequest
			return
		}

		// 各テーブルごとに処理を分岐
		if info.Table == "stock" {
			var stock Stock

			// 変換したjsonをStockに変換
			if err := json.Unmarshal(info_json, &stock); err != nil {
				message = fmt.Sprint("infoデコードエラー :", err)
				status = http.StatusBadRequest
				return
			}

			fmt.Println("stock :", stock)

			// 編集タイプごとに処理を分岐
			// Type 1:Update, 2:Insert, 3:Delete
			if info.Type == 1 {
				if err := stock.Update(client); err != nil {
					message = fmt.Sprint("Stock Updateエラー :", err)
					status = http.StatusBadRequest
					return
				}

			} else if info.Type == 2 {
				if err := stock.Insert(client); err != nil {
					message = fmt.Sprint("Stock Insertエラー :", err)
					status = http.StatusBadRequest
					return
				}

			} else if info.Type == 3 {
				if err := stock.Delete(client); err != nil {
					message = fmt.Sprint("Stock Deleteエラー :", err)
					status = http.StatusBadRequest
					return
				}

			} else {
				message = "エラー : Type is not found"
				status = http.StatusBadRequest
				return
			}
		}

		// 処理が正常終了したらManageテーブルに登録
		if err := info.Insert(client); err != nil {
			message = fmt.Sprint("EditInfo Insertエラー :", err)
			status = http.StatusBadRequest
			return
		}

		status = 10
	})

	http.ListenAndServe(":8080", nil)

	return nil
}
