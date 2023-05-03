package device

import (
	ssp "github.com/christianschmizz/go-smiley"
)

var (
	supportedCommands = map[ssp.CommandCode]Descriptor{
		ssp.ResetFixedEncryptionKey: {SupportedDevices: []Device{SmartHopper, SmartPayout, NV11}, Description: "", EncryptionRequired: false},
		ssp.GetSerialNumber:         {SupportedDevices: []Device{NV9USB, NV10USB, BV20, BV50, BV100, NV200, SmartHopper, SmartPayout, NV11}, Description: "", EncryptionRequired: false},
		ssp.Poll:                    {SupportedDevices: []Device{NV9USB, NV10USB, BV20, BV50, BV100, NV200, SmartHopper, SmartPayout, NV11}, Description: "", EncryptionRequired: false},
		ssp.DisplayOff:              {SupportedDevices: []Device{NV9USB, NV10USB, NV200, NV11}, Description: "", EncryptionRequired: false},
		ssp.DisplayOn:               {SupportedDevices: []Device{NV9USB, NV10USB, NV200, NV11}, Description: "", EncryptionRequired: false},
		ssp.Sync:                    {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.GetCounters:             {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.Enable:                  {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.SetupRequest:            {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.SetChannelInhibits:      {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.GetFirmwareVersion:      {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.Reset:                   {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.GetDatasetVersion:       {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
		ssp.PollWithACK:             {SupportedDevices: []Device{}, Description: "", EncryptionRequired: false},
	}
)

func (d Device) Supports(commands ...ssp.CommandCode) bool {
	for _, cmd := range commands {
		supported := false
		for _, dev := range supportedCommands[cmd].SupportedDevices {
			if dev == d {
				supported = true
			}
			if !supported {
				return false
			}
		}
	}
	return true
}
