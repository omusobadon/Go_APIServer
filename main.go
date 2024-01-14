package main

import "Go_APIServer/scheduler"

func main() {

	// スケジューラの起動
	go scheduler.Scheduler()

	// APIServerの起動
	if err := APIServer(); err != nil {
		panic(err)
	}
}
