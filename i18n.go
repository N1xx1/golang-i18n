package i18n

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	if scannerErr := fileScanner.Err(); scannerErr != nil {
		return nil, scannerErr
	}

	return func(key string, args ...interface{}) string {
		if translation, ok := translations[key]; ok {
			return fmt.Sprintf(translation, args...)
		}
		return key
	}, nil
}

func SetGlobalT(f TranslationFunction) {
	T = f
}

func GlobalTfunc(translationFile string) error {
	t, err := Tfunc(translationFile)
	if err != nil {
		return err
	}

	SetGlobalT(t)
	return nil
}

func parseFile(scanner *bufio.Scanner) (map[string]string, error) {
	translations := make(map[string]string)
	definitionRegexp := regexp.MustCompile(`^([\d\w\-_]+)\s*=\s*(".*")\s*(?:\#.*)?$`)
	emptyLineRegexp := regexp.MustCompile(`^(|\s*(\#.*)?)$`)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()

		if emptyLineRegexp.MatchString(line) {
			continue
		}

		matches := definitionRegexp.FindStringSubmatch(line)
		if len(matches) != 3 {
			return nil, errors.New(fmt.Sprintf("Malformed line :%d (%q)", lineNumber, line))
		}

		unquoted, err := strconv.Unquote(matches[2])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Malformed string :%d (%q)", lineNumber, line))
		}
		if _, ok := translations[matches[1]]; ok {
			return nil, errors.New(fmt.Sprintf("Multiple definitions of key %q", matches[1]))
		}

		translations[matches[1]] = unquoted
	}
	return translations, nil
}

func loadFile(filePath string) (*os.File, *bufio.Scanner, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	return file, bufio.NewScanner(file), nil
}
