package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var key_ string

func AESLoadKey() {
	// key_ = viper.GetString("secrets.aesgcm-key1")
	key_ = "secretAsdncHlaskjdk"
}

func NewEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}

// var NameEncrypt = func(buyerName string) string {
// 	buyerNameEncrypt, _ := Encrypt([]byte(buyerName))
// 	buyerNamestr := fmt.Sprintf("%x", buyerNameEncrypt)
// 	return buyerNamestr
// }

var Encrypt = func(plaintext []byte) (ciphertext []byte, err error) {
	key := []byte(key_)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	key := []byte(key_)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("error cipherText length invalid")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func PlainTextToCipherText(plainText string) (string, error) {
	cipherText := plainText
	dataPlainText, err := hex.DecodeString(cipherText)
	if err != nil {

		return "", errors.New("cannot decode string to byte error message:" + err.Error())
	}
	_, err = Decrypt(dataPlainText)
	if err != nil {
		dataEncrypt, err := Encrypt([]byte(plainText))
		if err != nil {
			return "", errors.New("cannot encrypt plainText to cipherText error message:" + err.Error())
		}
		cipherText = fmt.Sprintf("%x", dataEncrypt)
	}
	return cipherText, nil
}

func DecryptTest(cipherText string) (string, error) {
	data, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", errors.New("cannot decode string to byte error message:" + err.Error())
	}

	plainText, err := Decrypt(data)
	if err != nil {
		return "", errors.New("cannot decrypt cipherText to plainText error message:" + err.Error())
	}

	return fmt.Sprintf("%s\n", string(plainText)), nil
}
