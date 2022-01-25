package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"go-timer-push/logger"
	"strconv"
	"time"
)

type AesCryptor struct {
	Key string
	Iv  string
}

func (a *AesCryptor) AESBase64Encrypt(originData string) (result string, err error) {
	var block cipher.Block
	iv := []byte(a.Iv)
	if block, err = aes.NewCipher([]byte(a.Key)); err != nil {
		logger.Logger.Error(err.Error())
		return "", err
	}
	encrypt := cipher.NewCBCEncrypter(block, iv)
	source := a.PKCS5Padding([]byte(originData), 16)
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	return base64.RawStdEncoding.EncodeToString(dst), nil
}

func (a *AesCryptor) AESBase64Decrypt(encryptData string) (originData string, err error) {
	var block cipher.Block
	iv := []byte(a.Iv)
	if block, err = aes.NewCipher([]byte(a.Key)); err != nil {
		logger.Logger.Error(err.Error())
		return "", err
	}
	encrypt := cipher.NewCBCDecrypter(block, iv)
	var source []byte
	if source, err = base64.RawStdEncoding.DecodeString(encryptData); err != nil {
		logger.Logger.Error(err.Error())
		return "", err
	}
	dst := make([]byte, len(source))
	encrypt.CryptBlocks(dst, source)
	return string(a.PKCS5UnPadding(dst)), nil
}

func (a *AesCryptor) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (a *AesCryptor) PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func NewAesCryptor() *AesCryptor {
	return &AesCryptor{
		Key: "1hA0BBiz1dDE3jmo",
		Iv:  "Qm3ZlGwT5Aa3WZO8",
	}
}

//测试生成MeterSphere签名方法
func generateToken() string {
	a := time.Now().UnixNano() / 1e6
	aesClient := AesCryptor{
		"1hA0BBiz1dDE3jmo",
		"Qm3ZlGwT5Aa3WZO8",
	}
	originData := "Qm3ZlGwT5Aa3WZO8|" + strconv.FormatInt(a, 10)
	res, _ := aesClient.AESBase64Encrypt(originData)
	return res
}
