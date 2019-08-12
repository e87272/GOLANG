package socket

/* Method(-客戶端): type 必定為偶數 */
type Cmd_c_token_change_struct struct {
	Base_C
	Payload	 LoginInfo `json:"payload"`
}

type Cmd_c_get_member_list struct {
	Base_C
	Payload	 RoomInfo   `json:"payload"`
}

type Cmd_c_player_logout_struct struct {
	Base_C
	Payload	 RoomInfo `json:"payload"`
}

type Cmd_c_player_send_msg_struct struct {
	Base_C
	Payload	Chat     `json:"payload"`
}

type Cmd_c_player_exit_room_struct struct {
	Base_C
	Payload	RoomInfo     `json:"payload"`
}

type Cmd_c_player_enter_room_struct struct {
	Base_C
	Payload	RoomInfo     `json:"payload"`
}

type Cmd_c_get_room_list struct {
	Base_C
}

/* Message(-伺服器): type 必定為奇數 回覆 */

type Cmd_r_token_change_struct struct {
	Base_R
}
type Cmd_r_get_member_list_struct struct {
	Base_R
	Payload	[] User        `json:"payload"`
}

type Cmd_r_player_logout_struct struct {
	Base_R
}

type Cmd_r_player_send_msg_struct struct {
	Base_R
}

type Cmd_r_player_exit_room_struct struct {
	Base_R
}

type Cmd_r_player_enter_room_struct struct {
	Base_R
}

type Cmd_r_get_room_list_struct struct {
	Base_R
	Payload	[] RoomInfo        `json:"payload"`
}

/* Message(-伺服器): type 必定為奇數 廣播 */

type Cmd_b_player_enter_room_struct struct {
	Base_B
	Payload	ChatMessage     `json:"payload"`
}

type Cmd_b_player_exit_room_struct struct {
	Base_B
	Payload	ChatMessage     `json:"payload"`
}

type Cmd_b_player_speak_struct struct {
	Base_B
	Payload	ChatMessage `json:"payload"`
}
