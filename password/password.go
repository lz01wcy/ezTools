package password

import (
	"encoding/base64"
	"github.com/Anveena/ezTools/crypto"
)

func Encode(origPwd string) (string, error) {
	origPwdData := []byte(origPwd)
	encData, err := crypto.EZEncrypt(origPwdData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encData), nil
}
func Decode(base64Str string) (string, error) {
	encData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	origData, err := crypto.EZDecrypt(encData, "this code may be not working", 9458)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}
