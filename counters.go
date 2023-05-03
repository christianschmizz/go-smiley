package ssp

import (
	"encoding/binary"
)


type GetCountersResponse struct {
	Response
	NumberOfCountersSet uint8  `json:"number_of_counters_set"`
	NotesStacked        uint32 `json:"notes_stacked_count"`
	NotesStored         uint32 `json:"notes_stored_count"`
	NotesDispensed      uint32 `json:"notes_dispensed_count"`
	NotesTransferred    uint32 `json:"notes_transferred_count"`
	NotesRejected       uint32 `json:"notes_rejected_count"`
}

func NewGetCountersResponse(p *Packet) *GetCountersResponse {
	payload := p.Data()[1:]
	var (
		numberOfCountersSet uint8 = payload[0]
		notesStacked              = binary.LittleEndian.Uint32(payload[1:5])
		notesStored               = binary.LittleEndian.Uint32(payload[5:9])
		notesDispensed            = binary.LittleEndian.Uint32(payload[9:13])
		notesTransferred          = binary.LittleEndian.Uint32(payload[13:17])
		notesRejected             = binary.LittleEndian.Uint32(payload[17:22])
	)
	return &GetCountersResponse{
		Response: Response{
			Code: ResponseCode(p.Data()[0]),
			Args: p.Data()[1:],
		},
		NumberOfCountersSet: numberOfCountersSet,
		NotesStacked:        notesStacked,
		NotesStored:         notesStored,
		NotesDispensed:      notesDispensed,
		NotesTransferred:    notesTransferred,
		NotesRejected:       notesRejected,
	}
}


func (c *Connection) GetCounters() (*GetCountersResponse, error) {
	p, err := c.execute(Command{Code: GetCounters})
	if err != nil {
		return nil, err
	}
	return NewGetCountersResponse(p), nil
}
