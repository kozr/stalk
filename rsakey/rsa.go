package rsakey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"
	"time"
)

var (
	currentPublicKey *rsa.PublicKey
	lock             sync.RWMutex
)

func Init() {
	// Generate initial key
	err := rotateKey()
	if err != nil {
		panic(err)
	}

	// Schedule key rotation
	go scheduleKeyRotation()
}

func scheduleKeyRotation() {
	for {
		// Calculate duration until next midnight
		nextMidnight := time.Now().Add(time.Hour * 24)
		nextMidnight = time.Date(nextMidnight.Year(), nextMidnight.Month(), nextMidnight.Day(), 0, 0, 0, 0, nextMidnight.Location())
		durationUntilMidnight := time.Until(nextMidnight)

		// Wait until midnight
		time.Sleep(durationUntilMidnight)

		// Rotate key
		err := rotateKey()
		if err != nil {
			fmt.Println("Error rotating key:", err)
			continue
		}
	}
}

func rotateKey() error {
	pub, err := generatePublicKey()
	if err != nil {
		return fmt.Errorf("error generating public key: %w", err)
	}

	lock.Lock()
	currentPublicKey = pub
	lock.Unlock()

	return nil
}

func generatePublicKey() (*rsa.PublicKey, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %w", err)
	}

	// Extract public key
	publicKey := privateKey.PublicKey

	return &publicKey, nil
}

func publicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	// Marshal public key to ASN.1 DER
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("error marshaling public key to ASN.1 DER: %w", err)
	}

	// Convert to PEM
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	lock.RLock()
	defer lock.RUnlock()

	return currentPublicKey, nil
}

func GetPublicKeyPEM() (string, error) {
	lock.RLock()
	defer lock.RUnlock()

	return publicKeyToPEM(currentPublicKey)
}
