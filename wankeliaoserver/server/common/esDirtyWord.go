package common

import (
	"encoding/json"
	"io/ioutil"
)

func loaddirtyword(IndexName string) {

	// 可写方式打开文件
	// file, err := os.OpenFile(
	// 	"data/dirtyword.txt",
	// 	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	// 	0666,
	// )

	// log.Printf("loaddirtyword")

	dirtywordData, err := ioutil.ReadFile("data/dirtyword.json")
	if err != nil {
		// log.Printf("找不到dirtyword.json")
	}

	dirtywordDataObj := [][]interface{}{}

	err = json.Unmarshal(dirtywordData, &dirtywordDataObj)
	if err != nil {
		// log.Printf("err %+v \n", err)
	}

	for _, v := range dirtywordDataObj {
		// log.Printf("%v : %+v\n", k, v[0])

		dirtyword := make(map[string]interface{})
		dirtyword["dirtyword"] = v[0]

		// 写字节到文件中
		// byteSlice := []byte(dirtyword["dirtyword"].(string) + "\n")
		// file.Write(byteSlice)

		// Index a second tweet (by string)
		dirtywordJson, _ := json.Marshal(dirtyword)

		Esinsert("dirtyword", string(dirtywordJson))
	}

}
