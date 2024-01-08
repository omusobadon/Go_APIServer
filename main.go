package main

func main() {
	go scheduler()

	if err := APIServer(); err != nil {
		panic(err)
	}
}
