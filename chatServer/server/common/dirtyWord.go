package common

import (
	
	"regexp"
	"strings"
	"unicode/utf8"

	"server/database"
)

func Querydirtyword() {
	Mutexdirtywordlist.Lock()
	defer Mutexdirtywordlist.Unlock()
	rows, err := database.Query(
		"select dirtyWordUuid,dirtyWord from dirtyWord",
	)

	if err != nil {
		// log.Printf("Querydirtyword err : %+v\n", err)
	}
	var dirtyWord string
	var dirtyWordUuid string

	Dirtywordlist = make(map[string]string)

	for rows.Next() {

		rows.Scan(&dirtyWordUuid, &dirtyWord)
		Dirtywordlist[dirtyWordUuid] = dirtyWord

	}
	rows.Close()

	// log.Printf("Dirtywordlist : %+v\n", Dirtywordlist)
}

func Matchdirtyword(message string, maxWordLength int) (bool, string) {

	Mutexdirtywordlist.Lock()
	defer Mutexdirtywordlist.Unlock()
	var isDirtyWord = false
	for _, dirtyWordValue := range Dirtywordlist {
		var re = regexp.MustCompile(dirtyWordValue)
		for i := 0; i < maxWordLength; i++ {
			match := re.FindAllStringSubmatchIndex(message, -1)
			if len(match) == 0 {
				break
			}
			isDirtyWord = true
			for _, subMatch := range match {
				for matchIndex := 2; matchIndex < len(subMatch); matchIndex += 2 {
					var startIndex = subMatch[matchIndex]
					var endIndex = subMatch[matchIndex+1]
					var wordLength = utf8.RuneCountInString(message[startIndex:endIndex])
					var paddingLength = endIndex - startIndex - wordLength
					message = message[0:startIndex] + strings.Repeat("*", wordLength) + strings.Repeat("\a", paddingLength) + message[endIndex:]
				}
			}
		}
	}
	if isDirtyWord {
		message = regexp.MustCompile("\a+").ReplaceAllString(message, "")
	}
	return isDirtyWord, message
}
