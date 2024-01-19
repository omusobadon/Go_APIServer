package ini

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	seat_is bool = false // 座席を生成するか

	// 自動生成で使用される名前
	shop_name    string = "Car Shop"
	group_name   string = "Car"
	product_name string = "Color"
	price_name   string = "通常価格"
	seat_name    string = "test"
	stock_name   string = "time"

	// 自動生成する数。テーブルに1行でも情報がある場合は生成されない。
	shop_num    int = 1 // 店舗
	group_num   int = 1 // 商品グループ
	product_num int = 3 // 商品
	price_num   int = 1 // 価格
	seat_row    int = 3 // 座席（行）
	seat_column int = 5 // 座席（列）
	stock_num   int = 3 // 在庫
)

func AutoInsert() error {
	var (
		err     error
		shop    []db.ShopModel
		group   []db.ProductGroupModel
		product []db.ProductModel
		price   []db.PriceModel
		seat    []db.SeatModel
		stock   []db.StockModel
	)

	fmt.Println("[Auto Insert]")

	// 乱数シードの作成
	rand.NewSource(time.Now().UnixNano())

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

	// 店舗テーブルの取得
	shop, err = client.Shop.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ShopGetErr : %w", err)
	}

	// 店舗テーブルに情報がない場合は自動生成
	if len(shop) == 0 {
		fmt.Printf("店舗テーブルにインサート(%d件)...", shop_num)

		for i := 0; i < shop_num; i++ {

			// 電話番号の生成
			var n [5]int
			for j := range n {
				n[j] = rand.Intn(99)
			}

			phone := fmt.Sprintf("%02d%02d-%02d-%02d%02d", n[0], n[1], n[2], n[3], n[4])

			// insert
			_, err := client.Shop.CreateOne(
				db.Shop.Name.Set(fmt.Sprintf("%s %d", shop_name, i+1)),
				db.Shop.Mail.Set(fmt.Sprintf("shop%d@domain.jp", i+1)),
				db.Shop.Phone.Set(phone),
				db.Shop.Address.Set(fmt.Sprintf("Addr %d", i+1)),
			).Exec(ctx)
			if err != nil {
				fmt.Println("エラー")
				return fmt.Errorf("ShopInsertErr : %w", err)
			}
		}
		fmt.Println("完了")

		// 店舗テーブルの取得
		shop, err = client.Shop.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("ShopGetErr : %w", err)
		}
	}

	fmt.Printf("店舗が%d件見つかりました。\n", len(shop))

	// 商品グループテーブルの取得
	group, err = client.ProductGroup.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("GroupGetErr : %w", err)
	}

	if len(group) == 0 {
		fmt.Printf("商品グループテーブルにインサート(%d件)...", len(shop)*group_num)

		for _, i := range shop {
			for j := 0; j < group_num; j++ {

				// insert
				_, err := client.ProductGroup.CreateOne(
					db.ProductGroup.Name.Set(fmt.Sprintf("%s %d", group_name, j+1)),
					db.ProductGroup.Shop.Link(
						db.Shop.ID.Equals(i.ID),
					),
					// db.ProductGroup.StartBefore.Set(24),
					// db.ProductGroup.AvailableDuration.Set(72),
					// db.ProductGroup.UnitTime.Set(5),
					// db.ProductGroup.MaxTime.Set(72),
					// db.ProductGroup.Interval.Set(60),
				).Exec(ctx)
				if err != nil {
					fmt.Println("エラー")
					return fmt.Errorf("GroupInsertErr : %w", err)
				}
			}

		}
		fmt.Println("完了")

		// 商品グループテーブルの取得
		group, err = client.ProductGroup.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("GroupGetErr : %w", err)
		}
	}

	fmt.Printf("商品グループが%d件見つかりました。\n", len(group))

	// 商品テーブルの取得
	product, err = client.Product.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductGetErr : %w", err)
	}

	if len(product) == 0 {
		fmt.Printf("商品テーブルにインサート(%d件)...", len(group)*product_num)

		for _, i := range group {
			for j := 0; j < product_num; j++ {
				qty := (rand.Intn(10) + 1) * 10

				// insert
				_, err := client.Product.CreateOne(
					db.Product.Name.Set(fmt.Sprintf("%s %d", product_name, j+1)),
					db.Product.Group.Link(
						db.ProductGroup.ID.Equals(i.ID),
					),
					db.Product.Qty.Set(qty),
				).Exec(ctx)
				if err != nil {
					fmt.Println("エラー")
					return fmt.Errorf("ProductInsertErr : %w", err)
				}
			}

		}
		fmt.Println("完了")

		// 商品テーブルの取得
		product, err = client.Product.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("ProductGetErr : %w", err)
		}
	}

	fmt.Printf("商品が%d件見つかりました。\n", len(product))

	// 価格テーブルの取得
	price, err = client.Price.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("PriceGetErr : %w", err)
	}

	if len(price) == 0 {
		fmt.Printf("価格テーブルにインサート(%d件)...", len(product)*price_num)

		for _, i := range product {
			for j := 0; j < price_num; j++ {
				price := (rand.Intn(30) + 1) * 1000

				// insert
				_, err := client.Price.CreateOne(
					db.Price.Name.Set(fmt.Sprintf("%s %d", price_name, j+1)),
					db.Price.Value.Set(price),
					db.Price.Tax.Set(10),
					db.Price.Product.Link(
						db.Product.ID.Equals(i.ID),
					),
				).Exec(ctx)
				if err != nil {
					fmt.Println("エラー")
					return fmt.Errorf("PriceInsertErr : %w", err)
				}
			}

		}
		fmt.Println("完了")

		// 価格テーブルの取得
		price, err = client.Price.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("PriceGetErr : %w", err)
		}
	}

	fmt.Printf("価格が%d件見つかりました。\n", len(price))

	// 座席テーブルの取得
	seat, err = client.Seat.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("SeatGetErr : %w", err)
	}

	if len(seat) == 0 && seat_is {
		fmt.Printf("座席テーブルにインサート(%d件)...", len(product)*seat_row*seat_column)
		alphabets := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		for _, i := range product {
			for j := 0; j < seat_row; j++ {
				for k := 0; k < seat_column; k++ {

					// insert
					_, err := client.Seat.CreateOne(
						db.Seat.Row.Set(alphabets[j:j+1]),
						db.Seat.Column.Set(strconv.Itoa(k+1)),
						db.Seat.Product.Link(
							db.Product.ID.Equals(i.ID),
						),
						db.Seat.IsEnable.Set(true),
					).Exec(ctx)
					if err != nil {
						fmt.Println("エラー")
						return fmt.Errorf("SeatInsertErr : %w", err)
					}
				}
			}
		}
		fmt.Println("完了")

		// 座席テーブルの取得
		seat, err = client.Seat.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("SeatGetErr : %w", err)
		}
	}

	fmt.Printf("座席が%d件見つかりました。\n", len(seat))

	// 在庫テーブルの取得
	stock, err = client.Stock.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("StockGetErr : %w", err)
	}

	if len(stock) == 0 {
		fmt.Printf("在庫テーブルにインサート(%d件)...", len(price)*stock_num)

		// 現在時刻の取得
		now := time.Now()

		// 開始・終了時刻生成用の基準時間
		time_generated := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour()+1,
			0, 0, 0, time.Local,
		)

		for _, i := range price {
			p, err := client.Product.FindUnique(
				db.Product.ID.Equals(i.ProductID),
			).Exec(ctx)
			if err != nil {
				fmt.Println("エラー")
				return fmt.Errorf("")
			}

			qty, _ := p.Qty()

			for j := 0; j < stock_num; j++ {

				// 開始時刻の生成
				s := rand.Intn(10) + 1
				start := time_generated.Add(time.Duration(s) * time.Hour)

				// 終了時刻の生成
				e := s + rand.Intn(10) + 1
				end := time_generated.Add(time.Duration(e) * time.Hour)

				// insert
				_, err = client.Stock.CreateOne(
					db.Stock.Name.Set(fmt.Sprintf("%s %d", stock_name, j+1)),
					db.Stock.Price.Link(
						db.Price.ID.Equals(i.ID),
					),
					db.Stock.StartAt.Set(start),
					db.Stock.EndAt.Set(end),
					db.Stock.Qty.Set(qty),
					db.Stock.IsEnable.Set(true),
				).Exec(ctx)
				if err != nil {
					fmt.Println("エラー")
					return fmt.Errorf("StockInsertErr : %w", err)
				}
			}

		}
		fmt.Println("完了")

		// 在庫テーブルの取得
		stock, err = client.Stock.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("StockGetErr : %w", err)
		}
	}

	fmt.Printf("在庫が%d件見つかりました。\n", len(stock))

	return nil
}
