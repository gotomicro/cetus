package x

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"hash/fnv"
	"math/rand"
	"regexp"
	"time"
)

var regInputCheck = regexp.MustCompile(`^[a-zA-Z0-9].$`)

func InputCheck(input string) bool {
	return regInputCheck.MatchString(input)
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// 生成随机字符串
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// AESEncrypt Example of use
//
//	key := "myverystrongpasswordo32bitlength"
//	password := "mysecretpassword"
//
//	encrypted, err := AESEncrypt(password, key)
//	if err != nil {
//		fmt.Println("Error encrypting:", err)
//		return
//	}
//	fmt.Println("Encrypted:", encrypted)
func AESEncrypt(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, []byte(key)[:aes.BlockSize])
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt Example of use
//
//	decrypted, err := AESDecrypt(encrypted, key)
//	if err != nil {
//		fmt.Println("Error decrypting:", err)
//		return
//	}
//	fmt.Println("Decrypted:", decrypted)
func AESDecrypt(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherText, _ := base64.StdEncoding.DecodeString(text)
	cfb := cipher.NewCFBDecrypter(block, []byte(key)[:aes.BlockSize])
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
