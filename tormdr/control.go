package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"net"
	"net/textproto"
)

const (
	controlMsgAuth     = "AUTHENTICATE"
	controlMsgQuit     = "QUIT"
	controlMsgNewNym   = "SIGNAL NEWNYM"
	controlMsgShutdown = "SIGNAL SHUTDOWN"
)

func writeCtrlMsg(s string, ctrlConn net.Conn, reader *textproto.Reader) (err error) {
	if _, err = ctrlConn.Write([]byte(s + "\n")); err != nil {
		return err
	}
	code, message, err := reader.ReadResponse(250)
	if err != nil {
		err = fmt.Errorf("error at control message: %s", err)
	} else if code != 250 {
		err = fmt.Errorf("unexcepted control message: %s, code:%d", message, code)
	}
	return
}

func (tormdr *TorMDR) sendCtrlMsg(s string) (err error) {
	ctrlConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", localhost, tormdr.controlPort))
	if err != nil {
		return err
	}
	defer func() {
		if err := ctrlConn.Close(); err != nil {
			log.SInfo(fmt.Sprintf("%03d", tormdr.no), "%s %s", color.RedString("(Error)"),
				"Control socket can't closed")
		}
	}()
	reader := textproto.NewReader(bufio.NewReader(ctrlConn))
	if err = writeCtrlMsg(controlMsgAuth, ctrlConn, reader); err != nil {
		return err
	}
	if err = writeCtrlMsg(s, ctrlConn, reader); err != nil {
		return err
	}
	if err = writeCtrlMsg(controlMsgQuit, ctrlConn, reader); err != nil {
		return err
	}
	return nil
}
