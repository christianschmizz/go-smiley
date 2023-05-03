package device

import "fmt"

type Device string

const (
	Unknown     Device = ""
	NV9USB      Device = "NV9USB"
	NV10USB     Device = "NV10USB"
	BV20        Device = "BV20"
	BV50        Device = "BV50"
	BV100       Device = "BV100"
	NV200       Device = "NV200"
	NV11        Device = "NV11"
	SmartHopper Device = "SMART Hopper"
	SmartPayout Device = "SMART Payout"
)

func FromFirmwareVersion(version string) (Device, error) {
	if len(version) < 6 {
		return Unknown, fmt.Errorf("invalid version string: %s", version)
	}
	switch version[0:6] {
	case "NVS009":
		return NV9USB, nil
	case "NV0200":
		return NV200, nil
	}
	return Unknown, fmt.Errorf("failed to detect device from version string: %s", version)
}
