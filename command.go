package ssp


type CommandCode byte

const (
	ResetFixedEncryptionKey       CommandCode = 97
	SetFixedEncryptionKey         CommandCode = 96
	EnablePayoutDevice            CommandCode = 92
	DisablePayoutDevice           CommandCode = 91
	CoinMechOptions               CommandCode = 90
	ResetCounters                 CommandCode = 89
	GetCounters                   CommandCode = 88
	EventACK                      CommandCode = 87
	PollWithACK                   CommandCode = 86
	ConfigureBezel                CommandCode = 84
	CashboxPayoutOperationData    CommandCode = 83
	SmartEmpty                    CommandCode = 82
	GetHopperOptions              CommandCode = 81
	SetHopperOptions              CommandCode = 80
	GetBuildRevision              CommandCode = 79
	SetBaudRate                   CommandCode = 77
	RequestKeyExchange            CommandCode = 76
	SetModulus                    CommandCode = 75
	SetGenerator                  CommandCode = 74
	SetCoinMechGlobalInhibit      CommandCode = 73
	PayoutByDenomination          CommandCode = 70
	SetValueReportingType         CommandCode = 69
	FloatByDenomination           CommandCode = 68
	StackNote                     CommandCode = 67
	PayoutNote                    CommandCode = 66
	GetNotePositions              CommandCode = 65
	SetCoinMechInhibits           CommandCode = 64
	EmptyAll                      CommandCode = 63
	GetMinimumPayout              CommandCode = 62
	FloatAmount                   CommandCode = 61
	GetDenominationRoute          CommandCode = 60
	SetDenominationRoute          CommandCode = 59
	HaltPayout                    CommandCode = 56
	CommunicationPassThrough      CommandCode = 55
	GetDenominationLevel          CommandCode = 53
	SetDenominationLevel          CommandCode = 52
	PayoutAmount                  CommandCode = 51
	SetRefillMode                 CommandCode = 48
	GetBarCodeData                CommandCode = 39
	SetBarCodeInhibitStatus       CommandCode = 38
	GetBarCodeInhibitStatus       CommandCode = 37
	SetBarCodeConfiguration       CommandCode = 36
	GetBarCodeReaderConfiguration CommandCode = 35
	GetAllLevels                  CommandCode = 34
	GetDatasetVersion             CommandCode = 33
	GetFirmwareVersion            CommandCode = 32
	Hold                          CommandCode = 24
	LastRejectedCode              CommandCode = 23
	Sync                          CommandCode = 17
	ChannelReTeachData            CommandCode = 16
	ChannelSecurityData           CommandCode = 15
	ChannelValueRequest           CommandCode = 14
	UnitData                      CommandCode = 13
	GetSerialNumber               CommandCode = 12
	Enable                        CommandCode = 10
	Disable                       CommandCode = 9
	RejectBanknote                CommandCode = 8
	Poll                          CommandCode = 7
	HostProtocolVersion           CommandCode = 6
	SetupRequest                  CommandCode = 5
	DisplayOff                    CommandCode = 4
	DisplayOn                     CommandCode = 3
	SetChannelInhibits            CommandCode = 2
	Reset                         CommandCode = 1
)

var (
	genericCommandCodes = []CommandCode{
		Reset, HostProtocolVersion, GetSerialNumber, Sync,
		Disable, Enable, GetFirmwareVersion, GetDatasetVersion}
)

type Command struct {
	Code CommandCode
	Args []byte
}

func NewCommand(code CommandCode, args []byte) *Command {
	return &Command{
		Code: code,
		Args: args,
	}
}
