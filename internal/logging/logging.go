package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/awesome-gocui/gocui"
)

var (
	InfoChan = make(chan string)
	WarnChan = make(chan string)
	ErrChan  = make(chan error)
	StopChan = make(chan struct{})

	LogDir  = "./log"
	logFile *os.File
)

func init() {
	// create log file
	err := initLogDir()
	if err != nil {
		panic(err)
	}
	logFileName := fmt.Sprintf("%s/%s.log", LogDir, GetTimestamp())
	logFile, err = os.Create(logFileName)
	if err != nil {
		panic(err)
	}
}

// check if the log directory exists and if not create it
func initLogDir() error {
	_, err := os.Stat(LogDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(LogDir, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetTimestamp() string {
	return fmt.Sprintf("%d", int(time.Now().Unix()))
}

func FeedLog(g *gocui.Gui, msgLog <-chan string, warnLog <-chan string, errLog <-chan error, stopChan <-chan struct{}) {

	defer logFile.Close()

	for {
		select {
		case log := <-msgLog:
			log = "[INFO] " + log
			writeLogToView(g, log)
			writeLogToFile(log)

		case log := <-warnLog:
			log = "[WARN] " + log
			writeLogToView(g, log)
			writeLogToFile(log)

		case logE := <-errLog:
			log := "[ERRO] " + logE.Error()
			writeLogToView(g, log)
			writeLogToFile(log)

		case <-stopChan:
			return
		}

	}
}

func writeLogToView(g *gocui.Gui, log string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("log")
		if err != nil {
			return err
		}
		fmt.Fprintln(v, log)
		return nil
	})
}

func writeLogToFile(log string) {
	time := time.Now().Format("2006-01-02 15:04:05")
	logFile.WriteString(fmt.Sprintf("%s %s\n", time, log))
}
