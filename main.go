package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello")
	resp, err := http.Get("http://api.github.com/repos/fusor/fusor/pulls")
	if err != nil {
		fmt.Println("Error getting", err)
	}
	defer resp.Body.Close()
	var v map[string]interface{}
	prs := json.NewDecoder(resp.Body).Decode(&v)
	fmt.Println(prs)
}
