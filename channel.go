package ssp

import (
	"encoding/binary"
)

type setChannelInhibitsResponse struct{
	Response
}

func NewSetChannelInhibitsResponse(p *Packet) *setChannelInhibitsResponse {
	//_ := p.Data()[1:]
	return &setChannelInhibitsResponse{}
}

func (c *Connection) SetChannelInhibits(channels uint16) (*setChannelInhibitsResponse, error) {
	args := make([]byte, 2)
	binary.LittleEndian.PutUint16(args, channels)
	p, err := c.execute(Command{SetChannelInhibits, args})
	if err != nil {
		return nil, err
	}
	return NewSetChannelInhibitsResponse(p), nil
}
