package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	log       = logger.NewLog("TorMDR", color.FgMagenta)
	regexBoot = regexp.MustCompile(`\((.*?)\)`)
)

const (
	goLibVersion     string = "1.0S"
	localhost        string = "127.0.0.1"
	socksPortStart   int    = 20000
	controlPortStart int    = 40000
)

type TorMDR struct {
	cmd            *exec.Cmd
	no             int
	proxy          *url.URL
	bootStatus     string
	stdoutPipe     io.ReadCloser
	socksPort      int
	controlPort    int
	targetExitNode string
	cacheDir       string
}

func NewTorMDR(no int, cfg *TorMDRConfig) (tormdr *TorMDR, err error) {
	tormdr = &TorMDR{no: no, cmd: &exec.Cmd{}}
	tormdr.cmd.Path = cfg.TorMDRBinaryPath
	tormdr.cmd.Args = append(tormdr.cmd.Args, defaultArgs...)
	tormdr.socksPort = socksPortStart + no
	tormdr.controlPort = controlPortStart + no
	dataDir := path.Join(cfg.DataDirectory, strconv.Itoa(no))
	tormdr.cacheDir = path.Join(dataDir, "cache")
	tormdr.cmd.Args = append(tormdr.cmd.Args, "SocksPort", strconv.Itoa(tormdr.socksPort))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "ControlPort", strconv.Itoa(tormdr.controlPort))
	tormdr.cmd.Args = append(tormdr.cmd.Args, "DataDirectory", dataDir)
	tormdr.cmd.Args = append(tormdr.cmd.Args, "CacheDirectory", tormdr.cacheDir)
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
	if err = os.MkdirAll(cfg.DataDirectory, os.ModePerm); err != nil {
		return nil, err
	}
	if tormdr.proxy, err = url.Parse(fmt.Sprintf("socks5://%s:%d",
		localhost, tormdr.socksPort)); err != nil {
		return nil, err
	}
	return
}

func (tormdr *TorMDR) Start() (err error) {
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "%s-%s %s starting...",
		color.MagentaString("TorMDR"), color.CyanString("GoLib"), color.HiCyanString(goLibVersion))
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

	if err = populateExitNodeCache(tormdr.cacheDir); err != nil {
		return err
	}
	return nil
}

func (tormdr *TorMDR) Stop() (err error) {
	if tormdr.cmd.Process != nil {
		log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Sending shutdown signal...")
		_ = tormdr.sendCtrlMsg(controlMsgShutdown)
		time.Sleep(250 * time.Millisecond)
		log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Stopping process...")
		_ = tormdr.cmd.Process.Kill()
		_, _ = tormdr.cmd.Process.Wait()
	}
	return nil
}

func (tormdr *TorMDR) NewCircuit() error {
	if len(tormdr.targetExitNode) > 4 {
		log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Resetting and building new circuit...")
		if err := tormdr.sendCtrlMsg("SETCONF ExitNodes="); err != nil {
			return err
		} else {
			tormdr.targetExitNode = ""
		}
	} else {
		log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Building new circuit...")
	}
	return tormdr.sendCtrlMsg(controlMsgNewNym)
}

func (tormdr *TorMDR) SetExitNode(ip string) error {
	tormdr.targetExitNode = ip
	log.SInfo(fmt.Sprintf("%03d", tormdr.no), "Setting Exit Node: %s...", color.YellowString(ip))
	if tormdr.bootStatus == "done" {
		return tormdr.sendCtrlMsg(fmt.Sprintf("SETCONF ExitNodes=%s", ip))
	} else {
		tormdr.cmd.Args = append(tormdr.cmd.Args, "ExitNodes", ip)
		return nil
	}
}

func main() {
	tormdr, _ := NewTorMDR(1, &TorMDRConfig{
		TorMDRBinaryPath:    "/usr/bin/tormdr",
		DataDirectory:       "/tmp/tormdr_data",
		KeepalivePeriod:     60,
		UseSocks5Proxy:      true,
		Socks5ProxyAddress:  "40.70.243.118:32416",
		Socks5ProxyUserName: "e4cf6e290c0cf8ae8fb91fcf818e1e40",
		Socks5ProxyPassword: "a565ab1f3802afbf4d07c1674069d813",
	})

	if err := tormdr.Start(); err != nil {
		panic(err)
	}

	en, err := FindExitNode(nil, 20000, true, true, true)

	_ = tormdr.SetExitNode(en)

	if err != nil {
		panic(err)
	}

	fmt.Println(tormdr.TestIP())

	err = tormdr.Stop()
	if err != nil {
		panic(err)
	}

}
