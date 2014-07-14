package ttygif

import (
	"encoding/binary"
	"io"
)

// TimeVal type
type TimeVal struct {
	Sec  int32
	Usec int32
}

// Add returns sum of TimeVal data
func (tv1 TimeVal) Add(tv2 TimeVal) TimeVal {
	sec := tv1.Sec + tv2.Sec
	usec := tv1.Usec + tv2.Usec
	for usec >= 1000000 {
		sec++
		usec -= 1000000
	}
	return TimeVal{
		Sec:  sec,
		Usec: usec,
	}
}

// Subtract returns diff of TimeVal data
func (tv1 TimeVal) Subtract(tv2 TimeVal) TimeVal {
	sec := tv1.Sec - tv2.Sec
	usec := tv1.Usec - tv2.Usec
	if usec < 0 {
		sec--
		usec += 1000000
	}
	return TimeVal{
		Sec:  sec,
		Usec: usec,
	}
}

// Header type
type Header struct {
	tv  TimeVal
	len uint32
}

// TtyData type
type TtyData struct {
	TimeVal TimeVal
	Buffer  *[]byte
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
	// 4byte x 3 ?
	bufHeader := make([]byte, 12)
	_, err = r.reader.Read(bufHeader)
	if err != nil {
		return
	}

	header := &Header{
		tv: TimeVal{
			Sec:  int32(r.order.Uint32(bufHeader[0:4])),
			Usec: int32(r.order.Uint32(bufHeader[4:8])),
		},
		len: r.order.Uint32(bufHeader[8:12]),
	}

	bufBody := make([]byte, header.len)
	_, err = r.reader.Read(bufBody)
	if err != nil {
		return
	}
	return &TtyData{
		TimeVal: header.tv,
		Buffer:  &bufBody,
	}, nil
}
