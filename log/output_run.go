//go:build !ezDebug

package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Anveena/ezTools/hash"
	"github.com/Anveena/ezTools/log/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func startGRPCClient() {
	runtime.LockOSThread()
	ctx := context.Background()
	errMsg := ""
	defer func() {
		if enableDing {
			DingAtAllWithTag("ezlog", errMsg)
		}
		fmt.Printf("日志爆炸了!!错误:\n\t%s\n", errMsg)
		for {
			<-logChannel
		}
	}()
	clientConn, err := grpc.Dial(gRPCURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errMsg = fmt.Sprintf("grpc拨号错误:%s", err.Error())
		return
	}
	grpcClient := model.NewEzLogGrpcClient(clientConn)
	stream, err := grpcClient.Log(ctx)
	if err != nil {
		errMsg = fmt.Sprintf("grpc方法调用错误:%s", err.Error())
		return
	}
	// It is safe to have a goroutine calling SendMsg and another goroutine
	// calling RecvMsg on the same stream at the same time, but it is not safe
	// to call SendMsg on the same stream in different goroutines. It is also
	// not safe to call CloseSend concurrently with SendMsg.
	// 大老师说这个东西不能异步随便call 所以需要用chan
	for {
		msg := <-logChannel
		if err = stream.Send(msg); err != nil {
			errMsg = fmt.Sprintf("发送失败:%s\n文件:%s\n行号:%d\n消息内容:\n%s", err.Error(), msg.FileName, msg.FileLine, msg.Content)
			return
		}
	}
}
func Log(level LogLv, msg ...any) {
	if level < logLevel {
		return
	}
	_, file, line, _ := runtime.Caller(2)
	logChannel <- &model.EZLogReq{
		Level:    int32(level),
		FileLine: int32(line),
		Time:     timestamppb.New(time.Now()),
		FileName: file,
		AppName:  appName,
		Content:  fmt.Sprint(msg...),
	}
}
func LogWithTag(level LogLv, tag string, msg ...any) {
	if level < logLevel {
		return
	}
	_, file, line, _ := runtime.Caller(2)
	logChannel <- &model.EZLogReq{
		Level:    int32(level),
		FileLine: int32(line),
		Time:     timestamppb.New(time.Now()),
		FileName: file,
		AppName:  appName,
		Tag:      tag,
		Content:  fmt.Sprint(msg...),
	}
}
func sendToDing(logLv LogLv, tag string, msg string) {
	_, file, line, _ := runtime.Caller(2)
	fileSubArr := strings.SplitN(file, "/src/", 2)
	if len(fileSubArr) == 2 {
		file = fileSubArr[1]
	}
	realMsg := strings.ReplaceAll(msg, "\n", "\n>\n>")
	atInfo := struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	}{}
	var atStr string
	if logLv == LogLvDingLists && len(dingMobiles) == 0 {
		logLv = LogLvDingAll
	}
	if logLv == LogLvDingAll {
		atInfo.IsAtAll = true
	} else if logLv == LogLvDingLists {
		atStr = dingAtStr
		atInfo.AtMobiles = dingMobiles
	}
	m := dingRequestModel{
		MsgType: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: fmt.Sprintf("%s@%s", tag, appName),
			Text:  fmt.Sprintf("# %s@%s\n### **File**:%s\n### **Line**:%d\n### **Time**:%s\n### **Message**:\n>%s \n%s", tag, appName, file, line, time.Now().Format("15:04:05.999999"), realMsg, atStr),
		},
		At: atInfo,
	}
	tmpData, err := json.Marshal(m)
	if err != nil {
		E(err.Error())
		return
	}
	timestampString := strconv.Itoa(int(time.Now().Unix() * 1000))
	sign := hash.GetHMACSHA256Base64([]byte(fmt.Sprintf("%s\n%s", timestampString, dingSecretKey)), dingSecretKeyData)
	go func() {
		rsp, err := http.Post(fmt.Sprintf("%s&timestamp=%s&sign=%s", dingURL, timestampString, url.QueryEscape(sign)), "application/json;charset=utf-8", bytes.NewBuffer(tmpData))
		if err != nil {
			E(err.Error())
			return
		}
		rspInfo, _ := io.ReadAll(rsp.Body)
		rspStr := string(rspInfo)
		if !strings.Contains(rspStr, `"errcode":0,`) {
			E("rsp from ding:", rspStr)
		}
		_ = rsp.Body.Close()
	}()
}
