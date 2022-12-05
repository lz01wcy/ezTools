//go:build ezDebug

package ezLog

import (
	"fmt"
	"runtime"
	"strings"
)

func startGRPCClient() {}
func Log(level LogLv, msg ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	header := lvHeaderMap[level]
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func LogWithTag(level LogLv, tag string, msg ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	header := fmt.Sprintf("%s[%s]", lvHeaderMap[level], tag)
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func sendToDing(_ LogLv, _ string, _ string) {}
