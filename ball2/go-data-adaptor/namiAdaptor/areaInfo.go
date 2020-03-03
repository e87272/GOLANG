package namiAdaptor

import (
	"log"

	"../commonFunc"
	"../data"
	"../data/continentData"
	"../external/database"
)

func areaInfoMap(areas map[string]data.NamiAreaInfo) {
	areaList := []data.NamiAreaInfo{}
	for _, areaInfo := range areas {
		areaList = append(areaList, areaInfo)
	}
	areaInfoUpdate(areaList)
}

func areaInfoUpdate(areas []data.NamiAreaInfo) error {

	_, err := continentData.SearchContinentUuidByNami(0)

	//撈不到資料預塞洲編號0
	if err != nil {

		if err != database.ErrNoRows {
			log.Printf("areaInfoUpdate select err : %+v\n", err)
			return err
		}

		continentUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `continent` (`uuid` ,`name`) VALUES (?, ?)",
			continentUuid, "未知",
		)

		if err != nil {
			log.Printf("areaInfoUpdate INSERT continent err : %+v\n", err)
			return err
		}

		corporationUuid := commonFunc.GetUuid()
		_, err := database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "continent", continentUuid, "nami", 0,
		)

		if err != nil {
			log.Printf("areaInfoUpdate INSERT corporation_id err : %+v\n", err)
			return err
		}
	}

	log.Printf("areaInfoUpdate areas len : %+v\n", len(areas))
	for _, areaInfo := range areas {
		log.Printf("areaInfoUpdate areaInfo : %+v\n", areaInfo)
		continentUuid, err := continentData.SearchContinentUuidByNami(areaInfo.Id)

		if err != nil {
			if err != database.ErrNoRows {
				log.Printf("areaInfoUpdate select err : %+v\n", err)
				return err
			}

			continentUuid = commonFunc.GetUuid()
			_, err := database.Exec("INSERT INTO `continent` (`uuid` ,`name`) VALUES (?, ?)",
				continentUuid, areaInfo.Name_zh,
			)
			if err != nil {
				log.Printf("areaInfoUpdate INSERT DB err : %+v\n", err)
				return err
			}

			corporationUuid := commonFunc.GetUuid()
			_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
				corporationUuid, "continent", continentUuid, "nami", areaInfo.Id,
			)

			if err != nil {
				log.Printf("areaInfoUpdate INSERT DB err : %+v\n", err)
				return err
			}
			continue
		}

		_, err = database.Exec("UPDATE `continent` SET `name` = ?   WHERE `uuid` = ?", areaInfo.Name_zh, continentUuid)

		if err != nil {
			log.Printf("areaInfoUpdate update DB err : %+v\n", err)
			return err
		}
	}
	return nil
}
