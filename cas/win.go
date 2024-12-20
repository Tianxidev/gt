package cas

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

type Win struct {
}

func New() *Win {
	return &Win{}
}

// Execute 执行
func (w *Win) Execute(lpOperation, lpFile, lpParameters string) error {
	var param *uint16

	if lpParameters != "" {
		param = syscall.StringToUTF16Ptr(lpParameters)
	}
	err := windows.ShellExecute(
		0,
		syscall.StringToUTF16Ptr(lpOperation),
		syscall.StringToUTF16Ptr(lpFile),
		param,
		nil,
		int32(win.SW_HIDE))

	return err
}

// ExecuteWithOutput 带输出的执行
func (w *Win) ExecuteWithOutput(lpFile string, lpParameters string) (string, error) {
	cmdExec := fmt.Sprintf(`"%s" %s`, lpFile, lpParameters)
	cmd := exec.Command("cmd.exe")
	fmt.Println(cmdExec)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, cmdExec), HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetShortPathName 获取短路径名
func (w *Win) GetShortPathName(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("the path '%s' does not exist", path)
	}

	utf16Path, err := windows.UTF16FromString(path)
	if err != nil {
		return "", err
	}

	var shortPath [windows.MAX_PATH]uint16
	ret, err := windows.GetShortPathName(&utf16Path[0], &shortPath[0], windows.MAX_PATH)
	if err != nil {
		return "", err
	}

	if ret == 0 {
		return "", fmt.Errorf("failed to get short path name")
	}

	return windows.UTF16ToString(shortPath[:]), nil
}
