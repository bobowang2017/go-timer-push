package tests

import (
	"fmt"
	"go-timer-push/config"
	"go-timer-push/logger"
	"go-timer-push/utils/encrypt"
	"testing"
)

func SetUp() {
	config.LoadTestConfig()
	logger.Setup()
}

func TestAESBase64Encrypt(t *testing.T) {
	SetUp()
	aesClient := encrypt.NewAesCryptor()
	originData := "wangxiangbo"
	res, _ := aesClient.AESBase64Encrypt(originData)
	fmt.Println(res)
}

func TestAESBase64Decrypt(t *testing.T) {
	SetUp()
	aesClient := encrypt.NewAesCryptor()
	originData := "xZ7K5bvPMfiV0zQdxIvEdw"
	res, _ := aesClient.AESBase64Decrypt(originData)
	fmt.Println(res)
}
