package handlers

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"math/rand"
	"time"
)

// 自動生成で作成される商品タイプ
const product_type string = "car"

// 自動生成する数。在庫テーブルに1行でも情報がある場合は生成されない。
const gen_num int = 10

// 商品・在庫テーブルが空か判定し、空の場合は自動生成
func AutoInsert() error {

	// データベース接続用クライアントの作成
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return fmt.Errorf("クライアント接続エラー : %w", err)
	}
	defer func() {
		// クライアントの切断
		if err := client.Prisma.Disconnect(); err != nil {
			panic(fmt.Sprint("クライアント切断エラー : ", err))
		}
	}()

	ctx := context.Background()

	// 在庫テーブルの取得
	stock, err := client.Stock.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("在庫テーブル取得エラー : %w", err)
	}

	// 在庫テーブルが空のとき
	if len(stock) == 0 {
		var product []db.ProductModel
		fmt.Println("在庫がありません。商品テーブルから自動生成します。")

		// 乱数シードの作成
		rand.NewSource(time.Now().UnixNano())

		// 商品テーブルの取得
		product, err = client.Product.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("商品テーブル取得エラー : %w", err)
		}

		// 商品テーブルが空のとき
		if len(product) == 0 {
			fmt.Println("商品がありません。自動生成します。")

			for i := 0; i < gen_num; i++ {
				value := (rand.Intn(99) + 1) * 100
				num := (rand.Intn(9) + 1) * 10

				// Insert
				_, err := client.Product.CreateOne(
					db.Product.Category.Set(product_type),
					db.Product.Name.Set(fmt.Sprint(product_type, i+1)),
					db.Product.Value.Set(value),
					db.Product.Num.Set(num),
					db.Product.Note.Set(""),
				).Exec(ctx)
				if err != nil {
					return fmt.Errorf("商品テーブルインサートエラー : %w", err)
				}
			}

			// 商品テーブルの取得
			product, err = client.Product.FindMany().Exec(ctx)
			if err != nil {
				return fmt.Errorf("商品テーブル取得エラー : %w", err)
			}

		} else {
			fmt.Printf("商品が%d件見つかりました。在庫を生成します。\n", len(product))
		}

		// 現在時刻の取得
		now := GetTime()

		// 商品テーブルから在庫を生成
		for _, v := range product {

			// 開始時刻の生成
			s := rand.Intn(10) + 1
			start := now.Add(time.Duration(s) * time.Hour)

			// 終了時刻の生成
			e := s + rand.Intn(10) + 1
			end := now.Add(time.Duration(e) * time.Hour)

			// インターバルの生成
			interval := fmt.Sprintf("%dh", rand.Intn(10))

			// Insert
			_, err := client.Stock.CreateOne(
				db.Stock.Product.Set(v.ID),
				db.Stock.Start.Set(start),
				db.Stock.End.Set(end),
				db.Stock.Interval.Set(interval),
				db.Stock.Num.Set(v.Num),
				db.Stock.State.Set(true),
			).Exec(ctx)
			if err != nil {
				return fmt.Errorf("在庫テーブルインサートエラー : %w", err)
			}
		}

		fmt.Println("生成完了")

	} else {
		fmt.Printf("在庫が%d件見つかりました。処理を続行します。\n", len(stock))
	}

	return nil
}
