package curses

// #cgo LDFLAGS: -lcurses
// #include <curses.h>
import "C"
import "fmt"

func Begin() {
    C.initscr()
    C.noecho()
    C.cbreak()
    C.curs_set(C.int(0))
}

func Printf(x,y int, format string, args ...interface{}) {
    str := C.CString(fmt.Sprintf(format, args...))
    C.mvaddstr(C.int(y), C.int(x), str)
}

func Clear() {
    C.clear()
}

func Refresh() {
    C.refresh()
}

func End() {
    C.endwin()
}

func MaxXY() (x, y int) {
    x = int(C.getmaxx(C.stdscr))
    y = int(C.getmaxy(C.stdscr))
    return
}
