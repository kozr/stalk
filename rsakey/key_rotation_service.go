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

type KeyRotationService struct {
	currentPublicKeys []publicKeyWithTimestamp
	rotationInterval  time.Duration
	lock              sync.RWMutex
	stopChan          chan struct{}
	maxKeyAge         time.Duration
}

type publicKeyWithTimestamp struct {
	Key       *rsa.PublicKey
	Timestamp time.Time
	PEMString string
}

func NewKeyRotationService() *KeyRotationService {
	return &KeyRotationService{}
}

func (krs *KeyRotationService) SetRotationInterval(interval time.Duration) {
	krs.rotationInterval = interval
}

func (krs *KeyRotationService) SetMaxKeyAge(maxAge time.Duration) {
	krs.maxKeyAge = maxAge
}

func (krs *KeyRotationService) StartKeyRotation() error {
	if krs.stopChan != nil {
		return fmt.Errorf("there is already a key rotation in progress, call StopKeyRotation first")
	} else if krs.rotationInterval == 0 {
		return fmt.Errorf("rotation interval is not set, call SetRotationInterval first")
	}

	krs.stopChan = make(chan struct{})

	err := krs.rotateKey()
	if err != nil {
		panic(err)
	}

	go krs.rotationLoop()

	return nil
}

func (krs *KeyRotationService) StopKeyRotation() error {
	if krs.stopChan == nil {
		return fmt.Errorf("there is no key rotation in progress")
	}

	close(krs.stopChan)
	krs.stopChan = nil
	return nil
}

func (krs *KeyRotationService) rotationLoop() {
	for {
		nextStartTime := time.Now().Add(krs.rotationInterval)
		nextStartTime = time.Date(nextStartTime.Year(), nextStartTime.Month(), nextStartTime.Day(), 0, 0, 0, 0, nextStartTime.Location())
		duration := time.Until(nextStartTime)

		select {
		case <-time.After(duration):
			err := krs.rotateKey()
			if err != nil {
				fmt.Println("Error rotating key:", err)
				continue
			}
		case <-krs.stopChan:
			fmt.Println("Stopping key rotation")
			return
		}
	}
}

func (krs *KeyRotationService) rotateKey() error {
	krs.lock.Lock()
	defer krs.lock.Unlock()

	pub, err := krs.generatePublicKey()
	if err != nil {
		return fmt.Errorf("error generating public key: %w", err)
	}

	pemString := krs.publicKeyToPEM(pub) // Generate PEM string once and cache it
	// Append new key with current timestamp and its PEM string
	krs.currentPublicKeys = append(krs.currentPublicKeys, publicKeyWithTimestamp{
		Key:       pub,
		Timestamp: time.Now(),
		PEMString: pemString, // Assuming you add a PEMString field to your struct
	})

	krs.pruneOldKeys()

	return nil
}

func (krs *KeyRotationService) pruneOldKeys() {
	cutoff := time.Now().Add(-krs.maxKeyAge)
	var prunedKeys []publicKeyWithTimestamp
	for _, k := range krs.currentPublicKeys {
		if k.Timestamp.After(cutoff) {
			prunedKeys = append(prunedKeys, k)
		}
	}
	krs.currentPublicKeys = prunedKeys
}

func (krs *KeyRotationService) generatePublicKey() (*rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("error generating private key: %w", err)
	}

	return &privateKey.PublicKey, nil
}

func (krs *KeyRotationService) publicKeyToPEM(publicKey *rsa.PublicKey) string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic("error marshaling public key to ASN.1 DER: " + err.Error())
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM)
}

func (krs *KeyRotationService) GetLatestPublicKey() (string, error) {
	krs.lock.RLock()
	defer krs.lock.RUnlock()

	if len(krs.currentPublicKeys) == 0 {
		return "", fmt.Errorf("no keys available")
	}

	key := krs.currentPublicKeys[len(krs.currentPublicKeys)-1].PEMString

	// Return the most recently added key
	return key, nil
}

func (krs *KeyRotationService) GetAllPublicKeys() []string {
	krs.lock.RLock() // Ensure thread-safe access to the keys
	defer krs.lock.RUnlock()

	var keysPEM []string
	for _, k := range krs.currentPublicKeys {
		keysPEM = append(keysPEM, k.PEMString)
	}

	return keysPEM
}
