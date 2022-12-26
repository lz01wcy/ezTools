//go:build ezDebug

package log

import (
	"fmt"
	"runtime"
	"strings"
)

func startGRPCClient() {}
func Log(level LogLv, msg ...any) {
	_, file, line, _ := runtime.Caller(2)
	header := lvHeaderMap[level]
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func LogWithTag(level LogLv, tag string, msg ...any) {
	_, file, line, _ := runtime.Caller(2)
	header := fmt.Sprintf("%s[%s]", lvHeaderMap[level], tag)
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func sendToDing(_ LogLv, _ string, _ string) {}
