package ssp_test

import (
	"testing"

	ssp "github.com/christianschmizz/go-smiley"
	"github.com/stretchr/testify/assert"
)

func TestNewSetupRequestResponse(t *testing.T) {
	p, err := ssp.Decode([]byte{0x7F, 0x80, 0x17, 0xF0, 0x00, 0x30, 0x31, 0x30, 0x30, 0x45, 0x55, 0x52, 0x00, 0x00, 0x01, 0x03, 0x05, 0x0A, 0x14, 0x02, 0x02, 0x02, 0x00, 0x00, 0x64, 0x04, 0x2A, 0x25})
	assert.NoError(t, err)
	resp := ssp.NewSetupRequestResponse(p)
	assert.Equal(t, ssp.BanknoteValidator, resp.UnitType)
	assert.Equal(t, "EUR", resp.CountryCode)
	assert.Equal(t, "0100", resp.FirmwareVersion)

	validator, err := resp.BanknoteValidator()
	assert.NoError(t, err)

	assert.Equal(t, uint32(1), validator.ValueMultiplier)
	assert.Equal(t, uint8(3), validator.NumberOfChannels)
	assert.Equal(t, uint32(100), validator.RealValueMultiplier)
	assert.Equal(t, uint8(4), validator.ProtocolVersion)
	assert.Equal(t, validator.NumberOfChannels, uint8(len(validator.Channels)))
}
