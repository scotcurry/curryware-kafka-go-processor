package postgreshandlers

import (
	"bufio"
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
	"os"
	"regexp"
)

func GetSqlTemplate(templateName string) string {

	pathToFile := "/app/internal/postgreshandlers/sqltemplates/sqltemplate.txt"
	fileData, err := os.Open(pathToFile)
	if err != nil {
		logger.LogError(fmt.Sprintf("Path to file: %s", pathToFile))
		return ""
	}
	defer func(fileData *os.File) {
		err := fileData.Close()
		if err != nil {
			logger.LogInfo(fmt.Sprintf("Closing sqltemplate.txt %s", pathToFile))
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
