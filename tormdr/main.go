package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"midori-go/logger"
	"net/http"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultArgs = []string{"",
	"RunAsDaemon", "0",
	"ClientOnly", "1",
	"AvoidDiskWrites", "1",
	"FetchHidServDescriptors", "0",
	"FetchServerDescriptors", "1",
	"FetchUselessDescriptors", "0",
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

var log = logger.NewLog("TorMDR", color.FgMagenta)

const (
	socksPortStart   int = 20000
	controlPortStart int = 40000
)

type TorMDRConfig struct {
	TorMDRBinaryPath    string
	DataDirectory       string
	HardwareAccel       bool
	KeepalivePeriod     int
	UseSocks5Proxy      bool
	Socks5ProxyAddress  string
	Socks5ProxyUserName string
	Socks5ProxyPassword string
	UseObfs4Proxy       bool
	Obfs4ProxyPath      string
	Obfs4Bridges        []string
}

type TorMDR struct {
	cmd        *exec.Cmd
	no         int
	transport  *http.Transport
	dataDir    string
	bootStatus string
	mutex      *sync.Mutex
	stdoutPipe io.ReadCloser
}

func NewTorMDR(no int, cfg *TorMDRConfig) *TorMDR {
	tormdr := TorMDR{no: no, cmd: &exec.Cmd{}}
	tormdr.cmd.Path = cfg.TorMDRBinaryPath
	tormdr.cmd.Args = append(tormdr.cmd.Args, defaultArgs...)
	tormdr.cmd.Args = append(tormdr.cmd.Args, "SocksPort", strconv.Itoa(socksPortStart+no))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "ControlPort", strconv.Itoa(controlPortStart+no))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "DataDirectory", path.Join(cfg.DataDirectory, strconv.Itoa(no)))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "CacheDirectory", path.Join(cfg.DataDirectory, strconv.Itoa(no), "cache"))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "KeepalivePeriod", strconv.Itoa(cfg.KeepalivePeriod))
	if cfg.HardwareAccel {
		tormdr.cmd.Args = append(tormdr.cmd.Args, "HardwareAccel", "1")
	} else {
		tormdr.cmd.Args = append(tormdr.cmd.Args, "HardwareAccel", "0")
	}
	if cfg.UseSocks5Proxy {
		tormdr.cmd.Args = append(tormdr.cmd.Args, "Socks5Proxy", cfg.Socks5ProxyAddress)
		tormdr.cmd.Args = append(tormdr.cmd.Args, "Socks5ProxyUserName", cfg.Socks5ProxyUserName)
		tormdr.cmd.Args = append(tormdr.cmd.Args, "Socks5ProxyPassword", cfg.Socks5ProxyPassword)
	}
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Initialized")
	return &tormdr
}

func (tormdr *TorMDR) Start() (err error) {
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Starting...")
	tormdr.stdoutPipe, err = tormdr.cmd.StdoutPipe()
	if err != nil {
		return errors.New("can't access stdoutpipe of tormdr")
	}

	tormdr.cmd.Start()
	scanner := bufio.NewScanner(tormdr.stdoutPipe)
	if scanner.Scan() {
		fmt.Println("Version:", scanner.Text())
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Bootstrapped") {
			re := regexp.MustCompile(`\((.*?)\)`)
			tormdr.bootStatus = re.FindStringSubmatch(line)[1]
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Boot status: %s", tormdr.bootStatus)

			if tormdr.bootStatus == "done" {
				break
			}
		} else {
			//todo: warning & ERROR
			fmt.Printf("\t > %s\n", line)
		}
	}

	//todo: open control port and auth
	return nil
}

func main() {

	/*	args = append(args, "Socks5Proxy", "40.70.243.118:32416")
		args = append(args, "Socks5ProxyUserName", "e4cf6e290c0cf8ae8fb91fcf818e1e40")
		args = append(args, "Socks5ProxyPassword", "a565ab1f3802afbf4d07c1674069d813")*/

	fmt.Println(path.Join("/tmp", strconv.Itoa(1)))
	tormdr := NewTorMDR(1, &TorMDRConfig{
		TorMDRBinaryPath: "/home/anon/tormdr/tormdr",
		DataDirectory:    "/tmp",
		KeepalivePeriod:  30,
	})

	tormdr.Start()
	time.Sleep(1 * time.Minute)

}
