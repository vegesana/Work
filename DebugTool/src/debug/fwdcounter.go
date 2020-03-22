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
	fmt.Println("Inside processFwdCounters", servername)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "errors") || strings.Contains(line, "drop") {
			val := r.Split(line, -1)

			if len(val) > 0 {
				for i, value := range val[1:] {
					if value != "0000000000000000" {
						errstr := fmt.Sprintf("If:eth%d Error:%s,Count:%s\n",
							i, strings.TrimSpace(val[0]), value)

						fmt.Println("writetoDB processFwdCounters", servername)
						writeToDb(MyError{servername, errstr})
					}
				}

			}
		}
	}

	fmt.Println("End processFwdCounters", servername)
	return nil
}
