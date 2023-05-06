package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/yookoala/realpath"
	"golang.org/x/exp/slices"
)

const cmdName = "vscode-launcher-go"

const (
	vscodeExecutableName         = "code"
	vscodeInsidersExecutableName = "code-insiders"
)

const exitCodeError = 1

var ErrFlagParseFailed = errors.New("flag parse failed")

func main() {
	log.SetFlags(0)
	err := Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		if !errors.Is(err, ErrFlagParseFailed) {
			log.Println(err)
		}
		exitCode := exitCodeError
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}

type bypassPythonFlag struct {
	Value bool
}

func (f *bypassPythonFlag) Set(s string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("This flag can only be set on macOS")
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	f.Value = v
	return nil
}

func (f *bypassPythonFlag) String() string {
	return strconv.FormatBool(f.Value)
}

func (f *bypassPythonFlag) IsBoolFlag() bool {
	return true
}

func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	flags := flag.NewFlagSet(fmt.Sprintf("%s (v%s rev:%s)", cmdName, strings.TrimPrefix(version, "v"), revision), flag.ContinueOnError)
	flags.SetOutput(errStream)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), `
Usage of %s:
$ vscode-launcher-go [options] [-- vscode-arguments...]

Description:
  Launch Visual Studio Code.

Example:
  $ vscode-launcher-go -insiders -bypass-python -- --help

Options:
`, flags.Name())
		flags.PrintDefaults()
	}

	var (
		insiders     bool
		bypassPython bypassPythonFlag
		osProfile    bool
	)

	flags.BoolVar(&insiders, "insiders", false, "launch vscode insiders")
	flags.Var(&bypassPython, "bypass-python", "bypass python launcher (macOS only)")
	flags.BoolVar(&osProfile, "os-profile", false, "detect profile according to the OS")

	if err := flags.Parse(argv); err != nil {
		return errors.Join(err, ErrFlagParseFailed)
	}
	return run(ctx, insiders, bypassPython.Value, osProfile, flags.Args(), outStream, errStream)
}

func run(ctx context.Context, insiders, bypassPython, osProfile bool, argv []string, outStream, errStream io.Writer) error {
	if i := slices.Index(argv, "--"); i != -1 {
		argv = argv[i+1:]
	}

	if osProfile && !slices.Contains(argv, "--profile") {
		argv = append(argv, "--profile", strings.Title(runtime.GOOS))
	}

	code, err := getExecutablePath(insiders)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" && bypassPython {
		contents := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(code))))
		electron := fmt.Sprintf("%s/MacOS/Electron", contents)
		cli := fmt.Sprintf("%s/Resources/app/out/cli.js", contents)

		argv = append([]string{"--ms-enable-electron-run-as-node", cli}, argv...)
		cmd := exec.Command(electron, argv...)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "ELECTRON_RUN_AS_NODE=1")
	} else {
		cmd = exec.Command(code, argv...)
		cmd.Env = os.Environ()
	}
	cmd.Stdout = outStream
	cmd.Stderr = errStream
	return cmd.Run()
}

func getExecutablePath(insiders bool) (string, error) {
	var executable string
	if insiders {
		executable = vscodeInsidersExecutableName
	} else {
		executable = vscodeExecutableName
	}

	code, err := exec.LookPath(executable)
	if err != nil {
		return "", err
	}

	path, err := realpath.Realpath(code)
	if err != nil {
		return "", err
	}

	return path, nil
}
