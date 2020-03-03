package syncApi

import (
	"encoding/json"
	"os"

	"../commonFunc"
)

type jsObj map[string]interface{}

type testCase struct {
	Method string
	Title  jsObj
	Input  jsObj
	Output interface{}
}

var testSet = []testCase{}

func appendTestCase(data testCase) {
	testSet = append(testSet, data)
}

func createTestFile(fileName string) {

	file, err := os.Create("./unitTest/" + fileName)
	if err != nil {
		commonFunc.EsSysLog("createTestFile", err.Error())
		return
	}
	defer file.Close()

	content, _ := json.MarshalIndent(testSet, "", "\t")
	file.Write([]byte("var testSet = "))
	file.Write(content)
}
