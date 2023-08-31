package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func useLibrary(tkn string) {
	// アクセストークンを使用してクライアントを生成する
	api := slack.New(tkn)
	post(api)
}

// テキストを投稿する
func post(api *slack.Client) {
	//sendText := "<https://google.com|お店の名前> -ラーメン,中華,定食\n:star::star::star::star:☆\n>本社前の定食屋さんに行ってまいりました！\n>おいしかったです！"
	sendText := ":hello:"

	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := api.PostMessage("#a", slack.MsgOptionText(sendText, true))
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

// 画像を投稿する
func postImage(api *slack.Client) {
	file, err := os.Open("yopparai_sakeguse_warui_man.png")
	var param slack.FileUploadParameters
	param.Reader = file
	param.Filename = "upload file name"
	param.Channels = []string{"#a"}
	str, err := api.UploadFile(param)
	fmt.Println(str, err)
}

// ユーザーIDと名前を取得する
func group(api *slack.Client) {
	groups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(false))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}
}
