package ezUTCTime

import (
	"time"

	"github.com/beevik/ntp"
)

var timestampPadding int64
var defaultNTPServer = "ntp2.aliyun.com"

func SetupNTPServer(addr string) error {
	defaultNTPServer = addr
	return SyncTimeFromNTPServer()
}
func SyncTimeFromNTPServer() error {
	t1 := time.Now().UnixMilli()
	t2, err := ntp.Time(defaultNTPServer)
	if err != nil {
		return err
	}
	timestampPadding = t2.UnixMilli() - t1
	return nil
}
func GetFixedTimestamp() uint64 {
	return uint64(time.Now().UnixMilli() + timestampPadding)
}
func GetFixedTime() time.Time {
	return time.Now().Add(time.Millisecond * time.Duration(timestampPadding))
}
