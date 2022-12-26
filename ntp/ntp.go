package ntp

import (
	"time"

	"github.com/beevik/ntp"
)

var timestampPadding int64
var defaultNTPServer = "ntp1.aliyun.com"

func SetupServer(addr string) {
	defaultNTPServer = addr
	return
}
func SyncTime() error {
	t1 := time.Now().UnixMilli()
	t2, err := ntp.Time(defaultNTPServer)
	if err != nil {
		return err
	}
	timestampPadding = t2.UnixMilli() - t1
	return nil
}
func FixedTimestamp() uint64 {
	return uint64(time.Now().UnixMilli() + timestampPadding)
}
func FixedTime() time.Time {
	return time.Now().Add(time.Millisecond * time.Duration(timestampPadding))
}
