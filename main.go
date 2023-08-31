package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadJsonFile()

}

func readEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数の読み込みに失敗しました: %v", err)
	}
	tkn := os.Getenv("APIKEY")
	return tkn
}

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

func postExample() {
	url := "https://slack.com/api/chat.postMessage"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer access-token")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

func loadJsonFile() map[string]any {
	file, err := os.Open("reviewTest.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var rawJson map[string]any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rawJson); err != nil {
		log.Fatal(err)
	}
	fmt.Println(rawJson)
	return rawJson
}
