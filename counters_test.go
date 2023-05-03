package ssp_test

import (
	"testing"

	ssp "github.com/christianschmizz/go-smiley"
	"github.com/stretchr/testify/assert"
)

func TestNewGetCountersResponse(t *testing.T) {
	p, err := ssp.Decode([]byte{0x7F, 0x80, 0x16, 0xF0, 0x05, 0x2C, 0x01, 0x00, 0x00, 0xD2, 0x00, 0x00, 0x00, 0xB4, 0x00, 0x00, 0x00, 0x68, 0x01, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0xF1, 0x82})
	assert.NoError(t, err)
	resp := ssp.NewGetCountersResponse(p)
	assert.Equal(t, &ssp.GetCountersResponse{
		Response: ssp.Response{
			Code: ssp.OK,
			Args: p.Command().Args,
		},
		NumberOfCountersSet: 5,
		NotesStacked:        300,
		NotesStored:         210,
		NotesDispensed:      180,
		NotesTransferred:    360,
		NotesRejected:       25,
	}, resp)
}
