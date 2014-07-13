package main

import (
	"encoding/binary"
	"io"
)

// TimeVal type
type TimeVal struct {
	sec  uint32
	usec uint32
}

// Subtract returns diff of TimeVal data
func (tv1 TimeVal) Subtract(tv2 TimeVal) TimeVal {
	sec := tv1.sec - tv2.sec
	usec := tv1.usec - tv2.usec
	if usec < 0 {
		sec--
		usec += 1000000
	}
	return TimeVal{
		sec:  sec,
		usec: usec,
	}
}

// Header type
type Header struct {
	tv  TimeVal
	len uint32
}

// TtyData type
type TtyData struct {
	Header *Header
	Buffer *[]byte
}

// TtyReader type
type TtyReader struct {
	reader io.Reader
	order  binary.ByteOrder
}

// NewTtyReader returns TtyReader instance
func NewTtyReader(r io.Reader) *TtyReader {
	return &TtyReader{
		reader: r,
		order:  binary.LittleEndian,
	}
}

// ReadData returns next TtyData
func (r *TtyReader) ReadData() (data *TtyData, err error) {
	bufHeader := make([]byte, 12)
	_, err = r.reader.Read(bufHeader)
	if err != nil {
		return
	}

	header := &Header{
		tv: TimeVal{
			sec:  r.order.Uint32(bufHeader[0:4]),
			usec: r.order.Uint32(bufHeader[4:8]),
		},
		len: r.order.Uint32(bufHeader[8:12]),
	}

	bufBody := make([]byte, header.len)
	_, err = r.reader.Read(bufBody)
	if err != nil {
		return
	}
	return &TtyData{
		Header: header,
		Buffer: &bufBody,
	}, nil
}
