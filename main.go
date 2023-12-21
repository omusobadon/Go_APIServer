package main

import "fmt"

func main() {
	fmt.Println(auto_insert)
	if err := APIServer(); err != nil {
		panic(err)
	}
}
