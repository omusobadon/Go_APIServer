package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// World Time APIからのレスポンスを変換するための構造体
type TimeResponse struct {
	Dateline string `json:"datetime"`
}

// World Time APIから時刻を取得し返却
func GetTime() time.Time {
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/Asia/Tokyo")
	if err != nil {
		fmt.Println("時刻取得エラー : ", err)

		// エラーが出た場合はシステム時刻を返す
		return time.Now()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		fmt.Println("レスポンス読み込みエラー : ", err)

		return time.Now()
	}

	var time_res TimeResponse
	if err := json.Unmarshal(body, &time_res); err != nil {
		fmt.Println("レスポンスの構造体変換エラー : ", err)

		return time.Now()
	}

	time_parsed, err := time.Parse(time.RFC3339, time_res.Dateline)
	if err != nil {
		fmt.Println("時刻の解析エラー : ", err)

		return time.Now()
	}

	return time_parsed
}
