package main

import (
	"Go_APIServer/funcs"
	"Go_APIServer/scheduler"
	"fmt"
)

func main() {

	// TimeAPI動作チェック
	fmt.Println("[Test TimeAPI]", funcs.GetTime())

	// スケジューラの起動
	go scheduler.Scheduler()

	// APIServerの起動
	if err := APIServer(); err != nil {
		panic(err)
	}
}
