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
	INDENT_STRING = "    "
)

var indentOption = false

// Parse ltsv line and convert to json
func ltsvToJon(line string) ([]byte, error) {
	lvs := strings.Split(line, ELEMENT_DELIM)
	lvMap := map[string]string{}
	for _, lvString := range lvs {
		lv := strings.SplitN(lvString, LV_DELIM, 2)
		lvMap[lv[0]] = lv[1]
	}

	if indentOption {
		return json.MarshalIndent(lvMap, "", INDENT_STRING)
	} else {
		return json.Marshal(lvMap)
	}
}

// Read ltsv lines from given file
// parsed each line and write json to stdout
func readLines(f *os.File) int {
	stdin := bufio.NewReader(f)
	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				return 0
			}
			fmt.Printf("Error: %v\n", err)
			return 1
		}
		line = strings.TrimSpace(line)

		// ignore blank line
		if line == "" {
			continue
		}

		// convert ltsv to json
		lineJson, err := ltsvToJon(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return 1
		}
		fmt.Println(string(lineJson))
	}
}

// show help message and exit
func showHelpAndExit() {
	msg := `Usage of %s: [options] [FILE_NAME]
Options:
    -i --indent  Enable json indent
    -h --help    Show this message`

	fmt.Fprintf(os.Stderr, msg, os.Args[0])
	os.Exit(1)
}

// read args and options
func procArgs() (*os.File, bool, error) {
	indent := false
	toRead := os.Stdin
	for _, arg := range os.Args[1:] {
		// check help option
		if arg == "-h" || arg == "--help" {
			showHelpAndExit()
		}

		// check indent option
		if arg == "-i" || arg == "--indent" {
			indent = true
			continue
		}

		// check file name arg
		if !strings.HasPrefix("-", arg) {
			r, err := os.Open(arg)
			if err != nil {
				return toRead, indent, err
			}
			toRead = r
		}
	}
	return toRead, indent, nil
}

func main() {
	toRead, indent, err := procArgs()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	defer toRead.Close()
	indentOption = indent
	exitCode := readLines(toRead)
	os.Exit(exitCode)
}
