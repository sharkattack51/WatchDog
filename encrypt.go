package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
	"os"
)

var KEY = "0123456789abcdef"
var COMMON_IV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

// Mailパスワードを暗号化して保存
func EncryptToFile(path string, pwd string) error {
	c, err := aes.NewCipher([]byte(KEY))
	if err != nil {
		return err
	}

	cfbenc := cipher.NewCFBEncrypter(c, COMMON_IV)
	enc := make([]byte, len([]byte(pwd)))
	cfbenc.XORKeyStream(enc, []byte(pwd))

	ioutil.WriteFile(path, enc[:], os.ModePerm)
	return nil
}

// パスワードをファイルから複合化
func DecryptFromFile(path string) (string, error) {
	enc, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(KEY))
	if err != nil {
		return "", err
	}

	cfbdec := cipher.NewCFBDecrypter(c, COMMON_IV)
	dec := make([]byte, len(enc))
	cfbdec.XORKeyStream(dec, enc)

	return string(dec), nil
}
