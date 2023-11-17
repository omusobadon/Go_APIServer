package main

import "fmt"

type Editer interface {
	Insert()
}

func (s Stock) Insert() {
	fmt.Println("insert")
}

func test() {
	var e Editer

}
