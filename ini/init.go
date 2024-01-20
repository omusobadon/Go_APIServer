package ini

import (
	"Go_APIServer/db"
	"context"
	"fmt"
)

// 商品・在庫テーブルが空の場合、自動生成するかどうか（テスト用）
const auto_insert bool = false

// SeatReservationの自動生成
func generateSeatReservation() error {

	fmt.Println("[SeatReservation generate]")

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

	// Stockの取得
	stock, err := client.Stock.FindMany().With(
		db.Stock.Price.Fetch().With(
			db.Price.Product.Fetch().With(
				db.Product.Seat.Fetch(),
			),
		),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("stock取得エラー : %w", err)
	}

	for _, st := range stock {

		seat := st.RelationsStock.Price.RelationsPrice.Product.RelationsProduct.Seat

		for _, se := range seat {

			// SeatReservationにデータがない場合は生成
			_, err := client.SeatReservation.UpsertOne(
				db.SeatReservation.StockIDSeatID(
					db.SeatReservation.StockID.Equals(st.ID),
					db.SeatReservation.SeatID.Equals(se.ID),
				),
			).Create(
				db.SeatReservation.Stock.Link(
					db.Stock.ID.Equals(st.ID),
				),
				db.SeatReservation.Seat.Link(
					db.Seat.ID.Equals(se.ID),
				),
				db.SeatReservation.IsReserved.Set(false),
			).Update().Exec(ctx)
			if err != nil {
				return fmt.Errorf("seatreservationアップサートエラー : %w", err)
			}
		}
	}

	return nil
}

func init() {
	fmt.Println("[Init start]")

	// オプションの読み込み
	if err := loadOptions(); err != nil {
		panic(fmt.Sprint("オプション読み込みエラー: ", err))
	}

	// テーブルに情報がない場合に自動インサート（テスト用）
	if auto_insert {
		if err := AutoInsert(); err != nil {
			fmt.Println(err)
			fmt.Println("処理を続行します")
		}
	}

	// SeatReservationの作成
	if Options.Seat_enable {
		err := generateSeatReservation()
		if err != nil {
			panic(fmt.Sprint("SeatReservation作成エラー: ", err))
		}
	}

	// 各テーブルのチェック
	// 未実装

}
