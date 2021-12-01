//go:build ezDebug

package ezLog

import (
	"fmt"
	"runtime"
	"strings"
)

func startGRPCClient() {}
func ezlog(level int32, msg ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	header := lvHeaderMap[level]
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func ezlogWithTag(level int32, tag string, msg ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	header := fmt.Sprintf("%s[%s]", lvHeaderMap[level], tag)
	fmt.Printf(logFmtStr, header, file, header, line, header,
		strings.ReplaceAll(fmt.Sprintln(msg...), "\n", header))
}
func sendToDing(_ int32, _ string, _ string) {}
