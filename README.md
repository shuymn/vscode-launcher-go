# Visual Studio Code Launcher

This is Visual Studio Code launcher written in Go.

**:warning: This application is intended for personal use and only supports running on macOS for now. :warning:**

# Installation

```bash
go install github.com/shuymn/vscode-launcher-go@latest
```

Then, it is recommend to set up your `.bashrc` or `.zshrc` as follows.

```bash
alias code="vscode-launcher-go"
```

# Motivation

The original Visual Studio Code launcher for macOS is a Shell Script using Python, so if you have security software installed on your PC that interferes with Python execution, you may not be able to use the original launcher. Therefore, I have reimplemented the launcher using Go.

It is possible that non-macOS environments have similar problems, but since I don't use Visual Studio Code on non-macOS, I don't know how it really works or how many people have trouble with it. If you are facing a similar problems on non-macOS I would be very happy if you could send me a PR for this application.