package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"Go_APIServer/db"
)

// 注文処理後のレスポンス用
type PostResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

// 編集処理後のレスポンス用
type EditResponseBody struct {
	Message  string    `json:"message"`
	EditInfo *EditInfo `json:"edit_info"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func APIServer() error {
	post_cnt := 0 // POSTのカウント用
	get_cnt := 0  // GETのカウント用
	change_cnt := 0
	mpost_cnt := 0

	fmt.Println("Server started.")

	// POST--------------------------------------------------------------------------------------
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		var order *Order
		var res PostResponseBody
		var status int
		var message string
		post_cnt++

		fmt.Printf("*** Post No.%d ***\n", post_cnt)

		defer func() {

			// レスポンスボディの作成
			if order == nil {
				res = PostResponseBody{
					Status:  status,
					Message: message,
				}
			} else {
				res = PostResponseBody{
					Status:  status,
					Message: message,
					Order:   *order,
				}
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				status = http.StatusInternalServerError
				message = fmt.Sprint("レスポンスの作成エラー :", err)
			}

			// 処理結果メッセージの表示（サーバ側）
			if status == 0 || message == "" {
				fmt.Println("ステータスコードまたはメッセージがありません")
			} else {
				fmt.Printf("[%d] %s\n", status, message)
			}

			fmt.Printf("*** Post No.%d End ***\n", post_cnt)

		}()

		// 注文情報をデコード
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("POSTデコードエラー :", err)
			return
		}

		// 注文処理
		status, message = order.Process()

	})

	// GET-------------------------------------------------------------------------------------
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		var status int
		var message string
		get_cnt++

		fmt.Printf("* Get No.%d *\n", get_cnt)

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("クライアント接続エラー :", err)
			return
		}

		ctx := context.Background()

		// Stockテーブルの内容を一括取得
		stock, err := client.Stock.FindMany().Exec(ctx)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("在庫テーブル取得エラー :", err)
			return
		}

		// インデントしてJsonに変換
		stock_json, err := json.MarshalIndent(stock, "", "  ")
		// インデントなし
		// stock_json, err := json.Marshal(stock)
		if err != nil {
			status = http.StatusBadRequest
			message = fmt.Sprint("JSON変換エラー :", err)
			return
		}

		status = http.StatusOK
		message = "正常終了"

		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				fmt.Println("クライアント切断エラー")
				panic(err)
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			w.Write(stock_json)

			// 処理結果メッセージの表示（サーバ側）
			if status == 0 || message == "" {
				fmt.Println("ステータスコードまたはメッセージがありません")
			} else {
				fmt.Printf("[%d] %s\n", status, message)
			}

			fmt.Printf("* Get No.%d End *\n", get_cnt)
		}()
	})

	// Order変更（予約の終了、キャンセル処理）
	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		var order *Order
		var res *Response
		change_cnt++

		fmt.Printf("*** Change No.%d ***\n", change_cnt)

		// 注文情報をデコード
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("POSTデコードエラー :", err)
			return
		}

		fmt.Println("変更情報 :", order)

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("クライアント接続エラー :", err)
			return
		}
		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				panic(fmt.Sprint("クライアント切断エラー :", err))
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(res.Status)

			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				res.Status = http.StatusInternalServerError
				res.Message = fmt.Sprint("レスポンスの作成エラー :", err)
			}

			// 処理結果メッセージの表示（サーバ側）
			if res.Status == 0 || res.Message == "" {
				fmt.Println("ステータスコードまたはメッセージがありません")
			} else {
				fmt.Printf("[%d] %s\n", res.Status, res.Message)
			}

			fmt.Printf("*** Change No.%d End ***\n", change_cnt)

		}()

		ctx := context.Background()

		// OrderIDから注文情報を取得
		order_info, err := client.Order.FindUnique(
			db.Order.ID.Equals(order.ID),
		).Exec(ctx)
		if err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("注文テーブル取得エラー :", err)
			return
		}

		// StockIDから在庫情報を取得
		stock_info, err := client.Stock.FindUnique(
			db.Stock.ID.Equals(order_info.InnerOrder.Product),
		).Exec(ctx)
		if err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("在庫テーブル取得エラー :", err)
			return
		}

		// 変更処理 [2]:予約完了, [3]:予約キャンセル
		switch order.State {
		case 2:
		case 3:
		default:
			res.Status = http.StatusBadRequest
			res.Message = "不正なステータス"
			return
		}

		// Orderテーブルのステータスを変更
		_, err = client.Order.FindUnique(
			db.Order.ID.Equals(order.ID),
		).Update(
			db.Order.State.Set(order.State),
		).Exec(ctx)
		if err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("注文テーブルアップデートエラー :", err)
			return
		}

		// 在庫を元に戻す
		id := order_info.InnerOrder.Product
		num := stock_info.InnerStock.Num + order_info.InnerOrder.Num

		_, err = client.Stock.FindUnique(
			db.Stock.ID.Equals(id),
		).Update(
			db.Stock.Num.Set(num),
		).Exec(ctx)
		if err != nil {
			res.Status = http.StatusBadRequest
			res.Message = fmt.Sprint("在庫テーブルアップデートエラー :", err)
		}

		// 正常終了時
		res.Status = http.StatusOK
		res.Message = "正常終了"

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

		// mapを各テーブル用の構造体に変換するため、一度jsonに変換
		info_json, err := json.Marshal(info.Info)
		if err != nil {
			message = fmt.Sprint("infoエンコードエラー :", err)
			status = http.StatusBadRequest
			return
		}

		// 各テーブルごとに処理を分岐
		if info.Table == "stock" {
			var stock *Stock

			// 変換したjsonをStockに変換
			if err := json.Unmarshal(info_json, &stock); err != nil {
				message = fmt.Sprint("infoデコードエラー :", err)
				status = http.StatusBadRequest
				return
			}

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
				fmt.Printf("[%d] %s\n", status, message)
			}

			fmt.Printf("### Manage Post No.%d END ###\n", mpost_cnt)
		}()
	})

	http.ListenAndServe(":8080", nil)

	return nil
}
