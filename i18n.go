package i18n

import (
	"strconv"
	"fmt"
	"errors"
	"bufio"
	"os"
	"regexp"
)

type TranslationFunction func(string, ...interface{}) string

var T TranslationFunction

func Tfunc(translationFile string) (TranslationFunction, error) {
	file, fileScanner, err := loadFile(translationFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	translations, err := parseFile(fileScanner)
	if err != nil {
		return nil, err
	}
	if scannerErr = fileScanner.Err(); scannerErr != nil {
		return nil, scannerErr
	}
	
	return func(key string, args ...interface{}) string {
		if translation, ok := translations[key]; ok {
			return fmt.Sprintf(translation, args)
		}
		return key
	}
}

func parseFile(scanner bufio.Scanner) (map[string]string, error) {
	translations := make(map[string]string)
	r, _ := regexp.Compile(`([\d\w\-_])\s*=\s*(".*")`)
	lineNumber := 1
	
	for scanner.Scan() {
		line := scanner.Text()
		matches := r.FindStringSubmatch(line)
		if len(matches) != 2 {
			return errors.New(fmt.Sprintf("Malformed line %d", lineNumber))
		}
		
		unquoted, err := strconv.Unquote(matches[1])
		if err != nil {
			return errors.New(fmt.Sprintf("Malformed string %d", lineNumber))
		}
		if _, ok := translations[matches[0]]; ok {
			return errors.New(fmt.Sprintf("Multiple definitions of key %q", matches[0]))
		}
		
		translations[matches[0]] = unquoted
		lineNumber += 1
	}
	
	return translations, nil
}

func loadFile(filePath string) (os.File, bufio.Scanner, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	
	return file, bufio.NewScanner(file), nil
}