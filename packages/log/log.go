package log

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

type Log struct {
	Logtype string
	Loginfo string
}

var LOGch chan Log

var filename string

func init() {
	LOGch = make(chan Log, 50)
	if !checkFileIsExist("./logs/") {
		os.Mkdir("logs", os.ModePerm)
	}
	filename = "./logs/" + time.Now().Format("2006-01-02") + ".log"
	LOGI("CRA API LOGFile: ", filename)
	var err1 error
	var f *os.File
	f, err1 = os.Create(filename)
	if err1 != nil {
		panic(err1)
	}
	f.Close()
}
func LOG() {
	for {
		newlog := <-LOGch
		var f *os.File
		var err1 error
		f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
		if err1 != nil {
			panic(err1)
		}
		var loglevel string
		switch newlog.Logtype {
		case "LOGD":
			loglevel = color.MagentaString(newlog.Logtype)
		case "LOGI":
			loglevel = color.BlueString(newlog.Logtype)
		case "LOGW":
			loglevel = color.YellowString(newlog.Logtype)
		case "LOGE":
			loglevel = color.RedString(newlog.Logtype)
		}
		s := fmt.Sprint(
			"CRA API ",
			loglevel,
			" [",
			color.CyanString(time.Now().Format("2006-01-02 15:04:05")),
			"] ",
			newlog.Loginfo,
			"\n",
		)
		fmt.Print(s)
		f.WriteString(s)
		err1 = f.Close()
		if err1 != nil {
			panic(err1)
		}
	}
}
func LOGI(a ...any) {
	newlog := Log{
		Logtype: "LOGI",
		Loginfo: fmt.Sprint(a...),
	}
	LOGch <- newlog
}
func LOGD(a ...any) {
	newlog := Log{
		Logtype: "LOGD",
		Loginfo: fmt.Sprint(a...),
	}
	LOGch <- newlog
}
func LOGW(a ...any) {
	newlog := Log{
		Logtype: "LOGW",
		Loginfo: fmt.Sprint(a...),
	}
	LOGch <- newlog
}
func LOGE(a ...any) {
	newlog := Log{
		Logtype: "LOGE",
		Loginfo: fmt.Sprint(a...),
	}
	LOGch <- newlog
}
