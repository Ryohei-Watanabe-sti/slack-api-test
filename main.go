package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	url := "http://example.com"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

func readEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	tkn := os.Getenv("APIKEY")
	return tkn
}
