package socket

/* Method(-客戶端): type 必定為偶數 */
const CMD_C_TOKEN_CHANGE string = "2"

const CMD_C_PLAYER_LOGOUT string = "4"

const CMD_C_GET_MEMBER_LIST string = "6"

const CMD_C_PLAYER_EXIT_ROOM string = "8"

const CMD_C_PLAYER_ENTER_ROOM string = "10"

const CMD_C_ROOM_INFO_EDIT string = "16"

const CMD_C_GET_CHAT_HISTORY string = "18"

const CMD_C_FRIEND_INVITE string = "20"

const CMD_C_GET_FRIEND_LIST string = "22"

const CMD_C_CHATBLOCK string = "24"

const CMD_C_PING string = "26"

const CMD_C_PROCLAMATION string = "28"

const CMD_C_FRIEND_DELETE string = "30"

const CMD_C_SIDETEXT_DELETE string = "32"

const CMD_C_GET_SIDETEXT_HISTORY string = "34"

const CMD_C_GET_NEW_SIDETEXT string = "36"

const CMD_C_GET_FUNC_MANAGEMENT string = "38"

const CMD_C_TARGET_ADD_ROOM_BATCH string = "40"

const CMD_C_PLAYER_ENTER_ROOM_BATCH string = "42"

const CMD_C_MESSAGE_SEEN string = "44"

const CMD_C_KICK_ROOM_USER string = "46"

const CMD_C_CREATE_PRIVATE_ROOM string = "48"

const CMD_C_PLAYER_SEND_MSG string = "80"

const CMD_C_PLAYER_SIDETEXT string = "82"

const CMD_C_PLAYER_SEND_SHELL string = "84"

const CMD_C_DIS_MISS_ROOM string = "86"

const CMD_C_ROOM_ADMIN_ADD string = "90"

const CMD_C_ROOM_ADMIN_REMOVE string = "92"

const CMD_C_GET_LANG_LIST string = "96"

/* Message(-伺服器): type 必定為奇數 回覆 */

const CMD_R_TOKEN_CHANGE string = "3"

const CMD_R_PLAYER_LOGOUT string = "5"

const CMD_R_GET_MEMBER_LIST string = "7"

const CMD_R_PLAYER_EXIT_ROOM string = "9"

const CMD_R_PLAYER_ENTER_ROOM string = "11"

const CMD_R_ROOM_INFO_EDIT string = "17"

const CMD_R_GET_CHAT_HISTORY string = "19"

const CMD_R_FRIEND_INVITE string = "21"

const CMD_R_GET_FRIEND_LIST string = "23"

const CMD_R_CHATBLOCK string = "25"

const CMD_R_PING string = "27"

const CMD_R_PROCLAMATION string = "29"

const CMD_R_FRIEND_DELETE string = "31"

const CMD_R_SIDETEXT_DELETE string = "33"

const CMD_R_GET_SIDETEXT_HISTORY string = "35"

const CMD_R_GET_NEW_SIDETEXT string = "37"

const CMD_R_GET_FUNC_MANAGEMENT string = "39"

const CMD_R_TARGET_ADD_ROOM_BATCH string = "41"

const CMD_R_PLAYER_ENTER_ROOM_BATCH string = "43"

const CMD_R_MESSAGE_SEEN string = "45"

const CMD_R_KICK_ROOM_USER string = "47"

const CMD_R_CREATE_PRIVATE_ROOM string = "49"

const CMD_R_PLAYER_SEND_MSG string = "81"

const CMD_R_PLAYER_SIDETEXT string = "83"

const CMD_R_PLAYER_SEND_SHELL string = "85"

const CMD_R_DIS_MISS_ROOM string = "87"

const CMD_R_ROOM_ADMIN_ADD string = "91"

const CMD_R_ROOM_ADMIN_REMOVE string = "93"

const CMD_R_GET_LANG_LIST string = "97"

/* Message(-伺服器): type 必定為奇數 廣播 */

const CMD_B_PLAYER_ENTER_ROOM string = "13"

const CMD_B_PLAYER_EXIT_ROOM string = "15"

const CMD_B_PLAYER_SPEAK string = "51"

const CMD_B_SIDETEXT string = "53"

const CMD_B_SEND_GIFT string = "55"

const CMD_B_SUBSCRIPTION string = "57"

const CMD_B_PROCLAMATION string = "59"

const CMD_B_ADMIN_SHELL string = "61"

const CMD_B_SIDETEXT_DELETE string = "63"

const CMD_B_FUNC_MANAGEMENT string = "65"

const CMD_B_USER_INFO_UPDATE string = "67"

const CMD_B_ROOM_INFO_UPDATE string = "69"

const CMD_B_NOTIFY_ENTER_ROOM string = "71"

const CMD_B_MESSAGE_BE_SEEN string = "73"

const CMD_B_KICK_ROOM_USER string = "75"

const CMD_B_FRIEND string = "77"

const CMD_B_FRIEND_DELETE string = "89"

const CMD_B_GLOBAL_MESSAGE string = "95"

const CMD_B_ROOM_MEMBER_COUNT string = "99"
