package ssp

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type NVx struct {
	conn            *Connection
	FirmwareVersion string `json:"firmware_version"`
	DatasetVersion  string `json:"dataset_version"`
	channels        []Channel
}

func DialNVx(portName string) (*NVx, error) {
	conn, err := Dial(portName, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NVx at %s: %w", portName, err)
	}

	return &NVx{
		conn:     conn,
		channels: make([]Channel, 0, 5),
	}, nil
}

func (nv *NVx) Close() error {
	return nv.conn.Close()
}

func (nv *NVx) Init() error {
	{
		_, err := nv.conn.Sync()
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}
	}

	{
		_, err := nv.conn.SetChannelInhibits(0xFFFF)
		if err != nil {
			return fmt.Errorf("setting channel inhibits failed: %w", err)
		}
	}

	{
		version, err := nv.conn.GetFirmwareVersion()
		if err != nil {
			return fmt.Errorf("failed to read firmware version: %w", err)
		}
		nv.FirmwareVersion = version.FirmwareVersion
		log.Info().Str("firmware_version", version.FirmwareVersion).Msgf("read firmware version")
	}

	{
		version, err := nv.conn.GetDatasetVersion()
		if err != nil {
			return fmt.Errorf("failed to read dataset version: %w", err)
		}
		nv.DatasetVersion = version
		log.Info().Str("dataset_version", version).Msgf("read dataset version")
	}

	var l zerolog.Logger
	{
		resp, err := nv.conn.SetupRequest()
		if err != nil {
			return err
		}
		l = log.With().
			Str("country_code", resp.CountryCode).
			Str("firmware_version", resp.FirmwareVersion).
			Uint8("unit_type", uint8(resp.UnitType)).
			Logger()

		l.Info().
			Msgf("setup found unit: %s", resp.UnitType.String())

		if resp.UnitType != BanknoteValidator {
			return fmt.Errorf("unit-type is not supported: %s", resp.UnitType)
		}

		setup, err := resp.BanknoteValidator()
		if err != nil {
			return err
		}
		if log.Debug().Enabled() {
			supportedNotes := make([]string, len(setup.Channels))
			for i := range setup.Channels {
				supportedNotes[i] = strconv.Itoa(int(setup.Channels[i].Value))
			}

			l.Debug().
				Strs("supported_notes", supportedNotes).
				Uint8("num_channels", setup.NumberOfChannels).
				Msg("validator set up")
		}

		if setup.ProtocolVersion < 6 {
			// @todo check currency
			log.Info().Msgf("flashed currency is %s", setup.CountryCode)
		}

		nv.channels = setup.Channels
	}

	{
		resp, err := nv.conn.Enable()
		if err != nil {
			return err
		}
		log.Debug().Str("code", resp.Code.String()).Msg("enabled device")
	}

	return nil
}

func (nv *NVx) Poll(c chan os.Signal) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-c:
			return
		case <-ticker.C:
			resp, err := nv.conn.Poll()
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			log.Trace().Str("code", resp.Code.String()).Msg("polled")

			if resp.Code != OK {
				log.Fatal().Msg("poll failed")
			}

			if err := nv.handlePollResponse(resp); err != nil {
				log.Fatal().Err(err).Msg("")
			}
		}
	}
}

func (nv *NVx) handlePollResponse(r *PollResponse) error {
	var (
		channelNo uint8
	)
	for i := 0; i < len(r.Args); i++ {
		code := PollEventResponseCode(r.Args[i])

		l := log.With().Str("code", code.String()).Logger()
		l.Trace().Msgf("processing poll-event")

		// All events with have a channel as single byte arg
		switch code {
		case ReadNote: // 1 byte showing the channel of the recognised note. This will be zero if the note has not yet been recognised
			fallthrough
		case CreditNote: // 1 byte showing the channel of the credited note.
			fallthrough
		case FraudAttempt: // 1 byte showing the channel of the note that was in process when the fraud was detected. This will be zero if the note has not yet been recognised
			fallthrough
		case NoteClearedFromFront: // 1 byte showing the channel of the note rejected (0 if not known)
			fallthrough
		case NoteClearedToCashbox: // 1 byte showing the channel of the note stacked (0 if not known)
			i += 1
			if i > len(r.Args) {
				l.Fatal().Msg("got no params, but expected one")
			}
			channelNo = r.Args[i]
			if channelNo == 0 {
				l.Trace().Msg("note has not yet been recognised")
			}
		}

		switch code {
		case SafeNoteJam: // no data
			fallthrough
		case UnsafeNoteJam: // no data
			fallthrough
		case StackerFull: // no data
			fallthrough
		case SlaveReset: // no data
			fallthrough
		case Disabled: // no data
			l.Warn().Msg(code.Desc())
		}

		switch code {
		case ReadNote:
		case CreditNote:
			if err := nv.handleCredit(channelNo); err != nil {
				l.Warn().Err(err).Uint8("channel", channelNo).Msg("failed to handle credit")
			}
		case NoteRejected:
		case Disabled:
		}
	}
	return nil
}

func (nv *NVx) handleCredit(ch uint8) error {
	_, disErr := nv.conn.Disable()
	if disErr != nil {
		log.Fatal().Err(disErr).Msg("failed to disable device")
	}

	chIndex := ch - 1

	if chIndex > uint8(len(nv.channels)) {
		log.Fatal().Msgf("ch %d > %d", ch, len(nv.channels))
	}
	channel := nv.channels[chIndex]
	log.Info().Uint32("note_value", channel.Value).Msgf("received note")

	cresp, err := nv.conn.GetCounters()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	if log.Debug().Enabled() {
		log.Debug().
			Uint8("number_of_counters_set", cresp.NumberOfCountersSet).
			Uint32("notes_rejected_count", cresp.NotesRejected).
			Uint32("notes_stacked_count", cresp.NotesStacked).
			Uint32("notes_stored_count", cresp.NotesStored).
			Uint32("notes_transferred_count", cresp.NotesTransferred).
			Msg("checked counters")
	}

	_, enErr := nv.conn.Enable()
	if enErr != nil {
		log.Fatal().Err(enErr).Msg("")
	}

	return nil
}
