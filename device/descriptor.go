package device

type Descriptor struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	SupportedDevices   []Device `json:"supportedDevices"`
	EncryptionRequired bool     `json:"encryptionRequired"`
}
