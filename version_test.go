package ssp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFirmwareVersionResponse(t *testing.T) {
	testcases := []struct {
		Name            string
		Data            []byte
		FirmwareVersion string
		DeviceVersion   string
		ReleaseVersion  string
		BetaVersion     string
	}{
		{"from docs", []byte{0x7F, 0x80, 0x11, 0xF0, 0x4E, 0x56, 0x30, 0x32, 0x30, 0x30, 0x34, 0x31, 0x34, 0x31, 0x34, 0x39, 0x38, 0x30, 0x30, 0x30, 0xDE, 0x55}, "NV02004141498000", "NV0200", "1498", "000"},
		{"from device", []byte{0x7F, 0x80, 0x11, 0xf0, 0x4e, 0x56, 0x53, 0x30, 0x30, 0x39, 0x31, 0x30, 0x39, 0x32, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x7a, 0x18}, "NVS0091092500000", "NVS009", "2500", "000"},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			p, err := Decode(tc.Data)
			assert.NoError(t, err)
			resp := NewVersionResponse(p)
			assert.Equal(t, tc.FirmwareVersion, resp.FirmwareVersion)
			assert.Equal(t, tc.DeviceVersion, resp.DeviceVersion)
			assert.Equal(t, tc.ReleaseVersion, resp.ReleaseVersion)
			assert.Equal(t, tc.BetaVersion, resp.BetaVersion)
		})
	}
}
