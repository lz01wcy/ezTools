package ezLog

import (
	"fmt"
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	LogLvDebug = int32(iota + 1)
	LogLvInfo
	LogLvError
	LogLvDingMessage
	LogLvDingLists
	LogLvDingAll
)

var lvHeaderMap = map[int32]string{
	LogLvDebug:       "\n[Debug] ",
	LogLvInfo:        "\n[Info]  ",
	LogLvError:       "\n[Error] ",
	LogLvDingMessage: "\n[Ding]  ",
	LogLvDingLists:   "\n[DAL!]  ",
	LogLvDingAll:     "\n[DAA!]  ",
}
var logFmtStr = "%s File:%s%s Line:%d%s %s\n"

type EZLoggerModel struct {
	LogLevel      int32
	AppName       string
	DingTalkModel struct {
		Enable                 bool
		SecretKeyEncodedString string
		URLEncodedString       string
		Mobiles                []string
	}
	GRPCModel struct {
		ClientCounts     int
		URLEncodedString string
	}
}

var (
	enableDing        = false
	logLevel          = LogLvInfo
	appName           = ""
	dingAtStr         = ""
	dingURL           = ""
	dingSecretKey     = ""
	dingSecretKeyData []byte
	dingMobiles       []string
	gRPCURL           string
	gRPCClientCounts  int
	logChannel        = make(chan *ezLogPB.EZLogReq, 65536)
)

type dingRequestModel struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}

func SetUpEnv(m *EZLoggerModel) (err error) {
	appName = m.AppName
	logLevel = m.LogLevel
	enableDing = m.DingTalkModel.Enable
	if appName == "" {
		return fmt.Errorf("AppName必须要配置,不然无法区分对应的服务")
	}
	if enableDing {
		if m.DingTalkModel.URLEncodedString == "" {
			return fmt.Errorf("钉钉的那个URL没配置")
		}
		dingURL, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.DingTalkModel.URLEncodedString)
		if err != nil {
			return
		}
		if m.DingTalkModel.SecretKeyEncodedString == "" {
			return fmt.Errorf("钉钉的那个SecretKey没配置")
		}
		dingSecretKey, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.DingTalkModel.SecretKeyEncodedString)
		if err != nil {
			return
		}
		dingSecretKeyData = []byte(dingSecretKey)
		dingMobiles = m.DingTalkModel.Mobiles
		sb := strings.Builder{}
		sb.WriteString("### **责任人**:")
		for _, people := range dingMobiles {
			sb.WriteString(fmt.Sprintf(" @%s", people))
		}
		sb.WriteString("\n")
		dingAtStr = sb.String()
	}
	gRPCClientCounts = m.GRPCModel.ClientCounts
	if gRPCClientCounts > 0 {
		if m.GRPCModel.URLEncodedString == "" {
			return fmt.Errorf("gRPC日志服务没配地址")
		}
		gRPCURL, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.GRPCModel.URLEncodedString)
		if err != nil {
			return
		}
		for i := 0; i < gRPCClientCounts; i++ {
			go startGRPCClient()
		}
	}
	return nil
}
func D(msg ...interface{}) {
	ezlog(LogLvDebug, msg...)
}
func I(msg ...interface{}) {
	ezlog(LogLvInfo, msg...)
}
func E(msg ...interface{}) {
	ezlog(LogLvError, msg...)
}
func DingMessage(msg ...interface{}) {
	ezlog(LogLvDingMessage, msg...)
	if enableDing {
		sendToDing(LogLvDingMessage, "no tag", fmt.Sprintln(msg...))
	}
}
func DingAtAll(msg ...interface{}) {
	ezlog(LogLvDingAll, msg...)
	if enableDing {
		sendToDing(LogLvDingAll, "no tag", fmt.Sprintln(msg...))
	}
}
func DingList(msg ...interface{}) {
	ezlog(LogLvDingLists, msg...)
	if enableDing {
		sendToDing(LogLvDingLists, "no tag", fmt.Sprintln(msg...))
	}
}

func DWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDebug, tag, msg...)
}
func IWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvInfo, tag, msg...)
}
func EWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvError, tag, msg...)
}
func DingMessageWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingMessage, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingMessage, tag, fmt.Sprintln(msg...))
	}
}
func DingAtAllWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingAll, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingAll, tag, fmt.Sprintln(msg...))
	}
}
func DingListWithTag(tag string, msg ...interface{}) {
	ezlogWithTag(LogLvDingLists, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingLists, tag, fmt.Sprintln(msg...))
	}
}
