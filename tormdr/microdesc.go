package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	exitNodeCachePopulated   bool
	exitNodeCache            []ExitNode
	exitNodeCacheUpdateMutex sync.Mutex
)

type ExitNode struct {
	ip        string
	bandwidth int
	fast      bool
	guard     bool
	stable    bool
	valid     bool
}

func parseExitNodes(path string) (exitNodeList []ExitNode, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() && scanner.Text() == "network-status-version 3 microdesc" {
		for scanner.Scan() {
			lineR := scanner.Text()
			if lineR[0] == 'r' && lineR[1] == ' ' {
				ip := strings.Split(lineR, " ")[5]
				scanner.Scan()
				lineM := scanner.Text()
				if lineM[0] == 'a' {
					scanner.Scan()
				}
				scanner.Scan()
				lineS := scanner.Text()
				if !(lineS[0] == 's' && lineS[1] == ' ') {
					return nil, errors.New("unexcepted S line")
				}
				scanner.Scan()
				scanner.Scan()
				scanner.Scan()
				lineW := scanner.Text()
				if !(lineW[0] == 'w' && lineW[1] == ' ' && lineW[2] == 'B') {
					return nil, errors.New("unexcepted W line")
				}
				bandwidth, _ := strconv.Atoi(strings.Split(lineW, "=")[1])
				if strings.Contains(lineS, "Exit") && strings.Contains(lineS, "Running") &&
					!strings.Contains(lineS, "BadExit") {
					exitNodeList = append(exitNodeList, ExitNode{
						ip:        ip,
						bandwidth: bandwidth,
						fast:      strings.Contains(lineS, "Fast"),
						guard:     strings.Contains(lineS, "Guard"),
						stable:    strings.Contains(lineS, "Stable"),
						valid:     strings.Contains(lineS, "Valid"),
					})
				}
			}
		}
	} else {
		return nil, errors.New("unknown microdesc file")
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return
}
