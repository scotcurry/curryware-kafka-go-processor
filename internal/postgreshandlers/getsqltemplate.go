package postgreshandlers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func GetSqlTemplate(templateName string) string {

	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current path")
		panic(err)
	}
	fmt.Println(currentPath)
	pathToFile := filepath.Join(currentPath, "/sqltemplates/sqltemplate.txt")
	fileData, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println("Error opening sqltemplate.txt")
		panic(err)
	}
	defer func(fileData *os.File) {
		err := fileData.Close()
		if err != nil {

		}
	}(fileData)

	sqlCommand := ""
	scanner := bufio.NewScanner(fileData)
	for scanner.Scan() {
		line := scanner.Text()
		matchFound, foundSqlCommand := checkForTemplateMatch(line, templateName)
		if matchFound {
			sqlCommand = foundSqlCommand
		}
	}
	return sqlCommand
}

func checkForTemplateMatch(line string, templateName string) (bool, string) {

	fileFound := false
	sqlCommand := ""
	regexPattern := `^\s*(\S+)\s+(INSERT INTO .*)$`
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindAllStringSubmatch(line, -1)
	if matches[0][1] == templateName {
		fileFound = true
		sqlCommand = matches[0][2]
	}
	return fileFound, sqlCommand
}
