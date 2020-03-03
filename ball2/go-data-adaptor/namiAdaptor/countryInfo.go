package namiAdaptor

import (
	"log"

	"../commonFunc"
	"../data"
	"../data/continentData"
	"../data/countryData"
	"../external/database"
)

func countryInfoMap(countrys map[string]data.NamiCountryInfo) {
	countryList := []data.NamiCountryInfo{}
	for _, countryInfo := range countrys {
		countryList = append(countryList, countryInfo)
	}
	countryInfoUpdate(countryList)
}

func countryInfoUpdate(countrys []data.NamiCountryInfo) error {

	continentUuid, err := continentData.SearchContinentUuidByNami(0)

	if err != nil {
		log.Printf("countryInfoUpdate select continent err : %+v\n", err)
		return err
	}
	_, err = countryData.SearchCountryUuidByNami(0)
	//撈不到資料預塞國家編號0
	if err != nil {
		if err != database.ErrNoRows {
			log.Printf("countryInfoUpdate select err : %+v\n", err)
			return err
		}

		countryUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `country` (`uuid` ,`name` ,`continent_uuid`) VALUES (?, ?, ?)",
			countryUuid, "世界", continentUuid,
		)

		if err != nil {
			log.Printf("countryInfoUpdate insert country err : %+v\n", err)
			return err
		}

		corporationUuid := commonFunc.GetUuid()
		_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
			corporationUuid, "country", countryUuid, "nami", 0,
		)

		if err != nil {
			log.Printf("countryInfoUpdate insert corporation_id err : %+v\n", err)
			return err
		}

	}

	for _, countryInfo := range countrys {
		// log.Printf("countryInfoUpdate countryInfo : %+v\n", countryInfo)

		continentUuid, err := continentData.SearchContinentUuidByNami(countryInfo.Area_id)
		if err != nil {
			log.Printf("countryInfoUpdate select continent err : %+v\n", err)
			return err
		}

		countryUuid, err := countryData.SearchCountryUuidByNami(countryInfo.Id)
		if err != nil {
			if err != database.ErrNoRows {
				log.Printf("countryInfoUpdate select err : %+v\n", err)
				return err
			}
			countryUuid = commonFunc.GetUuid()
			_, err := database.Exec("INSERT INTO `country` (`uuid` ,`name` ,`continent_uuid`) VALUES (?, ?, ?)",
				countryUuid, countryInfo.Name_zh, continentUuid,
			)

			corporationUuid := commonFunc.GetUuid()
			_, err = database.Exec("INSERT INTO `corporation_id` (`uuid` ,`type` ,`my_id` , `source_name` , `source_id`) VALUES (?, ?, ?, ?, ?)",
				corporationUuid, "country", countryUuid, "nami", countryInfo.Id,
			)

			if err != nil {
				log.Printf("countryInfoUpdate update DB err : %+v\n", err)
				return err
			}

			continue

		}

		if countryInfo.Logo != "" {
			_, err = commonFunc.PostCdnUploadLink(countryInfo.Logo, countryUuid+".png", "/countryIcon/")
			if err != nil {
				log.Printf("CountryInfoUpdate err : %+v\n", err)
			}
		}

		_, err = database.Exec("UPDATE `country` SET `name` = ? , `continent_uuid` = ? WHERE `uuid` = ?",
			countryInfo.Name_zh, continentUuid, countryUuid)

		if err != nil {
			log.Printf("CountryInfoUpdate update DB err : %+v\n", err)
			return err
		}
	}
	return nil
}
