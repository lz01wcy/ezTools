package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetMD5Data(data []byte) []byte {
	md5Maker := md5.New()
	md5Maker.Write(data)
	return md5Maker.Sum(nil)
}
func GetMD5Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(GetMD5Data(data))
}
func GetMD5HexString(data []byte) string {
	return hex.EncodeToString(GetMD5Data(data))
}

func GetSHA1Data(data []byte) []byte {
	sha1Maker := sha1.New()
	sha1Maker.Write(data)
	return sha1Maker.Sum(nil)
}
func GetSHA1Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(GetSHA1Data(data))
}
func GetSHA1HexString(data []byte) string {
	return hex.EncodeToString(GetSHA1Data(data))
}

func GetSHA256Data(data []byte) []byte {
	sha256Maker := sha256.New()
	sha256Maker.Write(data)
	return sha256Maker.Sum(nil)
}
func GetSHA256Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(GetSHA256Data(data))
}
func GetSHA256HexString(data []byte) string {
	return hex.EncodeToString(GetSHA256Data(data))
}

func GetHMACSHA1Data(data []byte, key []byte) []byte {
	hmacMaker := hmac.New(sha1.New, key)
	hmacMaker.Write(data)
	return hmacMaker.Sum(nil)
}
func GetHMACSHA1Base64(data []byte, key []byte) string {
	return base64.StdEncoding.EncodeToString(GetHMACSHA1Data(data, key))
}
func GetHMACSHA1HexString(data []byte, key []byte) string {
	return hex.EncodeToString(GetHMACSHA1Data(data, key))
}

func GetHMACSHA256Data(data []byte, key []byte) []byte {
	hmacMaker := hmac.New(sha256.New, key)
	hmacMaker.Write(data)
	return hmacMaker.Sum(nil)
}
func GetHMACSHA256Base64(data []byte, key []byte) string {
	return base64.StdEncoding.EncodeToString(GetHMACSHA256Data(data, key))
}
func GetHMACSHA256HexString(data []byte, key []byte) string {
	return hex.EncodeToString(GetHMACSHA256Data(data, key))
}
