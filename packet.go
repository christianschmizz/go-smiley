package ssp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog/log"
)

const (
	stx         byte = 0x7F
	stex        byte = 0x7E
	maxDataSize      = 0xFF
)

type PacketI interface {
	SequenceFlag() bool
	SlaveID() uint8
	Data() []byte
	Bytes() []byte
}

type Packet struct {
	seq    uint8 // bit 7 sequence flag, bits 6-0 represent the address of the slave
	code   byte  // command or response code
	args   []byte
	secure bool
}

func NewPacket(code byte, args []byte, seq uint8) *Packet {
	return &Packet{
		seq:    seq,
		code:   code,
		args:   args,
		secure: false,
	}
}

func (p *Packet) SequenceFlag() bool {
	return p.seq&seqIDMask > 0
}

func (p *Packet) SlaveID() uint8 {
	return p.seq & ^seqIDMask
}

func (p *Packet) Command() *Command {
	return NewCommand(CommandCode(p.code), p.args)
}

func (p *Packet) Response() *Response {
	return NewResponse(ResponseCode(p.code), p.args)
}

func (p *Packet) encryptedData() []byte {
	data := p.Data()

	// @todo counter
	var count uint32
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, count)

	encData := append([]byte{byte(len(data))}, a...)
	encData = append(encData, data...)

	// packing = count (4 bytes) + len(data) + packing (? bytes) crcl (1 byte) + crch (1 byte)
	missingBytes := 16 % (6 + len(data))
	if missingBytes != 0 {
		packBytes := make([]byte, missingBytes)
		rand.Read(packBytes)
		encData = append(encData, packBytes...)
	}

	// Append checksums
	eCRCLow, eCRCHigh := crc16(encData)
	encData = append(encData, []byte{eCRCLow, eCRCHigh}...)

	// Prefix
	encData = append([]byte{stex}, encData...)

	return encData
}

func (p *Packet) payload() []byte {
	// seq/slave + length + data
	return append([]byte{p.seq, byte(len(p.args) + 1)}, p.Data()...)
}

func (p *Packet) Data() []byte {
	return append([]byte{p.code}, p.args...)
}

func (p *Packet) Bytes() []byte {
	crcLow, crcHigh := p.Checksum()
	return append(p.payload(), crcLow, crcHigh)
}

func (p *Packet) String() string {
	return fmt.Sprintf("seq: %v code: %x args: %s", p.seq, p.code, hexDump(p.args))
}

func (p *Packet) Checksum() (crcLow, crcHigh byte) {
	crcLow, crcHigh = crc16(p.payload())
	return
}

func Decode(data []byte) (*Packet, error) {
	// Strip prefix and unstuff data
	data = bytes.ReplaceAll(data, []byte{stx, stx}, []byte{stx})

	log.Trace().Hex("data", data).Msg("decoding")

	packets := make([]Packet, 0, 2)
	for {
		i := bytes.LastIndex(data, []byte{stx})
		if i == -1 {
			break
		}

		p := data[i+1:]

		if len(p) == 0 {
			continue
		}

		log.Trace().Hex("packet", p).Msg("processing packet")

		seq := p[0]
		length := p[1]

		crcLow := p[length+2]
		crcHigh := p[length+3]

		responseData := p[2 : length+2]
		code := responseData[0]
		args := responseData[1:]

		packet := &Packet{
			seq:  seq,
			code: code,
			args: args,
		}

		expectedCRCLow, expectedCRCHigh := packet.Checksum()
		if crcLow != expectedCRCLow || crcHigh != expectedCRCHigh {
			return nil, fmt.Errorf("invalid crc. low: %x - %x high: %x - %x", crcLow, expectedCRCLow, crcHigh, expectedCRCHigh)
		}

		if err := check(packet); err != nil {
			return nil, err
		}

		packets = append(packets, *packet)

		data = data[:i]
	}

	return &packets[0], nil
}

func Encode(packet *Packet) ([]byte, error) {
	if err := check(packet); err != nil {
		return nil, err
	}

	// Stuffing
	stuffedData := bytes.ReplaceAll(packet.Bytes(), []byte{stx}, []byte{stx, stx})

	// Prefix stuffed data with delimiter
	final := append([]byte{stx}, stuffedData...)

	return final, nil
}

func check(packet *Packet) error {
	if len(packet.Data()) > maxDataSize {
		return fmt.Errorf("length of data %s exceeds max data-size of %d", hexDump(packet.Data()), maxDataSize)
	}

	return nil
}

func crc16(source []byte) (byte, byte) {
	seed := 0xFFFF
	poly := 0x8005
	crc := seed
	for _, b := range source {
		crc ^= int(b) << 8
		for j := 0; j < 8; j++ {
			if crc&0x8000 != 0 {
				crc = ((crc << 1) & 0xffff) ^ poly
			} else {
				crc <<= 1
			}
		}
	}
	return byte(crc & 0xFF), byte((crc >> 8) & 0xFF)
}
