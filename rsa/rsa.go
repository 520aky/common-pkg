package rsa

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
	privateKey = "-----BEGIN RSA PRIVATE KEY-----\r\n" + privateKey + "\r\n-----END RSA PRIVATE KEY-----"
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
	publicKey = "-----BEGIN PUBLIC KEY-----\r\n" + publicKey + "\r\n-----END PUBLIC KEY-----"
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

func GenerateRSAKeyPair(bits int) (privateKeyPEM string, publicKeyPEM string, err error) {
	// 1. 生成 RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 2. 私钥转 PKCS#1 PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}))

	// 3. 公钥转 PKCS#1 PEM
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}))

	return privateKeyPEM, publicKeyPEM, nil
}

func StripPemHeader(pemStr string) (string, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return "", errors.New("invalid PEM format")
	}

	// Base64 编码 PKCS1 私钥内容
	return base64.StdEncoding.EncodeToString(block.Bytes), nil
}

func GenerateRSAKeyPairPKCS8(bits int) (privateKeyPEM, publicKeyPEM string, err error) {
	// 1. 生成 RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 2. 私钥转 PKCS#8（注意这里不是 MarshalPKCS1PrivateKey）
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY", // PKCS#8 标准头
		Bytes: privBytes,
	}))

	// 3. 公钥转 PKIX (X.509 SubjectPublicKeyInfo)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}))

	return privateKeyPEM, publicKeyPEM, nil
}
