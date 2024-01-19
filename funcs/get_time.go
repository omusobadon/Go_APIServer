// 時刻取得処理
package funcs

import (
	"Go_APIServer/ini"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// World Time APIからのレスポンスを変換するための構造体
type TimeResponse struct {
	Datetime string `json:"datetime"`
}

func GetTime() time.Time {

	got_time, err := GetTimeFromAPI()
	if err != nil {
		fmt.Println("APIからの時刻取得に失敗:", err)
		fmt.Println("システム時刻を使用")

		return time.Now().In(ini.Timezone)
	}

	fixed_time := got_time.In(ini.Timezone)

	return fixed_time
}

// World Time APIから時刻を取得し返却
func GetTimeFromAPI() (*time.Time, error) {

	res, err := http.Get("http://worldtimeapi.org/api/timezone/UTC")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(io.Reader(res.Body))
	if err != nil {
		return nil, err
	}

	var time_res TimeResponse
	if err := json.Unmarshal(body, &time_res); err != nil {
		return nil, err
	}

	time_parsed, err := time.Parse(time.RFC3339, time_res.Datetime)
	if err != nil {
		return nil, err
	}

	return &time_parsed, nil
}
