package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

// Package was a massage
type Package struct {
	Version [2]byte
	Length  int16
	Data    []byte
}

// Pack function
func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Data)
	return err
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if nil != err {
		log.Fatalf("Connect 127.0.0.1:8080 error: %s\n", err.Error())
	}
	defer conn.Close()

	pack := &Package{
		Version: [2]byte{'V', '1'},
		Data:    []byte(time.Now().String()),
	}
	pack.Length = int16(len(pack.Data))

	buf00 := new(bytes.Buffer)
	pack.Pack(buf00)
	buf01 := new(bytes.Buffer)
	pack.Pack(buf01)
	pack.Pack(buf01)
	pack.Pack(buf01)
	pack.Pack(buf01)
	pack.Pack(buf01)
	pack.Pack(buf01)

	conn.Write(buf00.Bytes())
	log.Println(buf00)
	conn.Write(buf01.Bytes())
	log.Println(buf01)
}
