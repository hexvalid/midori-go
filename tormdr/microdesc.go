package main

import (
	"bufio"
	"errors"
	"github.com/fatih/color"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	exitNodeCache            []ExitNode
	exitNodeCachePopulated   bool
	exitNodeCacheUpdateMutex sync.Mutex
)

type ExitNode struct {
	ip        string
	bandwidth int
	fast      bool
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

func populateExitNodeCache(cacheDir string) (err error) {
	exitNodeCacheUpdateMutex.Lock()
	if !exitNodeCachePopulated {
		log.SInfo("ALL", "Populating Exit Node Cache...")
		exitNodeCache, err = parseExitNodes(path.Join(cacheDir, "cached-microdesc-consensus"))
		totalbandwidth := 0
		for i := 0; i < len(exitNodeCache); i++ {
			totalbandwidth += exitNodeCache[i].bandwidth
		}
		avrbandwidth := totalbandwidth / len(exitNodeCache)
		log.SInfo("ALL", "%s Exit Node learned. Avarage bandwidth: %s.",
			color.YellowString(strconv.Itoa(len(exitNodeCache))),
			color.YellowString(strconv.Itoa(avrbandwidth)))
		exitNodeCachePopulated = true
	}
	exitNodeCacheUpdateMutex.Unlock()
	return
}

func FindExitNode(excludedIPs []string, minBandwidth int, fast, stable, valid bool) (string, error) {
	if exitNodeCachePopulated && len(exitNodeCache) > 1 {
		cache := make([]ExitNode, len(exitNodeCache))
		copy(cache, exitNodeCache)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(cache), func(i, j int) { cache[i], cache[j] = cache[j], cache[i] })
		for i := 0; i < len(cache); i++ {
			if cache[i].bandwidth >= minBandwidth &&
				reqOptional(fast, cache[i].fast) &&
				reqOptional(stable, cache[i].stable) &&
				reqOptional(valid, cache[i].valid) &&
				!containsIn(cache[i].ip, excludedIPs) {
				return cache[i].ip, nil
			}
		}
		return "", errors.New("suitable Exit Node not found")
	} else {
		return "", errors.New("cache of Exit Node not populated")
	}
}

func reqOptional(req, optional bool) bool {
	if req {
		return optional
	} else {
		return true
	}
}
func containsIn(ip string, list []string) bool {
	if list != nil {
		for i := range list {
			if list[i] == ip {
				return true
			}
		}
	}
	return false
}
