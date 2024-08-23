package ujwt

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Get path working directory
func getWokingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	var cwd = getWokingDir()

	// Read private key from file
	privateKeyBytes, err := os.ReadFile(fmt.Sprintf("%s%s", cwd, "/configs/keys/private.key"))
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	var cwd = getWokingDir()
	publicKeyBytes, err := os.ReadFile(fmt.Sprintf("%s%s", cwd, "/configs/keys/public.key"))
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
