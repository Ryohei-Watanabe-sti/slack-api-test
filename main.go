package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
	tkn := os.Getenv("APIKEY")

	// アクセストークンを使用してクライアントを生成する
	api := slack.New(tkn)
	post(api)
}

func post(api *slack.Client) {
	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := api.PostMessage("#general", slack.MsgOptionText("Hello World", true))
	if err != nil {
		panic(err)
	}
}

func group(api *slack.Client) {
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))
	groups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}
}
