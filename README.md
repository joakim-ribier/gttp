# GTTP - Go HTTP Client for Terminal UIs

## TOC

* [Dependencies](#dependencies)
* [Installation](#installation)
* [Troubleshooting](#troubleshooting)

## Dependencies

* tview

`GTTP` project is based on the awesome https://github.com/rivo/tview library to make terminal user interfaces.

Thanks to [Rivo](https://github.com/rivo)

* clipboard

It also used https://github.com/atotto/clipboard library "Provide copying and pasting to the Clipboard for Go".

Thanks to [Ato Araki](https://github.com/atotto)

## Installation

```bash
# Download the Go project
$ go get github.com/joakim-ribier/gttp

# Build
$ cd $HOME/go/src/gttp
$ go build

# Execute
$ ./gttp data.json
```

## Troubleshooting

### Linux, Unix

:warning:
If the shortcuts (copy/paste) do not work, read the README.md of [clipboard](https://github.com/atotto/clipboard) for more specific information.

> Linux, Unix (requires 'xclip' or 'xsel' command to be installed)

```bash
$ sudo apt install xclip
```

![](https://media0.giphy.com/media/1XgIXQEzBu6ZWappVu/giphy.gif)
