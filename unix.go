package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func TerminalSizes() (width, height int, err error) {
	cmd := exec.Command("stty", "size") //todo context
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected stty output: %s", out)
	}

	height, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	width, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func SetTerminalRowMod() (*unix.Termios, func(), error) {
	old, err := unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get old terminal state: %w", err)
	}

	term := *old
	term.Lflag &^= unix.ICANON | unix.ECHO

	if err = unix.IoctlSetTermios(unix.Stdin, unix.TCSETS, &term); err != nil {
		return nil, nil, fmt.Errorf("cannot set new terminal state: %w", err)
	}

	rollback := func() {
		if err = unix.IoctlSetTermios(unix.Stdin, unix.TCSETS, old); err != nil {
			fmt.Printf("cannot rollback terminal state, please close this terminal session")
		}
	}

	return &term, rollback, nil
}
