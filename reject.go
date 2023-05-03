package ssp

type RejectCode uint8

const (
	NoteAccepted                       RejectCode = 0x00
	NoteLengthIncorrect                RejectCode = 0x01
	RejectReason2                      RejectCode = 0x02
	RejectReason3                      RejectCode = 0x03
	RejectReason4                      RejectCode = 0x04
	RejectReason5                      RejectCode = 0x05
	ChannelInhibited                   RejectCode = 0x06
	SecondNoteInserted                 RejectCode = 0x07
	RejectReason8                      RejectCode = 0x08
	NoteRecognisedInMoreThanOneChannel RejectCode = 0x09
	RejectReason10                     RejectCode = 0x0A
	NoteTooLong                        RejectCode = 0x0B
	RejectReason12                     RejectCode = 0x0C
	MechanismSlowOrStalled             RejectCode = 0x0D
	StrimmingAttemptDetected           RejectCode = 0x0E
	FraudChannelReject                 RejectCode = 0x0F
	NoNotesInserted                    RejectCode = 0x10
	PeakDetectFail                     RejectCode = 0x11
	TwistedNoteDetected                RejectCode = 0x12
	EscrowTimeout                      RejectCode = 0x13
	BarCodeScanFail                    RejectCode = 0x14
	RearSensor2Fail                    RejectCode = 0x15
	SlotFail1                          RejectCode = 0x16
	SlotFail2                          RejectCode = 0x17
	LensOversample                     RejectCode = 0x18
	WidthDetectFail                    RejectCode = 0x19
	ShortNoteDetected                  RejectCode = 0x1A
	NotePayout                         RejectCode = 0x1B
	UnableToStackNote                  RejectCode = 0x1C
)
