package mj

import (
	"chatplus/utils"
	"regexp"
)

func extractProgress(input string) int {
	pattern := `\((\d+)\%\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return utils.IntValue(matches[1], 0)
	}
	return 100
}

func extractPercentage(input string) int {
	pattern := `(\d+)\%`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return utils.IntValue(matches[1], 0)
	}
	return 100
}
