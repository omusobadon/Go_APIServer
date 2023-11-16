package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"Go_APIServer/db"
)

// レスポンス用のJSONデータを格納する構造体
type ResponseBody struct {
	Status int   `json:"status"`
	Order  Order `json:"order_info"`
}

func APIServer() error {
	post_cnt := 0 // POSTのカウント用
	get_cnt := 0  // GETのカウント用
	mpost_cnt := 0

	fmt.Println("Server started.")

	// POST--------------------------------------------------------------------------------------
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		var order Order
		var order_stat int
		post_cnt++

		fmt.Printf("*** Post No.%d ***\n", post_cnt)

		// 注文時刻を取得
		order.Time = GetTime()

		// データベース接続用クライアントの作成
		client := db.NewClient()
		if err := client.Prisma.Connect(); err != nil {
			fmt.Println("クライアント接続エラー :", err)
			order_stat = 30
			return
		}

		// 関数の終了時に実行
		defer func() {
			// クライアントの切断
			if err := client.Prisma.Disconnect(); err != nil {
				panic(err)
			}

			// レスポンスの作成
			res := ResponseBody{
				Status: order_stat,
				Order:  order,
			}

			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // 200 OK
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				fmt.Println("レスポンスの作成エラー :", err)
				return
			}

			// ステータスメッセージの表示（サーバ側）
			switch order_stat {
			case 10:
				fmt.Println("正常終了")
			case 20:
				fmt.Println("在庫不足")
			case 30:
			default:
				fmt.Println("未解決エラー")
			}

			fmt.Printf("*** Post No.%d End ***\n", post_cnt)
		}()

		ctx := context.Background()

		// POSTされたjsonを注文テーブルOrderにデコード
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			fmt.Println("POSTデコードエラー :", err)
			order_stat = 30
			return
		}

		fmt.Println("注文情報 :", order)

		// 注文情報の商品idと一致する在庫情報を取得
		stock, err := client.Stock.FindFirst(
			db.Stock.ProductID.Equals(order.Product_id),
		).Exec(ctx)
		if err != nil {
			fmt.Println("在庫テーブル取得エラー :", err)
			order_stat = 30
			return
		}

		// 在庫が注文数を上回っていたら注文処理を行う
		if stock.InnerStock.StockNum >= order.Num {
			// 在庫テーブルに注文情報を反映
			_, err := client.Stock.FindMany(
				db.Stock.ProductID.Equals(order.Product_id),
			).Update(
				db.Stock.StockNum.Set(stock.InnerStock.StockNum - order.Num),
			).Exec(ctx)
			if err != nil {
				fmt.Println("在庫テーブルアップデートエラー :", err)
				order_stat = 30
				return
			}

			// 注文テーブルに注文情報をインサート
			_, err = client.Order.CreateOne(
				db.Order.CustID.Set(order.Cust_id),
				db.Order.ProductID.Set(order.Product_id),
				db.Order.OrderNum.Set(order.Num),
				db.Order.OrderTime.Set(order.Time),
			).Exec(ctx)
			if err != nil {
				fmt.Println("注文テーブルインサートエラー :", err)
				order_stat = 30

				// 注文を登録できなかった場合に在庫の数量を戻す
				_, err := client.Stock.FindMany(
					db.Stock.ProductID.Equals(order.Product_id),
				).Update(
					db.Stock.StockNum.Set(stock.InnerStock.StockNum + order.Num),
				).Exec(ctx)
				if err != nil {
					fmt.Println("在庫整合性エラー :", err)
					order_stat = 100
					return
				}
				return
			}

			// 正常終了のとき
			// 顧客IDと時刻をもとにテーブルを検索して注文IDを取得
			order_info, err := client.Order.FindFirst(
				db.Order.CustID.Equals(order.Cust_id),
				db.Order.OrderTime.Equals(order.Time),
			).Exec(ctx)
			if err != nil {
				fmt.Println("注文情報取得エラー :", err)
				order_stat = 30
				return
			}

			order.Order_id = order_info.OrderID
			fmt.Println("注文番号 :", order.Order_id)

			order_stat = 10

		} else {
			// 在庫不足のとき
			order_stat = 20
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
			fmt.Fprintf(w, "|%8d|%8s|%8d|\n", s.ProductID, s.ProductName, s.StockNum)
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
			fmt.Fprintf(w, "|%8d|%8d|%8d|%8d|%+v\n", o.OrderID, o.CustID, o.ProductID, o.OrderNum, o.OrderTime)
		}
	})

	http.HandleFunc("/manage_post", func(w http.ResponseWriter, r *http.Request) {
		var info ManageInfo
		var res ManageRes
		mpost_cnt++

		fmt.Printf("### Manage Post No.%d ###\n", mpost_cnt)

		defer func() {
			// レスポンスをJSON形式で返す
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // 200 OK
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, "レスポンスの作成エラー", http.StatusInternalServerError)
				fmt.Println("レスポンスの作成エラー :", err)
				res.Status = 30
			}

			// ステータスメッセージの表示（サーバ側）
			switch res.Status {
			case 10:
				fmt.Println("正常終了")
			case 20:
				fmt.Println("在庫不足")
			case 30:
			default:
				fmt.Println("未解決エラー")
			}

			fmt.Printf("### Manage Post No.%d END ###\n", mpost_cnt)
		}()

		// 更新情報をデコード
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			fmt.Println("POSTデコードエラー :", err)
			res.Status = 30
			return
		}

		// 更新情報を受け取って処理を行い、レスポンス用のデータを返す
		res = DBManage(info)
	})

	http.ListenAndServe(":8080", nil)

	return nil
}
