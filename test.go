package main

import (
	"fmt"
	"time"
)

func test() {
	s := "99h99m"
	time, _ := time.ParseDuration(s)

	fmt.Printf("型 : %T, 値 : %v\n", time, time)

	// test

}
