package funcs

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

type ExecCmd struct {
	cancel chan bool
	wd     string
	stdout io.Writer
}

func NewExecCmd(wd string, stdout io.Writer) *ExecCmd {
	return &ExecCmd{
		wd:     wd,
		cancel: make(chan bool),
		stdout: stdout,
	}
}

func (e *ExecCmd) Cancel() {
	e.cancel <- true
}

func (e *ExecCmd) Exec(args ...string) error {
	args = append([]string{"-c"}, args...)
	execCmd := exec.Command("/bin/sh", args...)
	execCmd.Env = os.Environ()
	execCmd.Dir = GetPwd()
	if len(e.wd) != 0 {
		execCmd.Dir = e.wd
	}
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdoutScan := bufio.NewScanner(stdout)
	go func() {
		for stdoutScan.Scan() {
			e.stdout.Write([]byte(stdoutScan.Text() + "\n"))
		}
	}()
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return err
	}
	stderrScan := bufio.NewScanner(stderr)
	go func() {
		for stderrScan.Scan() {
			e.stdout.Write([]byte(stderrScan.Text() + "\n"))
		}
	}()

	done := make(chan error)
	if err := execCmd.Start(); err != nil {
		return err
	}
	go func() {
		done <- execCmd.Wait()
	}()

	select {
	case <-e.cancel:
		logrus.Debugf("received cancel signal")
		logrus.Infof("kill process(%v) %v", execCmd.Process, args)
		if err := execCmd.Process.Kill(); err != nil {
			logrus.Infof("Kill command %v failed, error: %v\n", args, err)
		} else {
			logrus.Infof("process %v is killed", execCmd.Process)
		}
		return errors.New(fmt.Sprintf("%v is canceled", args))
	case err := <-done:
		return err
	}
}
