package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var defaultArgs = []string{
	"RunAsDaemon", "0",
	"ClientOnly", "1",
	"AvoidDiskWrites", "1",
	"FetchHidServDescriptors", "0",
	"FetchServerDescriptors", "1",
	"FetchUselessDescriptors", "0",
	"HardwareAccel", "1",
	"KeepalivePeriod", "30",
	"UseEntryGuards", "0",
	"NumEntryGuards", "0",
	"UseGuardFraction", "0",
	"DownloadExtraInfo", "0",
	"UseMicrodescriptors", "1",
	"ClientUseIPv4", "1",
	"ClientUseIPv6", "0",
	"DirCache", "0",
	"NewCircuitPeriod", "90 days",
	"MaxCircuitDirtiness", "30 days",
	"EnforceDistinctSubnets", "0",
	//"ExitNodes", "51.75.52.118",
}

const (
	socksPortStart   = 20000
	controlPortStart = 40000
)

type TorMDR struct {
	cmd        *exec.Cmd
	transport  *http.Transport
	dataDir    string
	bootStatus string
	mutex      *sync.Mutex
}

func main() {
	logger := logger.NewLog("tor", color.FgMagenta)
	logger.Info("Hello from logger")
	var args = []string{""}
	args = append(args, defaultArgs...)
	args = append(args, "SocksPort", "20000")
	args = append(args, "ControlPort", "40000")
	args = append(args, "DataDirectory", "/tmp/arya7")
	args = append(args, "CacheDirectory", "/tmp/arya7/cache")

	args = append(args, "Socks5Proxy", "40.70.243.118:32416")
	args = append(args, "Socks5ProxyUserName", "e4cf6e290c0cf8ae8fb91fcf818e1e40")
	args = append(args, "Socks5ProxyPassword", "a565ab1f3802afbf4d07c1674069d813")

	//ExitNodes
	cmd := &exec.Cmd{
		Path: "/home/hexvalid/tormdr/tormdr",
		Args: args,
	}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}
	cmd.Start()

	scanner := bufio.NewScanner(cmdReader)

	if scanner.Scan() {
		fmt.Println("Version:", scanner.Text())
	}

	var bootStatus string
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Bootstrapped") {
			re := regexp.MustCompile(`\((.*?)\)`)
			bootStatus = re.FindStringSubmatch(line)[1]
			if bootStatus == "done" {
				break
			}
		} else {
			//fmt.Printf("\t > %s\n", line)
		}
	}

	fmt.Println("ended")
}
