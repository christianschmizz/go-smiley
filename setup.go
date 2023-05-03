package ssp

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type SetupRequestResponse struct {
	UnitType UnitType `json:"unit_type"`

	// ASCII data of device's firmware (eg. 0123)
	FirmwareVersion string `json:"firmware_version"`

	// The ASCII code of the device dataset (eg EUR)
	CountryCode string `json:"country_code"`

	data []byte
}

func (r *SetupRequestResponse) BanknoteValidator() (*BanknoteValidatorSetup, error) {
	if r.UnitType != BanknoteValidator {
		return nil, fmt.Errorf("This unit-type %s does not support banknote validator setups", r.UnitType)
	}

	b := r.data[0:3]
	valueMultiplier := uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
	if valueMultiplier == 0 {
		log.Debug().Msg("This is a protocol version 6 or greater compatible dataset where the values are given in the expanded segment of the return data.")
	}

	n := uint8(r.data[3]) // number of channels

	b2 := r.data[4+(2*n) : 4+(2*n)+3]
	realValueMultiplier := uint32(b2[2]) | uint32(b2[1])<<8 | uint32(b2[0])<<16

	resp := &BanknoteValidatorSetup{
		SetupRequestResponse: r,
		ValueMultiplier:      valueMultiplier,
		NumberOfChannels:     n,
		RealValueMultiplier:  realValueMultiplier,
		ProtocolVersion:      r.data[7+(2*n)],
		Channels:             make([]Channel, 0, n),
	}

	if resp.ValueMultiplier == 0 && resp.ProtocolVersion < 6 {
		// If the value multiplier is 0 then it indicates that this is a
		// protocol version 6 or greater compatible dataset where the values
		// are given in the expanded segment of the return data.
		log.Fatal().Uint32("value_multiplier", resp.ValueMultiplier).
			Uint8("protocol_version", resp.ProtocolVersion).
			Msg("for devices whose value multiplier is 0 a protocol version 6 or greater compatible dataset is expected.")
	}

	for i := uint8(0); i < n; i++ {
		channelValue := uint32(r.data[4+i])
		if resp.ValueMultiplier != 0 {
			channelValue = channelValue * resp.ValueMultiplier
		}
		channelSecurity := r.data[4+n+i]
		if channelSecurity != 2 && resp.ValueMultiplier > 0 {
			log.Fatal().Uint32("value_multiplier", resp.ValueMultiplier).
				Uint8("channel_security", channelSecurity).
				Msg("for value multiplier greater than 0 a security level of 2 is expected.")
		}
		ch := Channel{Value: channelValue, Security: channelSecurity}
		if resp.ProtocolVersion >= 6 {
			ch.ExpandedCountryCode = string(r.data[8+(2*n)+(i*3) : 8+(2*n)+(i*3)+3])
			ch.ExpandedValue = r.data[8+(5*n)+(i*4) : 8+(5*n)+(i*3)+4]
		}
		resp.Channels = append(resp.Channels, ch)
	}

	return resp, nil
}

type Channel struct {
	Value               uint32 `json:"value"`
	Security            uint8  `json:"security"`
	ExpandedCountryCode string `json:"expanded_country_code"`
	ExpandedValue       []byte `json:"expanded_value"`
}

type BanknoteValidatorSetup struct {
	*SetupRequestResponse
	ValueMultiplier     uint32 `json:"value_multiplier"`
	RealValueMultiplier uint32 `json:"real_value_multiplier"`
	ProtocolVersion     uint8  `json:"protocol_version"`

	// The highest channel used in this device dataset [n] (1-16)
	NumberOfChannels uint8 `json:"number_of_channels"`

	Channels []Channel `json:"channels"`
}

type SMARTHopperResponse struct {
	*SetupRequestResponse
}

func NewSetupRequestResponse(p *Packet) *SetupRequestResponse {
	data := p.Data()[1:]
	return &SetupRequestResponse{
		UnitType:        UnitType(data[0]),
		FirmwareVersion: string(data[1:5]),
		CountryCode:     string(data[5:8]),
		data:            data[8:],
	}
}

func (c *Connection) SetupRequest() (*SetupRequestResponse, error) {
	p, err := c.execute(Command{Code: SetupRequest})
	if err != nil {
		return nil, err
	}
	return NewSetupRequestResponse(p), nil
}
