package ssp

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

const (
	readTimeout       = time.Second * 30
	seqIDMask   uint8 = 0x80
	maxSlaveID  uint8 = 0x7D
)

type Connection struct {
	port    serial.Port
	count   uint32
	seqFlag bool
	slaveID uint8
}

func hexDump(data []byte) []string {
	i := 0
	h := hex.EncodeToString(data)
	parts := make([]string, 0, len(h)/2)
	for i < len(h) {
		parts = append(parts, h[i:i+2])
		i += 2
	}
	return parts
}

func Dial(portName string, slaveID uint8) (*Connection, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list ports")
	}
	if len(ports) == 0 {
		log.Fatal().Msg("No serial ports found!")
	}

	i := sort.SearchStrings(ports, portName)
	if i == len(ports) || ports[i] != portName {
		log.Printf("port %s not found. available ports: %v", portName, ports)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.TwoStopBits,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	if err := port.SetReadTimeout(readTimeout); err != nil {
		return nil, err
	}

	{
		status, err := port.GetModemStatusBits()
		if err != nil {
			return nil, err
		}
		log.Debug().
			Bool("dsr", status.DSR).
			Bool("dcd", status.DCD).
			Bool("ri", status.RI).
			Bool("cts", status.CTS).
			Msg("modem status bits")
	}

	{
		err := port.SetRTS(true)
		if err != nil {
			return nil, err
		}
		log.Debug().Msg("set RTS ON")
	}

	{
		err := port.SetDTR(true)
		if err != nil {
			return nil, err
		}
		log.Debug().Msg("set DTR ON")
	}

	if slaveID > maxSlaveID {
		return nil, fmt.Errorf("invalid slaveID %b. max allowed slave id is %b", slaveID, maxSlaveID)
	}

	return &Connection{port: port, slaveID: slaveID}, nil
}

func (c *Connection) Close() error {
	return c.port.Close()
}

func (c *Connection) SequenceFlag() bool {
	return c.seqFlag
}

func (c *Connection) SlaveID() uint8 {
	return c.slaveID
}

func (c *Connection) flipSequenceFlag() {
	// Toggle (XOR) sequence flag (bit 7)
	c.seqFlag = !c.seqFlag
}

func (c *Connection) seq() uint8 {
	seq := c.slaveID
	if c.seqFlag {
		seq |= seqIDMask
	}
	return seq
}

func (c *Connection) send(packet *Packet) error {
	log.Trace().Msgf("sending packet: %v (count: %d)", packet, c.count)

	data, err := Encode(packet)
	if err != nil {
		return err
	}

	n, err := c.port.Write(data)
	if err != nil {
		return err
	}
	log.Trace().Msgf("sent %d bytes: %s", n, hexDump(data))

	c.count += 1
	return nil
}

func (c *Connection) execute(cmd Command) (*Packet, error) {
	// @todo Check if supported
	//_, exists := ssp.supportedCommands[cmd.Code]
	//if !exists {
	//	return nil, fmt.Errorf("command (code: %d) is not supported", cmd.Code)
	//}

	//if err := c.port.ResetInputBuffer(); err != nil {
	//	return nil, err
	//}
	cmdPacket := NewPacket(byte(cmd.Code), cmd.Args, c.seq())
	if err := c.send(cmdPacket); err != nil {
		return nil, err
	}

	responsePacket, err := c.checkForIncomingData()
	if err != nil {
		return nil, err
	}
	log.Trace().Hex("packet", responsePacket.Bytes()).Msg("incoming data")

	if c.seq()^responsePacket.seq > 0 {
		return nil, fmt.Errorf("seq broken: %b != %b", c.seq(), responsePacket.seq)
		// @todo implement retry mechanism
	}

	c.flipSequenceFlag()

	return responsePacket, nil
}

func (c *Connection) checkForIncomingData() (*Packet, error) {
	buff := make([]byte, 512)
	n, err := c.port.Read(buff)
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("read %d bytes: %s", n, hexDump(buff))
	if n == 0 {
		log.Debug().Msg("EOF")
	}

	log.Trace().Msgf("decoding: %s", hexDump(buff[:n]))
	packet, err := Decode(buff[:n])
	if err != nil {
		return nil, err
	}

	return packet, nil
}
