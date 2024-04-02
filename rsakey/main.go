package rsakey

var rotationService *KeyRotationService

func Init() error {
	rotationService = NewKeyRotationService()

	return nil
}

func GetRotationService() *KeyRotationService {
	return rotationService
}
