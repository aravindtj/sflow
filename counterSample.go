package sflow

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	TypeGenericInterfaceCountersRecord = 1
	TypeEthernetCountersRecord         = 2
	TypeTokenRingCountersRecord        = 3
	TypeVgCountersRecord               = 4
	TypeVlanCountersRecord             = 5

	TypeProcessorCountersRecord  = 1001
	TypeHostCpuCountersRecord    = 2003
	TypeHostMemoryCountersRecord = 2004
	TypeHostDiskCountersRecord   = 2005
	TypeHostNetCountersRecord    = 2006

	// Custom (Enterprise) types
	TypeApplicationCountersRecord = (1)<<12 + 1
)

type CounterSample struct {
	SequenceNum      uint32
	SourceIdType     byte
	SourceIdIndexVal uint32 // NOTE: this is 3 bytes in the datagram
	NumRecords       uint32
	Records          []Record
}

func (s *CounterSample) SampleType() int {
	return TypeCounterSample
}

func (s *CounterSample) GetRecords() []Record {
	return s.Records
}

func decodeCounterSample(r io.ReadSeeker) (Sample, error) {
	s := &CounterSample{}

	var err error

	err = binary.Read(r, binary.BigEndian, &s.SequenceNum)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.SourceIdType)
	if err != nil {
		return nil, err
	}

	var srcIdIndexVal [3]byte
	n, err := r.Read(srcIdIndexVal[:])
	if err != nil {
		return nil, err
	}

	if n != 3 {
		return nil, errors.New("sflow: counter sample decoding error")
	}

	s.SourceIdIndexVal = uint32(srcIdIndexVal[2]) | uint32(srcIdIndexVal[1]<<8) |
		uint32(srcIdIndexVal[0]<<16)

	err = binary.Read(r, binary.BigEndian, &s.NumRecords)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < s.NumRecords; i++ {
		format, length := uint32(0), uint32(0)

		err = binary.Read(r, binary.BigEndian, &format)
		if err != nil {
			return nil, err
		}

		err = binary.Read(r, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		var rec Record

		switch format {
		case TypeGenericInterfaceCountersRecord:
			rec, err = decodeGenericInterfaceCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeEthernetCountersRecord:
			rec, err = decodeEthernetCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeTokenRingCountersRecord:
			rec, err = decodeTokenRingCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeVgCountersRecord:
			rec, err = decodeVgCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeVlanCountersRecord:
			rec, err = decodeVlanCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeProcessorCountersRecord:
			rec, err = decodeProcessorCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeHostCpuCountersRecord:
			rec, err = decodeHostCpuCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeHostMemoryCountersRecord:
			rec, err = decodeHostMemoryCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeHostDiskCountersRecord:
			rec, err = decodeHostDiskCountersRecord(r)
			if err != nil {
				return nil, err
			}

		case TypeHostNetCountersRecord:
			rec, err = decodeHostNetCountersRecord(r)
			if err != nil {
				return nil, err
			}

		default:
			_, err := r.Seek(int64(length), 1)
			if err != nil {
				return nil, err
			}

			continue
		}

		s.Records = append(s.Records, rec)
	}

	return s, nil
}
