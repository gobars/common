package funcs

import (
	"github.com/sirupsen/logrus"
	"bufio"
	"os"
	"path"
	"sync"
	"time"
)

var (
	CancelCommandTimeout = 25 * time.Second
)

type LogConsole struct {
	stop    chan bool
	closed  chan bool
	write   chan []byte
	logPath string
	LogFile *os.File
	writer  *bufio.Writer
	mutex   *sync.RWMutex
}

func NewLogConsole(write chan []byte, logPath string) *LogConsole {
	console := &LogConsole{
		write:   write,
		stop:    make(chan bool),
		closed:  make(chan bool),
		logPath: logPath,
		mutex:   new(sync.RWMutex),
	}
	go console.writeLog()
	return console
}

func (console *LogConsole) writeLog() {
	filePath := path.Dir(console.logPath)
	logrus.Infof("logpath:%s", filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0755)
	}
	file, err := os.OpenFile(console.logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if isError(err) {
		return
	}
	writer := bufio.NewWriter(file)
	console.LogFile = file
	console.writer = writer
	for {
		select {
		case msg := <-console.write:
			str := string(msg)
			logrus.Infof("--->%s", str)
			console.writeByLock(msg)
		case <-console.stop:
			logrus.Info("log console stop ", console.logPath)
			console.LogFile.Close()
			return
		}
	}
}

func (console *LogConsole) writeByLock(msg []byte) {
	console.mutex.Lock()
	defer console.mutex.Unlock()
	console.writer.Write(msg)
	console.writer.Flush()
}

func (console *LogConsole) Write(data []byte) (int, error) {
	console.write <- data
	return len(data), nil
}

func (console *LogConsole) Close() error {
	defer console.LogFile.Close()
	return CloseAndWait(console.stop, console.closed, CancelCommandTimeout)
}
