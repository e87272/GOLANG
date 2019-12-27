package socket

/* Method(-客戶端): type 必定為偶數 */
type Cmd_c_token_change struct {
	Base_C
	Payload Logininfo `json:"payload"`
}

type Cmd_c_get_member_list struct {
	Base_C
	Payload Roomcore `json:"payload"`
}

type Cmd_c_player_logout struct {
	Base_C
}

type Cmd_c_player_send_msg struct {
	Base_C
	Payload Sendmessage `json:"payload"`
}

type Cmd_c_player_exit_room struct {
	Base_C
	Payload Roomcore `json:"payload"`
}

type Cmd_c_player_enter_room struct {
	Base_C
	Payload Roomstation `json:"payload"`
}

type Cmd_c_room_info_edit struct {
	Base_C
	Payload struct {
		Roomcore   Roomcore `json:"roomCore"`
		Roomname   string   `json:"roomName"`
		Roomicon   string   `json:"roomIcon"`
		Targetuuid string   `json:"targetUuid"`
	} `json:"payload"`
}

type Cmd_c_get_room_chat_history struct {
	Base_C
	Payload struct {
		Roomcore    Roomcore `json:"roomCore"`
		Historyuuid string   `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_c_friend_invite struct {
	Base_C
	Payload string `json:"payload"`
}

type Cmd_c_get_friend_list struct {
	Base_C
}

type Cmd_c_player_side_text struct {
	Base_C
	Payload Sendmessage `json:"payload"`
}

type Cmd_c_get_side_text_history struct {
	Base_C
	Payload struct {
		Chattarget  string `json:"chatTarget"`
		Historyuuid string `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_c_chatblock struct {
	Base_C
	Payload struct {
		Useruuid  string `json:"userUuid"`
		Roomuuid  string `json:"roomUuid"`
		Blocktime string `json:"blockTime"`
		Blockip   string `json:"blockIp"`
	} `json:"payload"`
}

type Cmd_c_healthcheck struct {
	Base_C
}

type Cmd_c_proclamation struct {
	Base_C
	Payload Roomcore `json:"payload"`
}

type Cmd_c_friend_delete struct {
	Base_C
	Payload string `json:"payload"`
}
type Cmd_c_side_text_delete struct {
	Base_C
	Payload string `json:"payload"`
}

type Cmd_c_get_chat_history struct {
	Base_C
	Payload struct {
		Chattarget  string `json:"chatTarget"`
		Historyuuid string `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_c_player_send_shell struct {
	Base_C
	Payload Sendmessage `json:"payload"`
}

type Cmd_c_get_new_side_text struct {
	Base_C
}

type Cmd_c_get_func_management struct {
	Base_C
}

type Cmd_c_target_add_room_batch struct {
	Base_C
	Payload struct {
		Room       []Roomcore `json:"room"`
		Targetuuid string     `json:"targetUuid"`
	} `json:"payload"`
}

type Cmd_c_player_enter_room_batch struct {
	Base_C
	Payload []Roomcore `json:"payload"`
}

type Cmd_c_message_seen struct {
	Base_C
	Payload struct {
		Chattype       string `json:"chatType"`
		Chattargetuuid string `json:"chatTargetUuid"`
	} `json:"payload"`
}

type Cmd_c_kick_room_user struct {
	Base_C
	Payload struct {
		Roomcore   Roomcore `json:"roomCore"`
		Targetuuid string   `json:"targetUuid"`
	} `json:"payload"`
}

type Cmd_c_create_private_room struct {
	Base_C
	Payload struct {
		Roomname   string `json:"roomName"`
		Roomicon   string `json:"roomIcon"`
		Targetuuid string `json:"targetUuid"`
	} `json:"payload"`
}

type Cmd_c_dis_miss_room struct {
	Base_C
	Payload string `json:"payload"`
}

type Cmd_c_room_admin_add struct {
	Base_C
	Payload struct {
		Targetuuid string `json:"targetUuid"`
		Roomuuid   string `json:"roomUuid"`
		Role       string `json:"role"`
	} `json:"payload"`
}

type Cmd_c_get_lang_list struct {
	Base_C
	Payload string `json:"payload"`
}

type Cmd_c_clear_user_msg struct {
	Base_C
	Payload Clearusermsg `json:"payload"`
}

type Cmd_c_forward_msg struct {
	Base_C
	Payload Sendmessage `json:"payload"`
}

/* Message(-伺服器): type 必定為奇數 回覆 */

type Cmd_r_token_change struct {
	Base_R
	Payload User `json:"payload"`
}
type Cmd_r_get_member_list struct {
	Base_R
	Payload []Userplatform `json:"payload"`
}

type Cmd_r_player_logout struct {
	Base_R
}

type Cmd_r_player_send_msg struct {
	Base_R
}

type Cmd_r_player_exit_room struct {
	Base_R
}

type Cmd_r_player_enter_room struct {
	Base_R
	Payload Roominfo `json:"payload"`
}

type Cmd_r_room_info_edit struct {
	Base_R
}

type Cmd_r_get_chat_history struct {
	Base_R
	Payload struct {
		Chattarget  string        `json:"chatTarget"`
		Message     []Chatmessage `json:"message"`
		Historyuuid string        `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_r_friend_invite struct {
	Base_R
}

type Cmd_r_get_friend_list struct {
	Base_R
	Payload struct {
		Friendlist     []Userplatform `json:"friendList"`
		Invitefromlist []Userplatform `json:"inviteFromList"`
		Invitetolist   []Userplatform `json:"inviteToList"`
	} `json:"payload"`
}

type Cmd_r_player_side_text struct {
	Base_R
}

type Cmd_r_chatblock struct {
	Base_R
	Payload struct {
		Useruuid string `json:"userUuid"`
		Roomuuid string `json:"roomUuid"`
	} `json:"payload"`
}

type Cmd_r_healthcheck struct {
	Base_R
	Payload string `json:"payload"`
}

type Cmd_r_proclamation struct {
	Base_R
	Payload map[string]Proclamation `json:"payload"`
}

type Cmd_r_player_send_shell struct {
	Base_R
	Payload string `json:"payload"`
}

type Cmd_r_friend_delete struct {
	Base_R
}

type Cmd_r_player_side_text_delete struct {
	Base_R
}

type Cmd_r_get_side_text_history struct {
	Base_R
	Payload struct {
		Chattarget  string        `json:"chatTarget"`
		Message     []Chatmessage `json:"message"`
		Historyuuid string        `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_r_get_new_side_text struct {
	Base_R
	Payload struct {
		Lastmessageuuid string        `json:"lastMessageUuid"`
		Newsidetextlist []Newsidetext `json:"newSidetextList"`
	} `json:"payload"`
}

type Cmd_r_get_func_management struct {
	Base_R
	Payload map[string]string `json:"payload"`
}

type Cmd_r_target_add_room struct {
	Base_R
}

type Cmd_r_target_add_room_batch struct {
	Base_R
	Payload []struct {
		Result   string   `json:"result"`
		Roomcore Roomcore `json:"roomCore"`
	} `json:"payload"`
}

type Cmd_r_player_enter_room_batch struct {
	Base_R
	Payload []struct {
		Result      string      `json:"result"`
		Roominfo    Roominfo    `json:"roomInfo"`
		Lastmessage Chatmessage `json:"lastMessage"`
		Membercount int         `json:"memberCount"`
	} `json:"payload"`
}

type Cmd_r_message_seen struct {
	Base_R
}

type Cmd_r_kick_room_user struct {
	Base_R
}

type Cmd_r_create_private_room struct {
	Base_R
	Payload Roomcore `json:"payload"`
}

type Cmd_r_dis_miss_room struct {
	Base_R
}

type Cmd_r_room_admin_add struct {
	Base_R
}

type Cmd_r_room_admin_remove struct {
	Base_R
}

type Cmd_r_get_lang_list struct {
	Base_R
	Payload map[string]string `json:"payload"`
}

type Cmd_r_clear_user_msg struct {
	Base_R
}

type Cmd_r_forward_msg struct {
	Base_R
}

/* Message(-伺服器): type 必定為奇數 廣播 */

type Cmd_b_player_room struct {
	Base_B
	Payload struct {
		Chatmessage Chatmessage `json:"chatMessage"`
		Chattarget  string      `json:"chatTarget"`
	} `json:"payload"`
}

type Cmd_b_player_speak struct {
	Base_B
	Payload struct {
		Chatmessage Chatmessage `json:"chatMessage"`
		Chattarget  string      `json:"chatTarget"`
	} `json:"payload"`
}

type Cmd_b_side_text struct {
	Base_B
	Payload struct {
		Chatmessage Chatmessage `json:"chatMessage"`
		Chattarget  string      `json:"chatTarget"`
	} `json:"payload"`
}

type Cmd_b_send_gift struct {
	Base_B
	Payload struct {
		Chatmessage Chatmessage `json:"chatMessage"`
		Chattarget  string      `json:"chatTarget"`
	} `json:"payload"`
}

type Cmd_b_subscription struct {
	Base_B
	Payload struct {
		Chatmessage Chatmessage `json:"chatMessage"`
		Chattarget  string      `json:"chatTarget"`
	} `json:"payload"`
}

type Cmd_b_proclamation struct {
	Base_B
	Payload map[string]Proclamation `json:"payload"`
}

type Cmd_b_admin_shell struct {
	Base_B
	Payload Sudoresult `json:"payload"`
}

type Cmd_b_sidetext_delete struct {
	Base_B
	Payload string `json:"payload"`
}

type Cmd_b_func_management struct {
	Base_B
	Payload map[string]string `json:"payload"`
}

type Cmd_b_user_info_update struct {
	Base_B
	Payload User `json:"payload"`
}

type Cmd_b_room_info_update struct {
	Base_B
	Payload Roominfo `json:"payload"`
}

type Cmd_b_target_add_room struct {
	Base_B
	Payload []Roomcore `json:"payload"`
}

type Cmd_b_message_be_seen struct {
	Base_B
	Payload struct {
		Chattarget  string `json:"chatTarget"`
		Historyuuid string `json:"historyUuid"`
	} `json:"payload"`
}

type Cmd_b_kick_room_user struct {
	Base_B
	Payload Roomcore `json:"payload"`
}

type Cmd_b_friend struct {
	Base_B
	Payload Friendplatform `json:"payload"`
}

type Cmd_b_friend_delete struct {
	Base_B
	Payload string `json:"payload"`
}

type Cmd_b_global_message struct {
	Base_B
	Payload Globalmessage `json:"payload"`
}

type Cmd_b_room_member_count struct {
	Base_B
	Payload struct {
		Roomuuid string `json:"roomUuid"`
		Count    int    `json:"count"`
	} `json:"payload"`
}

type Cmd_b_clear_user_msg struct {
	Base_B
	Payload Clearusermsg `json:"payload"`
}
