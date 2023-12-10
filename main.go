package main

func main() {
	if err := APIServer(); err != nil {
		panic(err)
	}
}
