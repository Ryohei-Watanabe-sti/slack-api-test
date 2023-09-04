package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	tkn := readEnv()
	file, err := os.Open("requestBody.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := io.Reader(file)
	postMessage(tkn, body)
}

func readEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数の読み込みに失敗しました: %v", err)
	}
	tkn := os.Getenv("APIKEY")
	return tkn
}

/*
func getExample() {
	url := "http://example.com"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}
*/

func postMessage(tkn string, body io.Reader) {
	url := "https://slack.com/api/chat.postMessage"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
	}
	value := "Bearer " + tkn
	req.Header.Set("Authorization", value)
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}
