package rsa

import "testing"

func TestGenerateRSAKeyPair(t *testing.T) {
	privatePem, publicPem, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Errorf("Error generating RSA key pair: %s", err)
		return
	}

	t.Log(privatePem, publicPem)

	privateKey, err := StripPemHeader(privatePem)
	if err != nil {
	}

	t.Logf("Private Key: %s", privateKey)
	publicKey, err := StripPemHeader(publicPem)
	if err != nil {
	}

	t.Logf("Public Key: %s", publicKey)

}

func TestGenerateRSAKeyPairPKCS8(t *testing.T) {
	privatePem, publicPem, err := GenerateRSAKeyPairPKCS8(2048)
	if err != nil {
		t.Errorf("Error generating RSA key pair: %s", err)
		return
	}
	t.Log(privatePem, publicPem)

	privateKey, err := StripPemHeader(privatePem)
	if err != nil {
	}

	t.Logf("Private Key: %s", privateKey)
	publicKey, err := StripPemHeader(publicPem)
	if err != nil {
	}

	t.Logf("Public Key: %s", publicKey)

}
