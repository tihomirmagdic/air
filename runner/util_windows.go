package runner

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

func (e *Engine) killCmd(cmd *exec.Cmd) (pid int, err error) {
	pid = cmd.Process.Pid
	// https://stackoverflow.com/a/44551450
	// kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
	//return pid, kill.Run()
	return pid, SendInterrupt(pid)
}

func SendInterrupt(pid int) error {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return fmt.Errorf("LoadDLL: %v", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return fmt.Errorf("FindProc: %v", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return fmt.Errorf("GenerateConsoleCtrlEvent: %v", e)
	}
	return nil
}

func (e *Engine) startCmd(cmd string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	var err error

	if !strings.Contains(cmd, ".exe") {
		e.runnerLog("CMD will not recognize non .exe file for execution, path: %s", cmd)
	}
	c := exec.Command("cmd", "/c", cmd)
	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, nil, nil, err
	}
	return c, stdout, stderr, err
}
