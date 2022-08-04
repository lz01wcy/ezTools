package ezCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"github.com/Anveena/ezTools/ezHash"
)

func freedomPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func freedomUnPadding(plainText []byte) []byte {
	length := len(plainText)
	number := int((plainText)[length-1])
	return (plainText)[:length-number]
}
func TripleDESEncrypt(origData []byte, keyStr string) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	key := []byte(keyStr)
	block, err := des.NewTripleDESCipher(key[:24])
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	tmpData := freedomPadding(origData, blockSize)
	tbc := make([]byte, len(tmpData))
	for i := 0; i+blockSize <= len(tmpData); i += blockSize {
		block.Encrypt(tbc[i:i+blockSize], tmpData[i:i+blockSize])
	}
	return tbc, nil
}
func MakeMD5Key(aesKey string, times int64) []byte {
	key := []byte(aesKey)
	for i := int64(0); i < times; i++ {
		key = ezHash.GetMD5Data(key)
	}
	return key
}
func AESCBCEncrypt(origData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	tmpData := freedomPadding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, (key)[:blockSize])
	crypted := make([]byte, len(tmpData))
	blockMode.CryptBlocks(crypted, tmpData)
	return crypted, nil
}
func AESCBCDecrypt(encData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, (key)[:blockSize])
	origData := make([]byte, len(encData))
	blockMode.CryptBlocks(origData, encData)
	origData = freedomUnPadding(origData)
	return origData, nil
}
func AESGCMEncrypt(origData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nil, key[8:20], origData, nil), nil
}
func AESGCMDecrypt(encData []byte, key []byte) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, key[8:20], encData, nil)
}
func EZEncrypt(origData []byte, ezKey string, salt uint64) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("很严重的错误呢,差点duang了%v", e)
		}
	}()
	aesKey := MakeMD5Key(ezKey, int64(salt%251))
	return AESCBCEncrypt(origData, aesKey)
}
func EZDecrypt(encData []byte, ezKey string, salt uint64) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	aesKey := MakeMD5Key(ezKey, int64(salt%251))
	return AESCBCDecrypt(encData, aesKey)
}
