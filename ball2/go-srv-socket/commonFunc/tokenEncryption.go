package commonFunc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"os"
)

type Token struct {
	Timestamp int64
	Content   string
}

func AesTokenEncrypt(token Token) string {

	timeByte := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeByte, uint64(token.Timestamp))
	plaintext := append(timeByte, []byte(token.Content)...)

	ciphertext := aesEncrypt(plaintext)
	str := base64.URLEncoding.EncodeToString(ciphertext)
	return str
}

func AesTokenDecrypt(str string) (Token, bool) {

	token := Token{}
	ciphertext, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return token, false
	}

	plaintext := aesDecrypt(ciphertext)
	if plaintext == nil || len(plaintext) < 8 {
		return token, false
	}
	timeByte := plaintext[:8]
	plaintext = plaintext[8:]

	token.Timestamp = int64(binary.LittleEndian.Uint64(timeByte))
	token.Content = string(plaintext)
	return token, true
}

func RandomByte(length int) []byte {
	a := make([]byte, length)
	for i := 0; i < length; i++ {
		a[i] = byte(rand.Intn(256))
	}
	return a
}

func aesEncrypt(plaintext []byte) []byte {

	padSize := aes.BlockSize - (len(plaintext))%aes.BlockSize
	padByte := bytes.Repeat([]byte{byte(padSize)}, padSize)
	plaintext = append(plaintext, padByte...)

	iv := RandomByte(aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))

	key := []byte(os.Getenv("aesKey"))
	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	ciphertext = append(iv, ciphertext...)
	return ciphertext
}

func aesDecrypt(ciphertext []byte) []byte {

	textSize := int(len(ciphertext) - aes.BlockSize)
	if textSize < 0 || textSize%aes.BlockSize != 0 {
		return nil
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, textSize)

	key := []byte(os.Getenv("aesKey"))
	block, _ := aes.NewCipher(key)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	padSize := int(plaintext[textSize-1])
	plaintext = plaintext[:textSize-padSize]
	return plaintext
}
