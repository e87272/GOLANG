package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	fileNameList := []string{}
	scanFolder("server", &fileNameList)
	errorCodeMap := map[string]string{}
	errorMessageList := []string{}
	re := regexp.MustCompile(`(?:return |Exception\(|Essyserrorlog\()"([A-Z_]+)",`)
	for _, fileName := range fileNameList {
		file, _ := os.Open(fileName)
		textByte, _ := ioutil.ReadAll(file)
		text := string(textByte)
		matchList := re.FindAllStringSubmatch(text, -1)
		for _, match := range matchList {
			message := match[1]
			_, ok := errorCodeMap[message]
			if !ok {
				errorCodeMap[message] = bkdrHash(message, 36, 5)
				errorMessageList = append(errorMessageList, message)
			}
		}
	}
	sort.Strings(errorMessageList)
	output := "\n"
	for _, message := range errorMessageList {
		output += message + " " + errorCodeMap[message] + "\n"
	}
	fmt.Print(output)
}

func scanFolder(folder string, nameList *[]string) {
	fileList, _ := ioutil.ReadDir(folder)
	for _, file := range fileList {
		name := folder + "/" + file.Name()
		if file.IsDir() {
			scanFolder(name, nameList)
		} else {
			*nameList = append(*nameList, name)
		}
	}
}

// BKDR-Hash
func bkdrHash(text string, base int64, length int) string {
	const seed = int64(131)

	var divisor = int64(1)
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
