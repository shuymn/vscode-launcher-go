# Visual Studio Code Launcher

This is Visual Studio Code launcher written in Go.

**:warning: This application is intended for personal use. :warning:**

# Installation

```bash
go install github.com/shuymn/vscode-launcher-go@latest
```

# Usage

```
vscode-launcher-go [options] [-- vscode-arguments...]
```

## Options

```
-insiders bool
    Launch VSCode Insiders
-bypass-python bool
    Bypass python based launcher (macOS only)
-os-profile bool
    Detect the profile according to the OS
```
