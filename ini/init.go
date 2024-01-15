package ini

import (
	"Go_APIServer/db"
	"context"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var OPTIONS Options

// config.ymlデコード用構造体
type Options struct {
	Time_free_enable  bool
	Seat_enable       bool
	Hold_enable       bool
	Payment_enable    bool
	User_end_enable   bool
	User_notification bool
	Timezone          string
	Delay             int
}

// 商品・在庫テーブルが空の場合、自動生成するかどうか（テスト用）
const auto_insert bool = true

// オプションの読み込み
func ReadOptions() error {

	// config.ymlの読み込み
	content, err := os.ReadFile("config/config.yml")
	if err != nil {
		return err
	}

	// ymlのデコード
	err = yaml.Unmarshal(content, &OPTIONS)
	if err != nil {
		return err
	}

	// optionの確認
	fmt.Printf("OPTIONS: %+v\n", OPTIONS)

	return nil
}

// SeatReservationの自動生成
func generateSeatReservation() error {

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
	stock, err := client.Stock.FindMany().Exec(ctx)
	if err != nil {
		return fmt.Errorf("stock取得エラー : %w", err)
	}

	for _, st := range stock {

		// Seatの取得
		seat, err := client.Seat.FindMany().Exec(ctx)
		if err != nil {
			return fmt.Errorf("seat取得エラー : %w", err)
		}

		for _, se := range seat {
			// SeatReservationにデータがある場合はスキップ
			_, err := client.SeatReservation.FindFirst(
				db.SeatReservation.StockID.Equals(st.ID),
				db.SeatReservation.SeatID.Equals(se.ID),
			).Exec(ctx)

			if errors.Is(err, db.ErrNotFound) {

			} else if err != nil {
				return fmt.Errorf("seat取得エラー : %w", err)

			} else {
				continue
			}

			// SeatReservationにインサート
			_, err = client.SeatReservation.CreateOne(
				db.SeatReservation.Stock.Link(
					db.Stock.ID.Equals(st.ID),
				),
				db.SeatReservation.Seat.Link(
					db.Seat.ID.Equals(se.ID),
				),
				db.SeatReservation.IsReserved.Set(false),
			).Exec(ctx)
			if err != nil {
				return fmt.Errorf("SeatReservationインサートエラー : %w", err)
			}
		}
	}

	return nil
}

func init() {
	fmt.Println("[Init start]")

	// オプションの読み込み
	if err := ReadOptions(); err != nil {
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
	if OPTIONS.Seat_enable {
		err := generateSeatReservation()
		if err != nil {
			panic(fmt.Sprint("SeatReservation作成エラー: ", err))
		}
	}

	// 各テーブルのチェック
	// 未実装

}
