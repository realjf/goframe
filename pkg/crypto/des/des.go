package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

const (
	NOPADDING = iota
	PKCS5PADDING
)

func EncryptECB(plainText []byte, key []byte, padding int) ([]byte, error) {
	text, err := Padding(plainText, padding)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}

func DecryptECB(cipherText []byte, key []byte, padding int) ([]byte, error) {
	text := make([]byte, len(cipherText))
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	plainText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}


// The length of the <key> should be either 16 or 24 bytes.
func EncryptECBTriple(plainText []byte, key []byte, padding int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length error")
	}

	text, err := Padding(plainText, padding)
	if err != nil {
		return nil, err
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	cipherText := make([]byte, len(text))
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}


// The length of the <key> should be either 16 or 24 bytes.
func DecryptECBTriple(cipherText []byte, key []byte, padding int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length error")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	text := make([]byte, len(cipherText))
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	plainText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}


func EncryptCBC(plainText []byte, key []byte, iv []byte, padding int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text, err := Padding(plainText, padding)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, len(text))

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText, text)

	return cipherText, nil
}


func DecryptCBC(cipherText []byte, key []byte, iv []byte, padding int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	plainText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}


func EncryptCBCTriple(plainText []byte, key []byte, iv []byte, padding int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length invalid")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text, err := Padding(plainText, padding)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText, text)

	return cipherText, nil
}


func DecryptCBCTriple(cipherText []byte, key []byte, iv []byte, padding int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, errors.New("key length invalid")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("iv length invalid")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	plainText, err := UnPadding(text, padding)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func PaddingPKCS5(text []byte, blockSize int) []byte {
	padding := blockSize - len(text) % blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

func UnPaddingPKCS5(text []byte) []byte {
	length := len(text)
	padText := int(text[length-1])
	return text[:(length - padText)]
}

func Padding(text []byte, padding int) ([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text) % 8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return PaddingPKCS5(text, 8), nil
	default:
		return nil, errors.New("padding type error")
	}

	return text, nil
}

func UnPadding(text []byte, padding int) ([]byte, error) {
	switch padding {
	case NOPADDING:
		if len(text) % 8 != 0 {
			return nil, errors.New("text length invalid")
		}
	case PKCS5PADDING:
		return UnPaddingPKCS5(text), nil
	default:
		return nil, errors.New("padding type error")
	}
	return text, nil
}

