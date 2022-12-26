package ezLog

import (
	"fmt"
	"github.com/Anveena/ezTools/ezLog/ezLogPB"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type LogLv int32

const (
	LogLvDebug = LogLv(iota + 1)
	LogLvInfo
	LogLvError
	LogLvDingMessage
	LogLvDingLists
	LogLvDingAll
)

var lvHeaderMap = map[LogLv]string{
	LogLvDebug:       "\n[Debug]  ",
	LogLvInfo:        "\n[Info]   ",
	LogLvError:       "\n[Error]  ",
	LogLvDingMessage: "\n[Ding]   ",
	LogLvDingLists:   "\n[DAL!]   ",
	LogLvDingAll:     "\n[DAA!]   ",
}
var logFmtStr = "%sFile:%s%sLine:%d%s%s\n"

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

func SetUpEnv(m *EZLoggerModel) {
	appName = m.AppName
	logLevel = LogLv(m.LogLevel)
	enableDing = m.DingTalkModel.Enable
	if appName == "" {
		panic(fmt.Sprintf("AppName必须要配置,不然无法区分对应的服务"))
	}
	var err error
	if enableDing {
		if m.DingTalkModel.URLEncodedString == "" {
			panic(fmt.Sprintf("钉钉的那个URL没配置"))
		}
		dingURL, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.DingTalkModel.URLEncodedString)
		if err != nil {
			panic(fmt.Sprintf("获取密码错误:%s", err.Error()))
		}
		if m.DingTalkModel.SecretKeyEncodedString == "" {
			panic(fmt.Sprintf("钉钉的那个SecretKey没配置"))
		}
		dingSecretKey, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.DingTalkModel.SecretKeyEncodedString)
		if err != nil {
			panic(fmt.Sprintf("获取密码错误:%s", err.Error()))
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
			panic(fmt.Sprintf("gRPC日志服务没配地址"))
		}
		gRPCURL, err = ezPasswordEncoder.GetPasswordFromEncodedStr(m.GRPCModel.URLEncodedString)
		if err != nil {
			panic(fmt.Sprintf("获取密码错误:%s", err.Error()))
		}
		for i := 0; i < gRPCClientCounts; i++ {
			go startGRPCClient()
		}
	}
}
func D(msg ...any) {
	Log(LogLvDebug, msg...)
}
func I(msg ...any) {
	Log(LogLvInfo, msg...)
}
func E(msg ...any) {
	Log(LogLvError, msg...)
}
func DingMessage(msg ...any) {
	Log(LogLvDingMessage, msg...)
	if enableDing {
		sendToDing(LogLvDingMessage, "no tag", fmt.Sprintln(msg...))
	}
}
func DingAtAll(msg ...any) {
	Log(LogLvDingAll, msg...)
	if enableDing {
		sendToDing(LogLvDingAll, "no tag", fmt.Sprintln(msg...))
	}
}
func DingList(msg ...any) {
	Log(LogLvDingLists, msg...)
	if enableDing {
		sendToDing(LogLvDingLists, "no tag", fmt.Sprintln(msg...))
	}
}

func DWithTag(tag string, msg ...any) {
	LogWithTag(LogLvDebug, tag, msg...)
}
func IWithTag(tag string, msg ...any) {
	LogWithTag(LogLvInfo, tag, msg...)
}
func EWithTag(tag string, msg ...any) {
	LogWithTag(LogLvError, tag, msg...)
}
func DingMessageWithTag(tag string, msg ...any) {
	LogWithTag(LogLvDingMessage, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingMessage, tag, fmt.Sprintln(msg...))
	}
}
func DingAtAllWithTag(tag string, msg ...any) {
	LogWithTag(LogLvDingAll, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingAll, tag, fmt.Sprintln(msg...))
	}
}
func DingListWithTag(tag string, msg ...any) {
	LogWithTag(LogLvDingLists, tag, msg...)
	if enableDing {
		sendToDing(LogLvDingLists, tag, fmt.Sprintln(msg...))
	}
}
