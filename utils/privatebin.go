package utils

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/pbkdf2"
)

type PrivateBinData struct {
	Ct    string        `json:"ct"`
	Adata []interface{} `json:"adata"`
}

func padStart(s string, minLength int, padRune rune) string {
	currentLength := len(s)
	if currentLength >= minLength {
		return s
	}
	padding := strings.Repeat(string(padRune), minLength-currentLength)
	return padding + s
}

func DecryptPrivateBin(url string, password string) (string, error) {
	if !strings.Contains(url, "#") {
		return "", errors.New("Missing Decrypt Key")
	}
	key := strings.Split(url, "#")[1]
	resp, err := Fetch(FetchConfig{
		Url: url,
		Headers: map[string]string{
			"Accept": "application/json, text/javascript, */*; q=0.01",
		},
	})
	if err != nil {
		return "", err
	}
	data := PrivateBinData{}
	err = json.Unmarshal(resp.Data, &data)
	if err != nil {
		return "", err
	}
	type pasteJson struct {
		Paste string `json:"paste"`
	}
	ret, err := decryptPrivateBin(key, data.Adata, data.Ct, password)
	if err != nil {
		return "", err
	}
	var j pasteJson
	err = json.Unmarshal([]byte(ret), &j)
	if err != nil {
		return "", err
	}
	return j.Paste, nil
}

func decryptPrivateBin(key string, data []interface{}, cipherMessage, password string) (string, error) {
	decodedKey := base58.Decode(key)
	key = padStart(string(decodedKey), 32, '\x00')
	additionalData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	spec := data[0].([]interface{})
	iterations := int(spec[2].(float64))
	iv, err := base64.StdEncoding.DecodeString(spec[0].(string))
	if err != nil {
		return "", err
	}
	salt, err := base64.StdEncoding.DecodeString(spec[1].(string))
	if err != nil {
		return "", err
	}
	cipherMessageBytes, err := base64.StdEncoding.DecodeString(cipherMessage)
	if err != nil {
		return "", err
	}

	keyArray := []byte(key)
	if password != "" {
		if spec[7].(string) == "rawdeflate" {
			hash := sha256.New()
			hash.Write([]byte(password))
			password = hex.EncodeToString(hash.Sum(nil))
		}
		passwordArray := []byte(password)
		keyArray = append(keyArray, passwordArray...)
	}
	aesKeyLength := int(spec[3].(float64)) / 8
	deriveKey := pbkdf2.Key(keyArray, salt, iterations, aesKeyLength, sha256.New)
	block, err := aes.NewCipher(deriveKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return "", err
	}
	plaintext, err := aesGCM.Open(nil, iv, cipherMessageBytes, additionalData)
	if err != nil {
		return "", err
	}
	if len(spec) >= 8 && spec[7].(string) == "zlib" {
		data, err := decompress(plaintext)
		if err != nil {
			return "", err
		}
		plaintext = data
	}
	return string(plaintext), err
}

func decompress(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}
