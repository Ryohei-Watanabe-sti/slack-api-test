package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	// アクセストークンを使用してクライアントを生成する
	tkn := "xoxb-5807191892403-5807304960274-aKfXEBIbtGDUX45j7adp4QSd"
	api := slack.New(tkn)
	group(api)
}

func post(api *slack.Client) {
	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := api.PostMessage("#public", slack.MsgOptionText("Hello World", true))
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
