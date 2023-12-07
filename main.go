package main

import (
	"fmt"
)

func main() {
	if err := APIServer(); err != nil {
		fmt.Println(err)
	}
}
