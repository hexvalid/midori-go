package logger

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

const timeFormat = "15:04:05.000"

//Examples:
//	log := NewLog("ana", color.FgCyan)
//	log.Info("Bu bir bilgi mesajı")
//	log.SInfo(color.YellowString("çok özel"), "Bu bir sayı: %d", 1212)
//	log.SInfo("normal özel", "Bu kırmızı bir string olacak: %s", color.GreenString("YEŞİL"))
//	log.Error("Bir hata oluştu")
//	log.SError("normal özel", "Özelin bir hatası")
type Logger struct {
	Name  string
	Color *color.Color
}

func NewLog(name string, attr color.Attribute) *Logger {
	return &Logger{
		Name:  name,
		Color: color.New(attr),
	}
}

func (logger *Logger) Info(msg string, a ...interface{}) {
	fmt.Fprintf(color.Output, color.HiBlackString(time.Now().Format(timeFormat))+
		logger.Color.Sprint(" "+logger.Name+": ")+msg+"\n", a...)
}

func (logger *Logger) Error(msg string, a ...interface{}) {
	fmt.Fprintf(color.Output, color.HiBlackString(time.Now().Format(timeFormat))+
		logger.Color.Sprint(" "+logger.Name+"(")+color.RedString("!")+
		logger.Color.Sprint("): ")+msg+"\n", a...)
}

func (logger *Logger) SInfo(smsg, msg string, a ...interface{}) {
	fmt.Fprintf(color.Output, color.HiBlackString(time.Now().Format(timeFormat))+
		logger.Color.Sprint(" "+logger.Name+"(")+smsg+logger.Color.Sprint("): ")+msg+"\n", a...)
}
