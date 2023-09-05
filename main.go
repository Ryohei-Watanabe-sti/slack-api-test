package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

	// body := io.Reader(file)
	// postMessage(tkn, body)
	// uploadFile(tkn)
	usersList := getUserIDsInChannel(tkn, "C05QEV93KTJ")
	for _, v := range usersList {
		var resp UsersInfo
		resp = getUserInfo(tkn, v)
		fmt.Println("id:", v)
		fmt.Println("name:", resp.User.Name)
		fmt.Println("email:", resp.User.Profile.Email)
		fmt.Println("iconImage:", resp.User.Profile.Image1024)
	}
}

func readEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数の読み込みに失敗しました: %v", err)
	}
	tkn := os.Getenv("APIKEY")
	return tkn
}

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
	fmt.Println("POST MESSAGE：")
	fmt.Println(string(byteArray))
}

func uploadFile(tkn string) {
	// リクエストボディを作成
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// ファイルを添付
	file, err := os.Open("yopparai_sakeguse_warui_man.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "yopparai_sakeguse_warui_man.png")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	writer.WriteField("channels", "#a")
	writer.WriteField("username", "My bot")

	writer.Close()

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", &requestBody)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	value := "Bearer " + tkn
	req.Header.Set("Authorization", value)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	fmt.Println("ファイルアップロード：")
	fmt.Println(string(byteArray))

}
