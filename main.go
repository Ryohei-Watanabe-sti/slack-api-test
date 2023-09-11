package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	tkn, socketKey := readEnv()
	//requestBody.jsonのメッセージを投稿
	// file, err := os.Open("requestBody.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	// body := io.Reader(file)
	// postMessage(tkn, body)
	// uploadFile(tkn)

	//チャンネルのユーザーのid,name,email,iconをprint表示
	// usersList := getUserIDsInChannel(tkn, "C05QEV93KTJ")
	// for _, v := range usersList {
	// 	var resp UsersInfo
	// 	resp = getUserInfo(tkn, v)
	// 	fmt.Println("id:", v)
	// 	fmt.Println("name:", resp.User.Name)
	// 	fmt.Println("email:", resp.User.Profile.Email)
	// 	fmt.Println("iconImage:", resp.User.Profile.Image1024)
	// }

	//socketモード
	socket(tkn, socketKey)
}

// APIキーの読み込み
func readEnv() (string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数の読み込みに失敗しました: %v", err)
	}
	tkn := os.Getenv("APIKEY")
	socketKey := os.Getenv("SOCKEY")
	return tkn, socketKey
}

// bodyに渡された投稿先のchannelとメッセージの内容を投稿する
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

// 画像ファイルをchannel:#aにアップロードする
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

// ソケットモード
func socket(botToken string, appToken string) {

	//tokenのチェック
	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	//clientの作成
	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	//SocketModeの起動、イベント検出時に処理を分岐
	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					//メンションされたときはここが実行
					case *slackevents.AppMentionEvent:
						//イベントを受信したチャンネルにイベント内容と同じテキストを投稿
						_, _, err := client.PostMessage(ev.Channel, slack.MsgOptionText(ev.Text, false))
						if err != nil {
							fmt.Printf("failed posting message: %v", err)
						}

					case *slackevents.MemberJoinedChannelEvent:
						//ユーザーがチャンネルに参加すると、ようこそとメッセージが送られる
						welcome(botToken, ev.User, ev.Channel)
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					fmt.Println("ブロックアクション！")

					// ブロックアクションが発生した場合の処理
					for _, action := range callback.ActionCallback.BlockActions {
						fmt.Println("ブロックアクション処理開始！")

						if action.ActionID == "static_select_arrive" {
							// static_select のアクションIDでフィルタリング
							fmt.Println("プルダウンメニューだあああ！")
							// 選択されたオプションを取得
							selectedOption := action.SelectedOption.Value

							var sendingText string
							// 選択されたオプションに対する処理を実行
							switch selectedOption {
							case "value-0":
								sendingText = "111"
							case "value-1":
								sendingText = "222"
							case "value-2":
								sendingText = "333"
							default:
								// 未知のオプションが選択された場合の処理を追加
							}
							_, _, err := api.PostMessage(callback.Channel.ID, slack.MsgOptionText(sendingText, false))
							if err != nil {
								fmt.Println("PostMessageError:", err)
							}
						}
					}
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:
					fmt.Println("未知のコールバックを受信")

				}

				client.Ack(*evt.Request, payload)
			//スラッシュコマンドを受け取ったとき
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				client.Debugf("Slash command received: %+v", cmd)

				switch cmd.Command {
				case "/入荷":
					payload := map[string]interface{}{
						"blocks": []slack.Block{
							slack.NewSectionBlock(
								&slack.TextBlockObject{
									Type: slack.MarkdownType,
									Text: "入荷記録",
								},
								nil,
								slack.NewAccessory(
									&slack.SelectBlockElement{
										Type: slack.OptTypeStatic,
										Placeholder: &slack.TextBlockObject{
											Type:  slack.PlainTextType,
											Text:  "商品を選択",
											Emoji: true,
										},
										ActionID: "static_select_arrive",
										Options: []*slack.OptionBlockObject{
											{
												Text: &slack.TextBlockObject{
													Type:  slack.PlainTextType,
													Text:  "111",
													Emoji: true,
												},
												Value: "value-0",
											},
											{
												Text: &slack.TextBlockObject{
													Type:  slack.PlainTextType,
													Text:  "222",
													Emoji: true,
												},
												Value: "value-1",
											},
											{
												Text: &slack.TextBlockObject{
													Type:  slack.PlainTextType,
													Text:  "333",
													Emoji: true,
												},
												Value: "value-2",
											},
										},
									},
								),
							),
						},
					}
					client.Ack(*evt.Request, payload)

				case "/hello":
					text := cmd.UserName + "さん、おはこんばんにちは！"
					payload := map[string]interface{}{
						"blocks": []slack.Block{
							slack.NewSectionBlock(
								&slack.TextBlockObject{
									Type: slack.MarkdownType,
									Text: text,
								},
								nil,
								nil,
							),
						},
					}
					client.Ack(*evt.Request, payload)
				}

			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()

}
