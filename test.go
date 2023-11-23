package main

import (
	"fmt"
	"time"
)

func test() {
	// Tokyo タイムゾーンを設定
	loc, _ := time.LoadLocation("Asia/Tokyo")

	// time.Time オブジェクトの生成
	t := time.Date(2014, 12, 31, 8, 4, 18, 0, loc)
	fmt.Println(t)

	// 時間間隔（Duration）の生成
	// ここでは、2020年1月2日3時4分5秒と123456789ナノ秒の間隔を設定します
	// ただし、このような間隔の設定は一般的ではなく、通常は時間や分などで間隔を指定します
	interval := time.Duration(5*time.Second + 4*time.Minute + 3*time.Hour + 2*time.Hour*24)

	// 時間の追加
	newTime := t.Add(interval)
	fmt.Println(newTime)
}
