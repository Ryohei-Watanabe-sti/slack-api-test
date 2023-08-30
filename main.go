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
	examplePost(api)
}

// テキストを投稿する
func post(api *slack.Client) {
	sendText := "<https://google.com|お店の名前> -ラーメン,中華,定食\n:star::star::star::star:☆\n>本社前の定食屋さんに行ってまいりました！\n>おいしかったです！"
	// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
	_, _, err := api.PostMessage("#a", slack.MsgOptionText(sendText, true))
	if err != nil {
		panic(err)
	}
}

// 複雑なテキストを投稿する例
func examplePost(api *slack.Client) {
	attachment := slack.Attachment{
		Pretext:    "This is slack post test by Go",
		Title:      "title",
		Color:      "#36a64f",
		AuthorName: "author_name",
		AuthorIcon: "https://placeimg.com/16/16/people",
		MarkdownIn: []string{"`textTomarkdown`"},
		Text:       "hello world `textTomarkdown`",
		ThumbURL:   "http://placekitten.com/g/200/200",
		FooterIcon: "https://platform.slack-edge.com/img/default_application_icon.png",

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Item1",
				Value: "this is value of item1",
				Short: false,
			}, slack.AttachmentField{
				Title: "Item2",
				Value: "this is value of item2",
				Short: true,
			}, slack.AttachmentField{
				Title: "Item3",
				Value: "```" + "this is value of item3" + "```",
				Short: false,
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(
		"#b",
		slack.MsgOptionText("This is Title", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
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
