package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("vim-go")
	file, _ := os.OpenFile("ipaddress.txt", os.O_RDONLY, 0666)

	scanner := bufio.NewScanner(file)
	image := "diamanti-cx-9.9.2-60.x86_64.rpm"
	for scanner.Scan() {
		text := scanner.Text()
		newtext := "diamanti@" + text + ":/home/diamanti"
		fmt.Println("newtext is ", newtext)
		cmd := exec.Command("sshpass", "-p", "diamanti", "scp", image, newtext)
		stderr, err := cmd.CombinedOutput()
		fmt.Printf("error is %s and str:%s\n", err.Error(), string(stderr))
	}

}
