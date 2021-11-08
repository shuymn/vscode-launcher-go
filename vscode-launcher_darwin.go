package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yookoala/realpath"
)

const (
	VSCODE_EXECUTABLE_NAME          = "code"
	VSCODE_INSIDERS_EXECUTABLE_NAME = "code-insiders"
)

const (
	ExitOK             = 0
	ExitFlagParseError = 1
	ExitError          = 1
)

func main() {
	var insiders bool

	flags := flag.NewFlagSet("vscode-launcher-go", flag.ContinueOnError)
	flags.BoolVar(&insiders, "insiders", false, "launch vscode insiders")

	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(ExitFlagParseError)
	}

	var file string
	if insiders {
		file = VSCODE_INSIDERS_EXECUTABLE_NAME
	} else {
		file = VSCODE_EXECUTABLE_NAME
	}

	out, err := run(file, flags.Args())
	if err != nil {
		fmt.Print(err)
		os.Exit(ExitError)
	}

	fmt.Print(string(out))
	os.Exit(ExitOK)
}

func run(file string, args []string) ([]byte, error) {
	code, err := exec.LookPath(file)
	if err != nil {
		return nil, err
	}

	path, err := realpath.Realpath(code)
	if err != nil {
		return nil, err
	}

	contents := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(path))))
	electron := fmt.Sprintf("%s/MacOS/Electron", contents)
	cli := fmt.Sprintf("%s/Resources/app/out/cli.js", contents)

	args = append([]string{cli}, args...)
	cmd := exec.Command(electron, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "ELECTRON_RUN_AS_NODE=1")

	return cmd.Output()
}
