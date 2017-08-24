package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	ELEMENT_DELIM = "\t"
	LV_DELIM      = ":"
)

func ltsvToJon(line string) ([]byte, error) {
	line = strings.TrimSpace(line)
	lvs := strings.Split(line, ELEMENT_DELIM)
	lvMap := map[string]string{}
	for _, lvString := range lvs {
		lv := strings.SplitN(lvString, LV_DELIM, 2)
		lvMap[lv[0]] = lv[1]
	}
	return json.Marshal(lvMap)
}

func readLines() int {
	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				return 0
			}
			fmt.Printf("Failed: %#v\n", err)
			return 1
		}
		lineJson, err := ltsvToJon(line)
		if err != nil {
			fmt.Printf("Failed: %#v\n", err)
			return 1
		}
		fmt.Println(string(lineJson))
	}
}

func main() {
	exitCode := readLines()
	os.Exit(exitCode)
}
