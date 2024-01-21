package main

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	fmt.Println(r.ContentLength)
	fmt.Println(r.Form)
	fmt.Println(r.Header)
	fmt.Println(r.Host)
	fmt.Println(r.Method)
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.RequestURI)
	fmt.Println(r.URL)
}
