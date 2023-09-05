package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// チャンネルに所属しているユーザーIDの一覧を取得する
func getUserIDsInChannel(tkn string, channelID string) []string {
	url := "https://slack.com/api/conversations.members?channel=" + channelID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	value := "Bearer " + tkn
	req.Header.Set("Authorization", value)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)

	var usersList UsersList
	json.Unmarshal(byteArray, &usersList)
	if usersList.Ok == true {
		return usersList.Members
	} else {
		return nil
	}
}

type UsersList struct {
	Ok               bool     `json:"ok"`
	Members          []string `json:"members"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

// 入力されたユーザーIDと対応する表示名とアイコンURLとメールアドレスを返す
func getUserInfo(tkn string, userID string) UsersInfo {
	url := "https://slack.com/api/users.info?user=" + userID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	value := "Bearer " + tkn
	req.Header.Set("Authorization", value)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)

	var usersInfo UsersInfo
	json.Unmarshal(byteArray, &usersInfo)
	if usersInfo.Ok == true {
		return usersInfo
	} else {
		var resp UsersInfo
		return resp
	}
}

type UsersInfo struct {
	Ok   bool `json:"ok"`
	User struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		Deleted  bool   `json:"deleted"`
		Color    string `json:"color"`
		RealName string `json:"real_name"`
		Tz       string `json:"tz"`
		TzLabel  string `json:"tz_label"`
		TzOffset int    `json:"tz_offset"`
		Profile  struct {
			Title                  string `json:"title"`
			Phone                  string `json:"phone"`
			Skype                  string `json:"skype"`
			RealName               string `json:"real_name"`
			RealNameNormalized     string `json:"real_name_normalized"`
			DisplayName            string `json:"display_name"`
			DisplayNameNormalized  string `json:"display_name_normalized"`
			Fields                 any    `json:"fields"`
			StatusText             string `json:"status_text"`
			StatusEmoji            string `json:"status_emoji"`
			StatusEmojiDisplayInfo []any  `json:"status_emoji_display_info"`
			StatusExpiration       int    `json:"status_expiration"`
			AvatarHash             string `json:"avatar_hash"`
			ImageOriginal          string `json:"image_original"`
			IsCustomImage          bool   `json:"is_custom_image"`
			Email                  string `json:"email"`
			FirstName              string `json:"first_name"`
			LastName               string `json:"last_name"`
			Image24                string `json:"image_24"`
			Image32                string `json:"image_32"`
			Image48                string `json:"image_48"`
			Image72                string `json:"image_72"`
			Image192               string `json:"image_192"`
			Image512               string `json:"image_512"`
			Image1024              string `json:"image_1024"`
			StatusTextCanonical    string `json:"status_text_canonical"`
			Team                   string `json:"team"`
		} `json:"profile"`
		IsAdmin                bool   `json:"is_admin"`
		IsOwner                bool   `json:"is_owner"`
		IsPrimaryOwner         bool   `json:"is_primary_owner"`
		IsRestricted           bool   `json:"is_restricted"`
		IsUltraRestricted      bool   `json:"is_ultra_restricted"`
		IsBot                  bool   `json:"is_bot"`
		IsAppUser              bool   `json:"is_app_user"`
		Updated                int    `json:"updated"`
		IsEmailConfirmed       bool   `json:"is_email_confirmed"`
		WhoCanShareContactCard string `json:"who_can_share_contact_card"`
	} `json:"user"`
}
