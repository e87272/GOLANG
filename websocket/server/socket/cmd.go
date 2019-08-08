package socket
/* Method(-客戶端): type 必定為偶數 */
const CMD_C_TOKEN_CHANGE	   string = "2"

const CMD_C_PLAYER_LOGOUT	   string = "4"

const CMD_C_GET_MEMBER_LIST	   string = "6"

const CMD_C_PLAYER_EXIT_ROOM   string = "8"

const CMD_C_PLAYER_ENTER_ROOM   string = "10"

const CMD_C_GET_ROOM_LIST	   string = "16"

const CMD_C_PLAYER_SEND_MSG	   string = "80"




/* Message(-伺服器): type 必定為奇數 回覆 */

const CMD_R_TOKEN_CHANGE	   string = "3"

const CMD_R_PLAYER_LOGOUT	   string = "5"

const CMD_R_GET_MEMBER_LIST	   string = "7"

const CMD_R_PLAYER_EXIT_ROOM   string = "9"

const CMD_R_PLAYER_ENTER_ROOM  string = "11"

const CMD_R_GET_ROOM_LIST	   string = "17"

const CMD_R_PLAYER_SEND_MSG	   string = "81"



/* Message(-伺服器): type 必定為奇數 廣播 */

const CMD_B_PLAYER_ENTER_ROOM  string = "13"

const CMD_B_PLAYER_EXIT_ROOM   string = "15"

const CMD_B_PLAYER_SPEAK	   string = "51"

