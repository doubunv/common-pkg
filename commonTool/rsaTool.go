package commonTool

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func GenerateSignature(secretKey string, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the private key")
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 签名
	hash := sha256.New()
	hash.Write([]byte(secretKey))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifySignature(secretKey string, signature string, publicKey string) (bool, error) {
	// 解析公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return false, errors.New("failed to parse PEM block containing the public key")
	}
	publicKey1, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPublicKey, ok := publicKey1.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not an RSA public key")
	}

	// 签名
	hash := sha256.New()
	hash.Write([]byte(secretKey))
	hashed := hash.Sum(nil)
	// 验证签名
	decodeString, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed, decodeString)
	if err != nil {
		return false, err
	}
	return true, nil
}
