package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	str := " very berry : 2 "
	re := regexp.MustCompile(`\s*([\w\s]+)\s*:\s*.*`)
	// matches := re.FindAllStringSubmatch(str, -1)
	matches := re.FindStringSubmatch(str)
	fmt.Println("len Matches :", len(matches))
	replacestr := `\s*` + strings.TrimSpace(matches[1]) + `\s*:\s*(.*)\s*`
	fmt.Println("Replce Str:", replacestr)
}
