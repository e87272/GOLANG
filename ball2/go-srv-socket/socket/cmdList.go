package socket

/* Method(-客戶端): type 必定為偶數 */

const CMD_C_TOKEN_CHANGE string = "2"

const CMD_C_PLAYER_ENTER_ROOM string = "4"

const CMD_C_PLAYER_SEND_MSG string = "6"

const CMD_C_LIVE_GAME_INFO string = "8"

/* Message(-伺服器): type 必定為奇數 回覆 */

const CMD_R_TOKEN_CHANGE string = "3"

const CMD_R_PLAYER_ENTER_ROOM string = "5"

const CMD_R_PLAYER_SEND_MSG string = "7"

const CMD_R_LIVE_GAME_INFO string = "9"

/* Message(-伺服器): type 必定為奇數 廣播 */

const CMD_B_PLAYER_ROOM_MSG string = "1001"
