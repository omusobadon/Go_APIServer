package ini

import (
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

// オプションの読み込み
func loadOptions() error {

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
