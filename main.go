package main

import (
	"fmt"
)

func main() {
	fmt.Println("[ Mode Select ]")

	var i int
	for {
		//1:Postプログラムの実行、9:テストプログラムの実行、0:終了
		fmt.Print("[1]:APIServer / [2]:NTPServer Test / [9]:Test / [0]:Exit -> ")
		fmt.Scan(&i)

		switch i {
		case 0:
		case 1:
			if err := APIServer(); err != nil {
				panic(err)
			}
		case 2:
			fmt.Println("*** NTP TEST ***")
			fmt.Println(GetTime())
			fmt.Println("*** NTP TEST End ***")
		case 9:
			test()
		default:
			fmt.Println("指定された数値を入力")
		}

		if i == 0 {
			break
		}
	}
}
