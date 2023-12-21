package handlers

import (
	"Go_APIServer/db"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 自動生成で作成される商品タイプ
const product_type string = "Car"

// 自動生成する数。テーブルに1行でも情報がある場合は生成されない。
const (
	shop_num    int = 3  // 店舗
	pgroup_num  int = 10 // 商品グループ
	product_num int = 3  // 商品
	price_num   int = 3  // 価格
	seat_num    int = 40 // 座席
	stock_num   int = 10 // 在庫
)

// 商品・在庫テーブルが空か判定し、空の場合は自動生成
func AutoInsert() error {
	var shop []db.ShopModel
	var pgroup []db.ProductGroupModel
	// var product []db.ProductModel
	// var price []db.PriceModel
	// var seat []db.SeatModel
	// var stock []db.StockModel

	// 商品タイプを小文字に変換
	ptype := strings.ToLower(product_type)

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
	shop, err := client.Shop.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ShopGetErr : %w", err)
	}

	// 店舗テーブルに情報がない場合は自動生成
	if len(shop) == 0 {
		fmt.Print("店舗テーブルにインサート...")

		for i := 1; i <= shop_num; i++ {

			// 電話番号の生成
			var n [5]int
			for j := range n {
				n[j] = rand.Intn(99)
			}

			phone := fmt.Sprintf("%02d%02d-%02d-%02d%02d", n[0], n[1], n[2], n[3], n[4])

			// insert
			_, err := client.Shop.CreateOne(
				db.Shop.Name.Set(fmt.Sprintf("%s Shop %d", product_type, i)),
				db.Shop.Mail.Set(fmt.Sprintf("%sshop%d@domain.jp", ptype, i)),
				db.Shop.Phone.Set(phone),
				db.Shop.Address.Set(fmt.Sprintf("Addr %d", i)),
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
	pgroup, err = client.ProductGroup.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("ProductGroupGetErr : %w", err)
	}

	if len(pgroup) == 0 {
		fmt.Print("商品グループテーブルにインサート...")

		for _, i := range shop {
			for j := 1; j <= pgroup_num; j++ {

				// insert
				_, err := client.ProductGroup.CreateOne(
					db.ProductGroup.Name.Set(fmt.Sprintf("%s %d", product_type, j)),
					db.ProductGroup.Shop.Link(
						db.Shop.ID.Equals(i.ID),
					),
					db.ProductGroup.StartBefore.Set("3d"),
					db.ProductGroup.AvailableDuration.Set("3d"),
					db.ProductGroup.UnitTime.Set("1h"),
					db.ProductGroup.MaxTime.Set("72h"),
					db.ProductGroup.Interval.Set("3h"),
				).Exec(ctx)
				if err != nil {
					fmt.Println("エラー")
					return fmt.Errorf("PGroupInsertErr : %w", err)
				}
			}

		}
		fmt.Println("完了")

		// 商品グループテーブルの取得
		pgroup, err = client.ProductGroup.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("ProductGroupGetErr : %w", err)
		}
	}

	fmt.Printf("商品グループが%d件見つかりました。\n", len(pgroup))

	return nil
}
