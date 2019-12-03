package common

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"../socket"
)

var mutexClientMsgList sync.Mutex

var clientMsgList = map[string]string{
	"API_BLOCKROOMCHAT_DB_DELETE_ERROR":                           "001df3d4a8622000", //系统发生错误
	"API_BLOCKROOMCHAT_DB_INSERT_ERROR":                           "001df3d4a8622000", //系统发生错误
	"API_BLOCKROOMCHAT_DB_NO_DATA":                                "001df3d4a8622000", //系统发生错误
	"API_BLOCKROOMCHAT_TIME_ERROR":                                "001df3d4a8622000", //系统发生错误
	"API_BLOCKSEARCHROOM_DB_SELECT_ERROR":                         "001df3d4a8622000", //系统发生错误
	"API_BLOCKSEARCHROOM_JSON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"API_BLOCKSEARCHUSER_JSON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"API_CREATEAPICLIENT_DB_INSERT_ERROR":                         "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_INSERT_ROOM_ERROR":                            "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_INSERT_ROOM_USER_ERROR":                       "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_ROOM_ICON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_ROOM_ICON_TOO_LARGE":                          "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_ROOM_TYPE_ERROR_LARGE":                        "001df3d4a8622000", //系统发生错误
	"API_CREATEROOM_SELECT_USER_ERROR":                            "001df3d4a8622000", //系统发生错误
	"API_ROOMINFOEDIT_ROOM_ICON_ERROR":                            "001df3d4a8622000", //系统发生错误
	"API_ROOMINFOEDIT_ROOM_ICON_TOO_LARGE":                        "001df3d4a8622000", //系统发生错误
	"API_ROOMINFOEDIT_ROOM_TYPE_ERROR_LARGE":                      "001df3d4a8622000", //系统发生错误
	"API_ROOMINFOEDIT_UPDATE_ROOM_ERROR":                          "001df3d4a8622000", //系统发生错误
	"API_ROOMPOPULATION_ROOMUUID_ERROR":                           "001df3d4a8622000", //系统发生错误
	"API_SENDGIFT_DB_NO_DATA":                                     "001df3d4a8622000", //系统发生错误
	"API_SERVERROOMUSERLIST_JSON_ERROR":                           "001df3d4a8622000", //系统发生错误
	"API_SUBSCRIPTION_DB_NO_DATA":                                 "001df3d4a8622000", //系统发生错误
	"API_TOKENCHANGE_INSERT_USER_ERROR":                           "001df3d4a8622000", //系统发生错误
	"API_TOKENCHANGE_SELECT_USER_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_BLOCKROOMCHAT_DB_DELETE_ERROR":                       "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_DB_INSERT_ERROR":                       "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_GUEST":                                 "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_JSON_ERROR":                            "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_NOT_ADMIN":                             "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_ROOMUUID_ERROR":                        "001df3d907c22000", //禁言失敗
	"COMMAND_BLOCKROOMCHAT_TARGET_IS_ADMIN":                       "001df3d1c0822000", //无法封锁管理员
	"COMMAND_BLOCKROOMCHAT_TIME_ERROR":                            "001df3d907c22000", //禁言失敗
	"COMMAND_CREATEPRIVATEROOM_DIRTY_WORD":                        "002032ee3956c000", //房名含有敏感字
	"COMMAND_CREATEPRIVATEROOM_GUEST":                             "0020d9959d86c000", //未有权限
	"COMMAND_CREATEPRIVATEROOM_INSERT_DB_ERROR":                   "001df3d4a8622000", //系统发生错误
	"COMMAND_CREATEPRIVATEROOM_INSERT_PRIVATEGROUPUSERLIST_ERROR": "001df3d4a8622000", //系统发生错误
	"COMMAND_CREATEPRIVATEROOM_INSERT_USER_LIST_ERROR":            "001df3d4a8622000", //系统发生错误
	"COMMAND_CREATEPRIVATEROOM_JSON_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_CREATEPRIVATEROOM_NOT_ADMIN":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_CREATEPRIVATEROOM_ROOM_ICON_ERROR":                   "0020d9978186c000", //图片错误
	"COMMAND_CREATEPRIVATEROOM_ROOM_ICON_TOO_LARGE":               "0020d996a1e6c000", //图片太大
	"COMMAND_CREATEPRIVATEROOM_ROOM_NAME_ERROR":                   "0020d9989f66c000", //请重新命名
	"COMMAND_DISMISSROOM_GUEST":                                   "0020d9959d86c000", //未有权限
	"COMMAND_DISMISSROOM_JSON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_DISMISSROOM_NOT_ADMIN":                               "001df3d4a8622000", //系统发生错误
	"COMMAND_DISMISSROOM_ROOM_UUID_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDDELETE_FRIEND_DELETELIST_ERROR":                "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDDELETE_FRIEND_STATE_ERROR":                     "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDDELETE_GUEST":                                  "0020d9959d86c000", //未有权限
	"COMMAND_FRIENDDELETE_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_FRIEND_LIST_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_FRIEND_STATE_ERROR":                     "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_GUEST":                                  "0020d9959d86c000", //未有权限
	"COMMAND_FRIENDINVITE_HAVE_ALREADY_INVITED":                   "0020d9945216c000", //已经邀请好友
	"COMMAND_FRIENDINVITE_INSERT_DB_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_SELECT_DB_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_FRIENDINVITE_TARGET_IS_FRIEND":                       "0020d9935c96c000", //已经成为好友
	"COMMAND_FRIENDINVITE_USER_UUID_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_GETFRIENDLIST_GUEST":                                 "0020d9959d86c000", //未有权限
	"COMMAND_GETFRIENDLIST_JSON_ERROR":                            "001df3d4a8622000", //系统发生错误
	"COMMAND_GETFUNCMANAGEMENT_JSON_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_GETMEMBERLIST_JSON_ERROR":                            "001df3d4a8622000", //系统发生错误
	"COMMAND_GETMEMBERLIST_NOT_LIVEGROUP_LIST":                    "001df3d4a8622000", //系统发生错误
	"COMMAND_GETMEMBERLIST_ROOM_TYPE_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_GETMEMBERLIST_ROOM_UUID_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_GETNEWSIDETEXT_GUEST":                                "0020d9959d86c000", //未有权限
	"COMMAND_GETNEWSIDETEXT_JSON_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_GETNEWSIDETEXT_SEARCH_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_GETROOMHISTORY_ES_SEARCH_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_GETROOMHISTORY_JSON_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_GETROOMHISTORY_ROOM_UUID_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_GETSIDETEXTHISTORY_ES_SEARCH_ERROR":                  "001df3d4a8622000", //系统发生错误
	"COMMAND_GETSIDETEXTHISTORY_GUEST":                            "0020d9959d86c000", //未有权限
	"COMMAND_GETSIDETEXTHISTORY_JSON_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_KICKROOMUSER_GUEST":                                  "0020d9959d86c000", //未有权限
	"COMMAND_KICKROOMUSER_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_KICKROOMUSER_NOT_ADMIN":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_KICKROOMUSER_ROOM_TYPE_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_KICKROOMUSER_ROOM_UUID_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_MESSAGESEEN_GUEST":                                   "0020d9959d86c000", //未有权限
	"COMMAND_MESSAGESEEN_JSON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_MESSAGESEEN_TARGET_ROOM_TYPE_ERROR":                  "001df3d4a8622000", //系统发生错误
	"COMMAND_MESSAGESEEN_TARGET_ROOM_UUID_ERROR":                  "001df3d4a8622000", //系统发生错误
	"COMMAND_MESSAGESEEN_TARGET_SIDE_TEXT_UUID_ERROR":             "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERENTERROOMBATCH_IN_ROOM":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERENTERROOMBATCH_JSON_ERROR":                     "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERENTERROOMBATCH_ROOM_TYPE_NULL":                 "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERENTERROOMBATCH_ROOM_UUID_NULL":                 "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERENTERROOM_IN_ROOM":                             "001df3d65f422000", //直播间连线失败
	"COMMAND_PLAYERENTERROOM_JSON_ERROR":                          "001df3d65f422000", //直播间连线失败
	"COMMAND_PLAYERENTERROOM_ROOM_UUID_NULL":                      "001df3d65f422000", //直播间连线失败
	"COMMAND_PLAYEREXITROOM_JSON_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYEREXITROOM_ROOM_TYPE_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYEREXITROOM_ROOM_UUID_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYEREXITROOM_ROOM_UUID_NULL":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERLOGOUT_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERSENDMSG_CHAT_BLOCK":                            "001df3d59a122000", //用户已被封锁
	"COMMAND_PLAYERSENDMSG_ES_INSERT_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERSENDMSG_GUEST":                                 "001df3d82a722000", //请登入后再进行发言
	"COMMAND_PLAYERSENDMSG_JSON_ERROR":                            "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERSENDMSG_MSG_TOO_LONG":                          "001df3d2bae22000", //发言字数过多
	"COMMAND_PLAYERSENDMSG_NOT_IN_ROOM":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_PLAYERSENDMSG_SPEAK_CD":                              "001df3d394522000", //发言过于频繁
	"COMMAND_PROCLAMATIONSEARCH_JSON_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_GUEST":                                  "0020d9959d86c000", //未有权限
	"COMMAND_ROOMADMINADD_INSERT_ROLE_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_NOT_ADMIN":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_ROOMUUID_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_SELECT_ROLE_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_UPDATE_GROUP_ERROR":                     "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_UPDATE_ROLE_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMADMINADD_USER_ROLE_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_GUEST":                                  "0020d9959d86c000", //未有权限
	"COMMAND_ROOMINFOEDIT_INSERT_PRIVATEGROUPUSERLIST_ERROR":      "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_NOT_ADMIN":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_ROOM_ICON_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_ROOM_ICON_TOO_LARGE":                    "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_ROOM_NAME_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_ROOMINFOEDIT_ROOM_UUID_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTDELETE_DELETE_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTDELETE_GUEST":                                "0020d9959d86c000", //未有权限
	"COMMAND_SIDETEXTDELETE_JSON_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTDELETE_NOT_SIDETEXT":                         "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_CHATBLOCK_INSERT_ERROR":                 "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_DELETE_CHATBLOCK_ERROR":                 "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_ES_CHAT_HISTORY_INSERT_ERROR":           "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_ES_DIRTYWORD_HISTORY_INSERT_ERROR":      "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_GUEST":                                  "001df3d82a722000", //请登入后再进行发言
	"COMMAND_SIDETEXTSEND_INSERT_CHATTARGET_ERROR":                "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_JSON_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_MSG_TOO_LONG":                           "001df3d2bae22000", //发言字数过多
	"COMMAND_SIDETEXTSEND_QUERY_SIDETEXT_ERROR":                   "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_SELECT_UUID_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_SIDE_TEXT_YOURSELF":                     "001df3d4a8622000", //系统发生错误
	"COMMAND_SIDETEXTSEND_SPEAK_CD":                               "001df3d394522000", //发言过于频繁
	"COMMAND_TARGETADDROOMBATCH_INSERT_DB_ERROR":                  "001df3d4a8622000", //系统发生错误
	"COMMAND_TARGETADDROOMBATCH_JSON_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMAND_TARGETADDROOMBATCH_NOT_ADMIN":                        "001df3d4a8622000", //系统发生错误
	"COMMAND_TARGETADDROOMBATCH_ROOM_TYPE_ERROR":                  "001df3d4a8622000", //系统发生错误
	"COMMAND_TARGETADDROOMBATCH_TARGET_IN_ROOM":                   "001df3d4a8622000", //系统发生错误
	"COMMAND_TOKENCHANGE_DB_ERROR":                                "001df3d4a8622000", //系统发生错误
	"COMMAND_TOKENCHANGE_JSON_ERROR":                              "001df3d4a8622000", //系统发生错误
	"COMMAND_TOKENCHANGE_SIDETEXTMAP_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMON_ALIVECHECK_ERROR":                                     "001df3d4a8622000", //系统发生错误
	"COMMON_CHECKPLATFORMUSER_BODY_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMON_CHECKPLATFORMUSER_JSON_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMON_CHECKPLATFORMUSER_PLATFORM_UUID_ERROR":                "001df3d4a8622000", //系统发生错误
	"COMMON_CHECKPLATFORMUSER_REQUEST_ERROR":                      "001df3d4a8622000", //系统发生错误
	"COMMON_GETREDISFIRSTENTERROOM_ERROR":                         "001df3d4a8622000", //系统发生错误
	"COMMON_GETREDISROOMLASTMESSAGE_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMON_GETREDISROOMLASTSEEN_ERROR":                           "001df3d4a8622000", //系统发生错误
	"COMMON_GETREDISSIDETEXTLASTMESSAGE_ERROR":                    "001df3d4a8622000", //系统发生错误
	"COMMON_GETREDISSIDETEXTLASTSEEN_ERROR":                       "001df3d4a8622000", //系统发生错误
	"COMMON_HIERARCHYROOMINFOSEARCH_ROOM_READ_ERROR":              "001df3d4a8622000", //系统发生错误
	"COMMON_HIERARCHYROOMINFOSEARCH_ROOM_TYPE_NOT_WORD":           "001df3d4a8622000", //系统发生错误
	"COMMON_HIERARCHYTARGETINFOSEARCH_SELECT_USER_ERROR":          "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYBLOCKLIST_ERROR":                                 "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYFUNCTIONMANAGEMENT_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYPROCLAMATION_SELECT_DB_ERROR":                    "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYROOMINFO_ERROR":                                  "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYROOMINFO_ROOM_READ_ERROR":                        "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYROOMSINFO_ERROR":                                 "001df3d4a8622000", //系统发生错误
	"COMMON_QUERYSIDETEXTLIST_ERROR":                              "001df3d4a8622000", //系统发生错误
	"COMMON_REDISPUBDATA_REDIS_ERROR":                             "001df3d4a8622000", //系统发生错误
	"COMMON_SIDETEXTSEND_SIDETEXT_UUID_ERROR":                     "001df3d4a8622000", //系统发生错误
	"MAIN_ECHOHANDLER_ERROR":                                      "001df3d4a8622000", //系统发生错误
	"SHELL_BLOCKUSER_DB_DELETE_ERROR":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_DB_INSERT_ERROR":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_NOT_ADMIN":                                   "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_PARAMETER_ERROR":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_ROOM_UUID_ERROR":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_TARGET_IS_ADMIN":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_TIME_ERROR":                                  "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_USER_UUID_ERROR":                             "001df3d907c22000", //禁言失敗
	"SHELL_BLOCKUSER_UUID_NULL":                                   "001df3d907c22000", //禁言失敗
	"SHELL_LINKPROCLAMATION_CONTENT_TOOLONG":                      "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_DELETE_DB_ERROR":                      "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_INSERT_DB_ERROR":                      "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_ORDER_ERROR":                          "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_ROLE_ERROR":                           "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_SHELL_ERROR":                          "001df3da1f922000", //群连结设定失败
	"SHELL_LINKPROCLAMATION_TITLE_TOOLONG":                        "001df3da1f922000", //群连结设定失败
	"SHELL_NORMALPROCLAMATION_CONTENT_TOOLONG":                    "001df3d008622000", //公告设定失败
	"SHELL_NORMALPROCLAMATION_DELETE_DB_ERROR":                    "001df3d008622000", //公告设定失败
	"SHELL_NORMALPROCLAMATION_INSERT_DB_ERROR":                    "001df3d008622000", //公告设定失败
	"SHELL_NORMALPROCLAMATION_ROLE_ERROR":                         "001df3d008622000", //公告设定失败
	"SHELL_NORMALPROCLAMATION_SHELL_ERROR":                        "001df3d008622000", //公告设定失败
	"SHELL_NORMALPROCLAMATION_TITLE_TOOLONG":                      "001df3d008622000", //公告设定失败
	"SHELL_QUERYBLOCKLIST_JSON_ERROR":                             "001df3d742e22000", //指令发送失败
	"SHELL_QUERYBLOCKLIST_PARAMETER_ERROR":                        "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_JSON_ADMINSET_ERROR":                         "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_OCTOPUS_URL_ERROR":                           "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_PARAMETER_ERROR":                             "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_PATCH_ERROR":                                 "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_REQUEST_ERROR":                               "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_SELECT_ADMINSET_ERROR":                       "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_SELECT_USER_ERROR":                           "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_SHELL_ERROR":                                 "001df3d742e22000", //指令发送失败
	"SHELL_SHELLDEMO_TIME_ERROR":                                  "001df3d742e22000", //指令发送失败
	"SHELL_SHELL_SHELL_ERROR":                                     "001df3d742e22000", //指令发送失败
}

func Exception(msg string, userUuid string, err error) socket.Exception {

	if msg == "" {
		return socket.Exception{}
	}

	var code = Essyserrorlog(msg, userUuid, err)

	mutexClientMsgList.Lock()
	defer mutexClientMsgList.Unlock()
	clientMsg, ok := clientMsgList[msg]
	if !ok {
		clientMsg = "001df3d4a8622000" //系统发生错误
	}

	//待客端實作功能
	Mutexmutilangerrormsg.Lock()
	defer Mutexmutilangerrormsg.Unlock()
	clientMsg = Mutilangerrormsg["zh-CN"][clientMsg]

	if os.Getenv("environmentId") != "Online" {
		code = msg
	}

	return socket.Exception{Code: code, Message: clientMsg}
}

func Essyslog(msg string, loginUuid string, userUuid string) {

	if msg == "" {
		return
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sysErrorLog := Syserrorlog{Code: loginUuid, Useruuid: userUuid, Message: msg, Stamp: timeUnix}
	sysErrorLogJson, _ := json.Marshal(sysErrorLog)
	Esinsert(os.Getenv("sysLog"), string(sysErrorLogJson))

	return
}

func Essyserrorlog(msg string, userUuid string, err error) string {

	if msg == "" {
		return ""
	}

	var errType string
	if err != nil {
		errType = err.Error()
	}

	var code = bkdrHash(msg, 36, 5)

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sysErrorLog := Syserrorlog{Useruuid: userUuid, Code: code, Message: msg, Error: errType, Stamp: timeUnix}
	sysErrorLogJson, _ := json.Marshal(sysErrorLog)
	Esinsert(os.Getenv("sysErrorLog"), string(sysErrorLogJson))

	return code
}

// BKDR-Hash
func bkdrHash(text string, base int64, length int) string {
	const seed = int64(131)

	var divisor int64 = 1
	for i := 0; i < length; i++ {
		divisor *= base
	}

	var hash = int64(0)
	var textByte = []byte(text)
	var textLength = len(textByte)
	for i := 0; i < textLength; i++ {
		hash = (hash*seed + int64(textByte[i])) % divisor
	}

	var code = strconv.FormatInt(hash, int(base))
	var codeLength = len(code)
	if codeLength < length {
		code = strings.Repeat("0", length-codeLength) + code
	}

	return code
}
