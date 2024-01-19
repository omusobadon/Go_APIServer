package ini

import (
	"Go_APIServer/db"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var Options LoadedOptions
var Timezone *time.Location

// config.ymlデコード用構造体
type LoadedOptions struct {
	Time_free_enable  bool   `yaml:"time_free_enable"`
	Seat_enable       bool   `yaml:"seat_enable"`
	Hold_enable       bool   `yaml:"hold_enable"`
	Payment_enable    bool   `yaml:"payment_enable"`
	User_end_enable   bool   `yaml:"user_end_enable"`
	User_notification bool   `yaml:"user_notification"`
	Delay             int    `yaml:"delay"`
	Margin            int    `yaml:"margin"`
	Local_time_enable bool   `yaml:"local_time_enable"`
	Timezone          string `yaml:"timezone"`
	Time_difference   int    `yaml:"time_difference"`
}

// 商品・在庫テーブルが空の場合、自動生成するかどうか（テスト用）
const auto_insert bool = true

// オプションの読み込み
func LoadOptions() error {

	// config.ymlの読み込み
	content, err := os.ReadFile("config/config.yml")
	if err != nil {
		return err
	}

	// ymlのデコード
	err = yaml.Unmarshal(content, &Options)
	if err != nil {
		return err
	}

	// ローカルタイムの設定
	if Options.Local_time_enable {
		Timezone = time.Local

	} else {
		Timezone, err = time.LoadLocation(Options.Timezone)
		if err != nil {
			Timezone = time.FixedZone(Options.Timezone, Options.Time_difference*60*60)
		}
	}

	// optionの確認
	fmt.Println("[Option]")
	// fmt.Printf("Options: %+v\n", Options)

	fmt.Print("ユーザ時刻指定: ")
	if Options.Time_free_enable {
		fmt.Println("有効")
	} else {
		fmt.Println("無効")
	}

	fmt.Print("座席指定: ")
	if Options.Seat_enable {
		fmt.Println("有効")
	} else {
		fmt.Println("無効")
	}

	fmt.Printf("スケジューラチェック間隔: %ds, マージン: %ds\n", Options.Delay, Options.Margin)
	fmt.Println("タイムゾーン:", Timezone)

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
	if err := LoadOptions(); err != nil {
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
