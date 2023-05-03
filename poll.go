package ssp

import (
	"github.com/rs/zerolog/log"
)

type PollEventResponseCode byte

const (
	SlaveReset                   PollEventResponseCode = 0xF1
	ReadNote                     PollEventResponseCode = 0xEF
	CreditNote                   PollEventResponseCode = 0xEE
	NoteRejecting                PollEventResponseCode = 0xED
	NoteRejected                 PollEventResponseCode = 0xEC
	NoteStacking                 PollEventResponseCode = 0xCC
	NoteStacked                  PollEventResponseCode = 0xEB
	SafeNoteJam                  PollEventResponseCode = 0xEA
	UnsafeNoteJam                PollEventResponseCode = 0xE9
	Disabled                     PollEventResponseCode = 0xE8
	FraudAttempt                 PollEventResponseCode = 0xE6
	StackerFull                  PollEventResponseCode = 0xE7
	NoteClearedFromFront         PollEventResponseCode = 0xE1
	NoteClearedToCashbox         PollEventResponseCode = 0xE2
	CashboxRemoved               PollEventResponseCode = 0xE3
	CashboxReplaced              PollEventResponseCode = 0xE4
	BarCodeTicketValidated       PollEventResponseCode = 0xE5
	BarCodeTicketAcknowledge     PollEventResponseCode = 0xD1
	NotePathOpen                 PollEventResponseCode = 0xE0
	ChannelDisable               PollEventResponseCode = 0xB5
	Initialising                 PollEventResponseCode = 0xB6
	Dispensing                   PollEventResponseCode = 0xDA
	Dispensed                    PollEventResponseCode = 0xD2
	Jammed                       PollEventResponseCode = 0xD5
	Halted                       PollEventResponseCode = 0xD6
	Floating                     PollEventResponseCode = 0xD7
	Floated                      PollEventResponseCode = 0xD8
	Timeout                      PollEventResponseCode = 0xD9
	IncompletePayout             PollEventResponseCode = 0xDC
	IncompleteFloat              PollEventResponseCode = 0xDD
	CashboxPaid                  PollEventResponseCode = 0xDE
	CoinCredit                   PollEventResponseCode = 0xDF
	CoinMechJammed               PollEventResponseCode = 0xC4
	CoinMechReturnPressed        PollEventResponseCode = 0xC5
	Emptying                     PollEventResponseCode = 0xC2
	Emptied                      PollEventResponseCode = 0xC3
	SmartEmptying                PollEventResponseCode = 0xB3
	SmartEmptied                 PollEventResponseCode = 0xB4
	CoinMechError                PollEventResponseCode = 0xB7
	NoteStoredInPayout           PollEventResponseCode = 0xDB
	PayoutOutOfService           PollEventResponseCode = 0xC6
	JamRecovery                  PollEventResponseCode = 0xB0
	ErrorDuringPayout            PollEventResponseCode = 0xB1
	NoteTransferedToStacker      PollEventResponseCode = 0xC9
	NoteHeldInBezel              PollEventResponseCode = 0xCE
	NotePaidIntoStoreAtPowerUp   PollEventResponseCode = 0xCB
	NotePaidIntoStackerAtPowerUp PollEventResponseCode = 0xCA
	NoteDispensedAtPowerUp       PollEventResponseCode = 0xCD
	NoteFloatRemoved             PollEventResponseCode = 0xC7
	NoteFloatAttached            PollEventResponseCode = 0xC8
	DeviceFull                   PollEventResponseCode = 0xC9
)

var (
	pollEventResponseCodeDesc = map[PollEventResponseCode]string{
		SafeNoteJam:   "The note is stuck in a position not retrievable from the front of the device (user side)",
		UnsafeNoteJam: "The note is stuck in a position where the user could possibly remove it from the front of the device.",
		StackerFull:   "The banknote stacker unit attached to this device has been detected as at its full limit",
		SlaveReset:    "The device has undergone a power reset.",
		Disabled:      "The device is not active and unavailable for normal validation functions.",
	}
)

func (c PollEventResponseCode) Desc() string {
	s, ok := pollEventResponseCodeDesc[c]
	if !ok {
		return "unknown error"
	}
	return s
}

type PollEvent struct {
	Code PollEventResponseCode
	Args []byte
}

func NewPollEvent(cmd byte, args []byte) *PollEvent {
	return &PollEvent{
		Code: PollEventResponseCode(cmd),
		Args: args,
	}
}

type PollResponse struct {
	Response
}

func (c *Connection) Poll() (*PollResponse, error) {
	packet, err := c.execute(Command{Code: Poll})
	if err != nil {
		return nil, err
	}

	resp := packet.Response()
	if resp.Code != OK {
		log.Warn().Msg("poll NOK")
	}

	return &PollResponse{*resp}, nil
}
