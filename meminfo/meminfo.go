package meminfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	. "github.com/t-hg/syncp/must"
)

type MemInfo struct {
	Dirty     int
	Writeback int
}

func ReadMemInfo() MemInfo {
    meminfo := MemInfo{}
	file := Must(os.Open("/proc/meminfo"))
    defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		name := strings.TrimSuffix(fields[0], ":")
        value := fields[1]
        dirtySeen := false
        writebackSeen := false
        switch name {
        case "Dirty":
            dirtySeen = true
            meminfo.Dirty = Must(strconv.Atoi(value))
        case "Writeback":
            writebackSeen = true
            meminfo.Writeback = Must(strconv.Atoi(value))
        }
        if dirtySeen && writebackSeen {
            break
        }
	}
    return meminfo
}
