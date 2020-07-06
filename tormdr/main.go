package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"midori-go/logger"
	"net"
	"net/http"
	"os"
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

var (
	log       = logger.NewLog("TorMDR", color.FgMagenta)
	regexBoot = regexp.MustCompile(`\((.*?)\)`)
)

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
	ctrlConn   net.Conn
	transport  *http.Transport
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
	_ = os.MkdirAll(cfg.DataDirectory, os.ModePerm)
	return &tormdr
}

func (tormdr *TorMDR) Start() (err error) {
	if tormdr.stdoutPipe, err = tormdr.cmd.StdoutPipe(); err != nil {
		return err
	}

	if err = tormdr.cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(tormdr.stdoutPipe)
	if scanner.Scan() {
		versionLine := scanner.Text()
		if strings.Contains(versionLine, "TorMDR=") {
			versionParts := strings.Split(versionLine, " ")
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "%s %s started with %s %s, %s %s and %s %s.",
				color.MagentaString(strings.Split(versionParts[1], "=")[0]),
				color.HiMagentaString(strings.Split(versionParts[1], "=")[1]),
				color.CyanString(strings.Split(versionParts[2], "=")[0]),
				color.HiCyanString(strings.Split(versionParts[2], "=")[1]),
				color.CyanString(strings.Split(versionParts[3], "=")[0]),
				color.HiCyanString(strings.Split(versionParts[3], "=")[1]),
				color.CyanString(strings.Split(versionParts[4], "=")[0]),
				color.HiCyanString(strings.Split(versionParts[4], "=")[1]),
			)
		} else {
			errMsg := "TorMDR is responded unexpectedly"
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), color.RedString("(Error) %s"), errMsg)
			return errors.New(errMsg)
		}
	}

	//todo: make timeout

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Bootstrapped") {
			tormdr.bootStatus = regexBoot.FindStringSubmatch(line)[1]
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Boot status: %s", color.BlueString(tormdr.bootStatus))
			if tormdr.bootStatus == "done" {
				break
			}
		} else if strings.Contains(line, "[warn]") {
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "%s%s", color.YellowString("(Warning)"),
				strings.Split(line, "[warn]")[1])
		} else if strings.Contains(line, "[err]") {
			errMsg := strings.Split(line, "[err]")[1]
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "%s%s", color.RedString("(Error)"), errMsg)
			return errors.New(errMsg)
		}
	}

	tormdr.ctrlConn, err = net.Dial("tcp", "127.0.0.1:9051")
	//auth ctrl

	return nil
}

func (tormdr *TorMDR) Stop() (err error) {
	//todo: exit over telnet

	if tormdr.cmd.Process != nil {
		_ = tormdr.cmd.Process.Kill()
		_, _ = tormdr.cmd.Process.Wait()
	}
	return nil
}

func main() {
	tormdr := NewTorMDR(1, &TorMDRConfig{
		TorMDRBinaryPath: "/home/anon/tormdr/tormdr",
		DataDirectory:    "/tmp/tormdr_data",
		KeepalivePeriod:  30,
	})
	fmt.Println(tormdr.cmd.Process == nil)
	tormdr.Start()
	time.Sleep(3 * time.Second)
	tormdr.Stop()
}
