package funcs

import (
	"os"
	"strconv"
	"time"
	"path/filepath"
	"fmt"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/gobars/commons/lang/dates"
)

func GetWorkPath() string {
	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, "tmp", GetTag())
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0755)
	}
	return filePath
}

func GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

func GetTag() string {
	return dates.ToDateStr() + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

func CloseAndWait(stop, closed chan bool, timeout time.Duration) error {
	select {
	case _, ok := <-stop:
		if !ok {
			return nil
		}
	default:
	}

	close(stop)

	select {
	case <-closed:
		return nil
	case <-time.After(timeout):
		return Err("Wait for closed timeout")
	}
}

func Sprintf(f string, args ...interface{}) string {
	return fmt.Sprintf(f, args...)
}

func Err(f string, args ...interface{}) error {
	return errors.New(Sprintf(f, args...))
}

func isError(err error) bool {
	if err != nil {
		logrus.Error(err.Error())
	}
	return err != nil
}

func GetFilenameExtension(path string) (ext string) {
	var extIndex int
	if path == "" {
		return
	}
	if extIndex = strings.LastIndex(path, "."); extIndex == -1 {
		return
	}
	folderIndex := strings.LastIndex(path, "/")
	if folderIndex > extIndex {
		return
	}
	runes := []rune(path)
	ext = string(runes[extIndex+1:])
	return
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
