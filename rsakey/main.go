package rsakey

var rotationService *KeyRotationService

func Init() {
	rotationService = NewKeyRotationService()
}

func GetRotationService() *KeyRotationService {
	return rotationService
}
