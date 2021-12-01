package ezUTCTime

import (
	"time"

	"github.com/beevik/ntp"
)

var timestamPadding int64

func SyncTimeFromAliyun() error {
	t1 := time.Now().UnixMilli()
	t2, err := ntp.Time("ntp2.aliyun.com")
	if err != nil {
		return err
	}
	timestamPadding = t2.UnixMilli() - t1
	return nil
}
func GetAliyunTimestamp() uint64 {
	return uint64(time.Now().UnixMilli() + timestamPadding)
}
func GetAliyunTime() time.Time {
	return time.Now().Add(time.Millisecond * time.Duration(timestamPadding))
}
