package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
)

// TimeVal type
type TimeVal struct {
	sec  uint32
	usec uint32
}

// Subtract returns diff of TimeVal data
func (tv1 *TimeVal) Subtract(tv2 *TimeVal) *TimeVal {
	sec := tv1.sec - tv2.sec
	usec := tv1.usec - tv2.usec
	if usec < 0 {
		sec--
		usec += 1000000
	}
	return &TimeVal{
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

// TtyRead returns channel
func TtyRead(file *os.File) (ch chan *TtyData) {
	reader := bufio.NewReader(file)

	ch = make(chan *TtyData)
	go func() {
		defer close(ch)
		for {
			header, err := readHeader(reader)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}

			buf := make([]byte, header.len)
			_, err = reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			ch <- &TtyData{
				Header: header,
				Buffer: &buf,
			}
		}
	}()
	return ch
}

func readHeader(reader *bufio.Reader) (header *Header, err error) {
	buf := make([]byte, 12)
	_, err = reader.Read(buf)
	if err != nil {
		return
	}

	byteOrder := binary.LittleEndian
	return &Header{
		tv: TimeVal{
			sec:  byteOrder.Uint32(buf[0:4]),
			usec: byteOrder.Uint32(buf[4:8]),
		},
		len: byteOrder.Uint32(buf[8:12]),
	}, nil
}
