package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	nonceTxt = "46562ffddb9890d8eb946527"
)

func AESGCMEncrypt(plainText, key string) (string, error) {
	var (
		keyBytes = genKeyBytes(key)
	)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	nonce, _ := hex.DecodeString(nonceTxt)
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	cipherText := aesGcm.Seal(nil, nonce, []byte(plainText), nil)

	return fmt.Sprintf("%x", cipherText), nil
}

func AESGCMDecrypt(cipherText, key string) (string, error) {
	var (
		keyBytes = genKeyBytes(key)
	)

	cipherTextBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	nonce, _ := hex.DecodeString(nonceTxt)
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aesGcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// AES/CBC/PKCS5Padding encrypt/decrypt mode
func AESEncrypt(plaintext, key string) (string, error) {

	keyHash := formatKey(key)

	//	new ase cipher key hash
	block, err := aes.NewCipher([]byte(keyHash))
	if err != nil {
		return "", err
	}

	//	data padding with pkcs5 type
	dataPadding := []byte(plaintext)
	dataPadding = pkcs5Padding(dataPadding, block.BlockSize())
	if len(dataPadding)%aes.BlockSize != 0 {
		err = errors.New("plaintext is not a multiple of the block size")
		return "", err
	}

	//	ASE CBC encrypting
	encryptBuf := make([]byte, len(dataPadding))
	iv := []byte(keyHash)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encryptBuf, dataPadding)

	return base64.RawURLEncoding.EncodeToString(encryptBuf), nil
}

// AES/CBC/PKCS5Padding decrypt mode
func AESDecrypt(cipherTxt, key string) (plaintext string, err error) {
	keyFormatted := formatKey(key)

	block, err := aes.NewCipher([]byte(keyFormatted))
	if err != nil {
		return
	}

	cipherBuf, err := base64.RawURLEncoding.DecodeString(cipherTxt)
	if err != nil {
		return
	}

	if len(cipherBuf) < aes.BlockSize {
		err = errors.New("cipher text too short")
		return
	}

	iv := []byte(getVector(keyFormatted))
	dataUnPadding := make([]byte, len(cipherBuf))

	if len(cipherBuf)%aes.BlockSize != 0 {
		err = errors.New("cipher text is not a multiple of the block size")
		return
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dataUnPadding, cipherBuf)

	decryptBuf, err := pkcs5Trimming(dataUnPadding, mode.BlockSize())
	plaintext = string(decryptBuf)
	return
}

func genKeyBytes(keyString string) []byte {
	keyInSha := sha256.Sum256([]byte(keyString))
	keyStringInSha := fmt.Sprintf("%x", keyInSha)
	byteArr := []byte(keyStringInSha)
	keyBytes := make([]byte, aes.BlockSize*2)
	copy(keyBytes, byteArr)

	return keyBytes
}

func pkcs5Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

func pkcs5Trimming(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, errors.New("padding size error")
	}

	return src[:srcLen-paddingLen], nil
}

func formatKey(key string) string {
	if len(key) >= 32 {
		return key[0:32]
	}

	return key
}

func getVector(key string) string {
	length := len(key)
	if length >= 16 {
		return key[length-16 : length]
	}

	return key
}
