package captcha

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Anveena/ezTools/captcha/ezCaptchaPB"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	_crtPath    = ""
	_serverName = ""
	_serverUrl  = ""
)

func SetGrpcPathAndCert(serverUrl string, serverName string, crtPath string) {
	_crtPath = crtPath
	_serverUrl = serverUrl
	_serverName = serverName
}
func Get() ([]byte, string, error) {
	crt, err := credentials.NewClientTLSFromFile(_crtPath, _serverName)
	if err != nil {
		return nil, "", err
	}
	clientConn, err := grpc.Dial(_serverUrl, grpc.WithTransportCredentials(crt))
	if err != nil {
		return nil, "", err
	}
	defer func() {
		_ = clientConn.Close()
	}()
	client := ezCaptchaPB.NewEZCaptchaServiceClient(clientConn)
	rsp, err := client.GetCaptcha(context.Background(), &ezCaptchaPB.EZCaptchaEmpty{})
	if err != nil {
		return nil, "", err
	}
	if !rsp.Suc {
		return nil, "", fmt.Errorf("对方说不行:%s", rsp.ErrDesc)
	}
	data, err := base64.StdEncoding.DecodeString(rsp.PngBase64)
	if err != nil {
		return nil, "", err
	}
	return data, rsp.CorrectAnswer, nil
}
