package sflow

import (
	"encoding/binary"
	"io"
)

type decoder struct {
	reader io.ReadSeeker
}

func NewDecoder(r io.ReadSeeker) *decoder {
	return &decoder{
		reader: r,
	}
}

func (d *decoder) Use(r io.ReadSeeker) {
	d.reader = r
}

func (d *decoder) Decode() (*Datagram, error) {
	// Decode headers first
	dgram := &Datagram{}
	var err error

	err = binary.Read(d.reader, binary.BigEndian, &dgram.Version)
	if err != nil {
		return nil, err
	}

	err = binary.Read(d.reader, binary.BigEndian, &dgram.IpVersion)
	if err != nil {
		return nil, err
	}

	ipLen := 4
	if dgram.IpVersion == 2 {
		ipLen = 16
	}

	ipBuf := make([]byte, ipLen)
	_, err = d.reader.Read(ipBuf)
	if err != nil {
		return nil, err
	}

	dgram.IpAddress = ipBuf

	err = binary.Read(d.reader, binary.BigEndian, &dgram.SubAgentId)
	if err != nil {
		return nil, err
	}

	err = binary.Read(d.reader, binary.BigEndian, &dgram.SequenceNumber)
	if err != nil {
		return nil, err
	}

	err = binary.Read(d.reader, binary.BigEndian, &dgram.Uptime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(d.reader, binary.BigEndian, &dgram.NumSamples)
	if err != nil {
		return nil, err
	}

	return dgram, nil
}
