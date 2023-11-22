package main

import (
	"fmt"
	"time"
)

func test() {
	// func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
	loc, _ := time.LoadLocation("Asia/Tokyo")
	time := time.Date(2014, 12, 31, 8, 4, 18, 0, loc)

	fmt.Println(time)

	interval := time.Date(2020, 1, 2, 3, 4, 5, 123456789, loc)

	fmt.Println(time.Add(interval))
}
