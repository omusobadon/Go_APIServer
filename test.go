package main

import (
	"errors"
	"fmt"
)

func test() {
	e := fmt.Errorf("NotFound")
	err := fmt.Errorf("Order : %w", e)

	fmt.Println(err)

	fmt.Println(errors.Is(err, err))

}
