package ezCaptcha

import (
	"fmt"
	"github.com/Anveena/ezTools/ezFile"
	"github.com/Anveena/ezTools/ezRandom"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	SetGrpcPathAndCert("127.0.0.1:12580", "JiangNan7Guai", "../../../JiangNan7Guai/commonTools/x509/crt")
	for i := 0; i < 10; i++ {
		data, ans, err := Get()
		if err != nil {
			panic(err)
		}
		f, err := ezFile.CreateFile("./test", fmt.Sprintf("%s------ans%s.png", ezRandom.RandomString(ezRandom.NumberAndAllLetter, 8), ans), true, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
		if err != nil {
			panic(err)
		}
		if _, err = f.Write(data); err != nil {
			panic(err)
		}
		if err = f.Close(); err != nil {
			panic(err)
		}
	}
}
