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
	postImage(api)
}

func post(api *slack.Client) {
	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := api.PostMessage("#b", slack.MsgOptionText("Hello World", true))
	if err != nil {
		panic(err)
	}
}

// チャンネルに参加する
func join(api *slack.Client, channelID string) {
	_, _, _, err := api.JoinConversation(channelID)
	if err != nil {
		panic(err)
	}
}

func postImage(api *slack.Client) {
	file, err := os.Open("yopparai_sakeguse_warui_man.png")
	var param slack.FileUploadParameters
	param.Reader = file
	param.Filename = "upload file name"
	param.Channels = []string{"#a"}
	str, err := api.UploadFile(param)
	fmt.Println(str, err)
}

func group(api *slack.Client) {
	groups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}
}
