package ipData

import (
	"server/common"
	"server/external/database"
)

func QueryWhiteListDB() bool {

	rows, err := database.Query("SELECT `ip` FROM `whiteList` WHERE  1 ")

	if err != nil {
		common.SysErrorLog(map[string]interface{}{
			"name": "QueryWhiteListDB select err",
		}, err)
		return false
	}

	for rows.Next() {
		ip := ""
		err := rows.Scan(&ip)
		if err != nil {
			common.SysErrorLog(map[string]interface{}{
				"name": "QueryWhiteListDB err",
			}, err)
			return false
		}
		SetWhiteList(ip)
	}
	return true
}
