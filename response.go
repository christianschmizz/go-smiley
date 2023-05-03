package ssp

type ResponseCode byte

// Generic response codes
const (
	OK                       ResponseCode = 0xF0
	CommandNotKnown          ResponseCode = 0xF2
	WrongNumberOfParameters  ResponseCode = 0xF3
	ParameterOutOfRange      ResponseCode = 0xF4
	CommandCannotBeProcessed ResponseCode = 0xF5
	SoftwareError            ResponseCode = 0xF6
	Fail                     ResponseCode = 0xF8
	KeyNotSet                ResponseCode = 0xFA
)

type Response struct {
	Code ResponseCode `json:"code"`
	Args []byte       `json:"-"`
}

func NewResponse(code ResponseCode, args []byte) *Response {
	return &Response{
		Code: code,
		Args: args,
	}
}

func (r *Response) String() string {
	switch r.Code {
	default:
		return "unknown code"
	case OK:
		return "ok"
	case CommandNotKnown:
		return "command not known"
	case WrongNumberOfParameters:
		return "wrong number of params"
	case ParameterOutOfRange:
		return "parameter out of range"
	case CommandCannotBeProcessed:
		return "cannot be processed"
	case SoftwareError:
		return "software error"
	case Fail:
		return "fail"
	case KeyNotSet:
		return "key not set"
	}
}
