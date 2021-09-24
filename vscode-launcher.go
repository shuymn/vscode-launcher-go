package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yookoala/realpath"
)

const VSCODE_EXECUTABLE_NAME = "code"

const (
	ExitOK = iota
	ExitErr
)

func main() {
	code, err := exec.LookPath(VSCODE_EXECUTABLE_NAME)
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitErr)
	}
	path, err := realpath.Realpath(code)
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitErr)
	}

	contents := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(path))))
	electron := fmt.Sprintf("%s/MacOS/Electron", contents)
	cli := fmt.Sprintf("%s/Resources/app/out/cli.js", contents)

	args := []string{cli}
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	cmd := exec.Command(electron, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "ELECTRON_RUN_AS_NODE=1")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitErr)
	}
	os.Exit(ExitOK)
}
