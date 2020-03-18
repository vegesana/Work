package debug

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func processFwdCounters(path string, servername string) error {
	/* Check err */
	f, _ := os.Open(path)
	r := regexp.MustCompile(`[ ]+`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "errors") || strings.Contains(line, "drop") {
			val := r.Split(line, -1)

			if len(val) > 0 {
				for i, value := range val[1:] {
					if value != "0000000000000000" {
						errstr := fmt.Sprintf("Server:%s:If:eth%d Error:%s,Count:%s\n",
							servername, i, strings.TrimSpace(val[0]), value)
						writeToDb(MyError{servername, errstr})
					}
				}

			}
		}
	}
	return nil
}
