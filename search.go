package rde

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

type RegExp struct {
	Brand         string
	RegExpression []string
}

func UnmarshalRegex() ([]RegExp, error) {
	configPath, err := GetConfigPath("regex.json")
	if err != nil {
		return nil, err
	}
	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("faild open file with regular expressions: %v", err)
	}

	var regExps map[string][]string
	err = json.Unmarshal(jsonFile, &regExps)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal regular expressions: %v", err)
	}

	var regExpList []RegExp
	for brand, patterns := range regExps {
		regExpList = append(regExpList, RegExp{
			Brand:         brand,
			RegExpression: patterns,
		})
	}
	return regExpList, nil
}

func Match(data []string, regExp []RegExp) []string {
	for _, reg := range regExp {
		for _, pattern := range reg.RegExpression {
			expr, err := regexp.Compile(pattern)
			if err != nil {
				log.Printf("WARNING: %s is not valid regular expression. Skiping.", pattern)
				continue
			}
			for _, entry := range data {
				if expr.Match([]byte(entry)) {
					log.Printf("%s matched %s", pattern, entry)
				}
			}
		}
	}
	return []string{}
}
