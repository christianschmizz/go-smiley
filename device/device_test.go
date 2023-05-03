package device

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFirmwareVersion(t *testing.T) {
	testcases := []struct {
		FirmwareVersion string
		ExpectedDevice  Device
		Error error
	}{
		{"NVS0091092500000", NV9USB, nil},
		{"NV02004141498000", NV200, nil},
		{"XXX", Unknown, errors.New("invalid version string: XXX")},
		{"XXXXXX", Unknown, errors.New("failed to detect device from version string: XXXXXX")},
	}
	for _, tc := range testcases {
		t.Run("detect device from "+tc.FirmwareVersion, func(t *testing.T) {
			dev, err := FromFirmwareVersion(tc.FirmwareVersion)
			if tc.Error != nil {
				assert.Error(t, tc.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedDevice, dev)
			}
		})
	}
}
