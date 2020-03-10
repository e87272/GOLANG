package socket

import "../commonData"

type base_C struct {
	Cmd  string `json:"cmd"`
	Idem string `json:"idem"`
}

type base_R struct {
	Cmd    string               `json:"cmd"`
	Idem   string               `json:"idem"`
	Stamp  string               `json:"stamp"`
	Result string               `json:"result"`
	Exp    commonData.Exception `json:"exp"`
}

type base_B struct {
	Cmd   string `json:"cmd"`
	Stamp string `json:"stamp"`
}

/* Method(-客戶端): type 必定為偶數 */

type cmd_c_token_change struct {
	base_C
	Payload string `json:"payload"`
}

type cmd_c_player_enter_room struct {
	base_C
	Payload string `json:"payload"`
}

type cmd_c_player_send_msg struct {
	base_C
	Payload commonData.SendMessage `json:"payload"`
}

type cmd_c_live_game_info struct {
	base_C
	Payload string `json:"payload"`
}

/* Message(-伺服器): type 必定為奇數 回覆 */

type cmd_r_token_change struct {
	base_R
	Payload commonData.UserInfo `json:"payload"`
}

type cmd_r_player_enter_room struct {
	base_R
}

type cmd_r_player_send_msg struct {
	base_R
}

type cmd_r_live_game_info struct {
	base_R
	Payload []commonData.ChatMessage `json:"payload"`
}

/* Message(-伺服器): type 必定為奇數 廣播 */

type cmd_b_player_room_msg struct {
	base_B
	Payload struct {
		ChatTarget  string                 `json:"chatTarget"`
		ChatMessage commonData.ChatMessage `json:"chatMessage"`
	} `json:"payload"`
}
