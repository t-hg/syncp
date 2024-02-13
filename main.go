package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/t-hg/syncp/curses"
	. "github.com/t-hg/syncp/meminfo"
)

var running = true
var resized = false

func main() {
	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt)
		select {
		case <-signalChan:
			running = false
		}
	}()

	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGWINCH)
		for {
			select {
			case <-signalChan:
				resized = true
			}
		}
	}()

	defer func() {
		curses.End()
		if val := recover(); val != nil {
			fmt.Fprintln(os.Stderr, val)
		}
	}()

	go func() {
		cmd := exec.Command("sync")
		err := cmd.Run()
		//running = false
		if err != nil {
			panic(err)
		}
	}()

	dirtyMax := 0

	curses.Begin()
	for running {
		meminfo := ReadMemInfo()
		if meminfo.Dirty > dirtyMax {
			dirtyMax = meminfo.Dirty
		}
		percentage := -(1/float64(dirtyMax))*float64(meminfo.Dirty) + 1
		x, _ := curses.MaxXY()
		progressFormat := fmt.Sprintf("[%%-%ds]", x-6)
		progressFill := strings.Repeat("#", int((float64(x)-6)*percentage))
		progress := fmt.Sprintf(progressFormat, progressFill)
		curses.Printf(0, 4, progressFormat)
		curses.Clear()
		curses.Printf(0, 0, "Sync in progress...")
		curses.Printf(0, 1, "Dirty: %d, Writeback: %d\n", meminfo.Dirty, meminfo.Writeback)
		curses.Printf(0, 2, progress+"%3d%%", int(percentage*100))
		if resized {
			resized = false
			curses.End()
		}
		curses.Refresh()
		time.Sleep(1 * time.Second)
	}
}
