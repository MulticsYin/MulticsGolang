package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// Package was a massage
type Package struct {
	Version [2]byte
	Length  int16
	Data    []byte
}

func (p *Package) String() string {
	return fmt.Sprintf("Version: %v, Length: %v, Data: %v, Data(String): %s\n", p.Version, p.Length, p.Data, string(p.Data))
}

// Unpack function
func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	p.Data = make([]byte, p.Length)
	err = binary.Read(reader, binary.BigEndian, &p.Data)
	return err
}

func handleReq(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && 'V' == data[0] {
			if 4 < len(data) {
				dataLen := int16(0)
				binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &dataLen)
				if int(dataLen)+4 <= len(data) {
					return int(dataLen) + 4, data[:int(dataLen)+4], nil
				}
			}
		}
		return
	})

	for scanner.Scan() {
		scanPack := new(Package)
		scanPack.Unpack(bytes.NewReader(scanner.Bytes()))
		log.Println(scanPack)
		log.Println(scanner.Bytes())
	}
}

func main() {
	lister, err := net.Listen("tcp", ":8080")
	if nil != err {
		log.Fatalf("Listen :8080 error: %s\n", err.Error())
	}
	for {
		conn, err := lister.Accept()
		if nil != err {
			log.Fatalf("Lister accept error: %s\n", err.Error())
		}
		go handleReq(conn)
	}
}
