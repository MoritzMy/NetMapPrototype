package main

import (
	"encoding/binary"
	"fmt"
)

const (
	echoCode       = 0
	echoType       = 8
	maxPayload     = 56
	echoHeaderSize = 4
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type ICMPHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
}

type ICMPPacket interface {
	Marshaler
	Unmarshaler
	GetHeaders() ICMPHeader
	SetHeaders(header ICMPHeader)
}

type EchoICMPPacket struct {
	ICMPHeader
	Identifier     uint16
	SequenceNumber uint16
	Payload        []byte
}

func CreateEchoRequest(identifier uint16, sequenceNumber uint16, payload []byte) EchoICMPPacket {
	return EchoICMPPacket{
		ICMPHeader: ICMPHeader{
			Type: echoType,
			Code: echoCode,
		},
		Identifier:     identifier,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
	}
}

func (req EchoICMPPacket) GetHeaders() ICMPHeader {
	return req.ICMPHeader
}

func (req EchoICMPPacket) SetHeaders(header ICMPHeader) {
	req.ICMPHeader = header
}

func Marshal[T ICMPPacket](packet T) ([]byte, error) {
	base, err := packet.GetHeaders().Marshal()

	if err != nil {
		return nil, err
	}

	content, err := packet.Marshal()

	if err != nil {
		return nil, err
	}

	res := append(base, content...)

	cs := computeChecksum(res)

	binary.BigEndian.PutUint16(res[2:4], cs)

	return res, nil
}

// Marshal parses the ICMP Type and ICMP Code of the Packet and sets the Checksum Placeholder
func (packet ICMPHeader) Marshal() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, packet.Type, packet.Code)
	b = binary.BigEndian.AppendUint16(b, packet.Checksum)

	return b, nil
}

func (headers ICMPHeader) Unmarshal(b []byte) error {
	headers.Type = b[0]
	headers.Code = b[1]
	headers.Checksum = binary.BigEndian.Uint16(b[2:4])
	return nil
}

func (packet EchoICMPPacket) Marshal() ([]byte, error) {
	if len(packet.Payload) > maxPayload {
		return nil, fmt.Errorf("marshal icmp request: payload size %d exceeds limit of %d Bytes", len(packet.Payload), maxPayload)
	}

	b := make([]byte, 0, echoHeaderSize+len(packet.Payload))
	b = binary.BigEndian.AppendUint16(b, packet.Identifier)
	b = binary.BigEndian.AppendUint16(b, packet.SequenceNumber)
	b = append(b, packet.Payload...)

	return b, nil
}

func Unmarshal[T ICMPPacket](data []byte, zero T) error {
	if err := zero.GetHeaders().Unmarshal(data[0:4]); err != nil {
		return err
	}

	if err := zero.Unmarshal(data[4:]); err != nil {
		return err
	}
	return nil

}

func (packet EchoICMPPacket) Unmarshal(data []byte) error {
	packet.Identifier = binary.BigEndian.Uint16(data[0:2])
	packet.SequenceNumber = binary.BigEndian.Uint16(data[2:4])
	packet.Payload = data[4:]
	return nil
}

// computeChecksum computes the checksum of a package, by splitting it up into 16 Bit words,
// adding those words together and performing an end around carry until the sum is also a 16 Bit word.
// In the Case of ICMP, while the Checksum is not computed, a Placeholder should be used of which the 16 Bit
// word value is 0
func computeChecksum(request []byte) uint16 {
	sum := uint32(0)

	// Turn the bytes into 16 Bit Words and add them up
	for i := 0; i+1 < len(request); i += 2 {
		sum += (uint32(request[i]) << 8) + uint32(request[i+1])
	}

	if len(request)%2 != 0 {
		sum += uint32(request[len(request)-1]) << 8
	}

	// sum needs to be a valid uint16, otherwise an end around carry is performed
	for sum>>16 != 0 {
		sum = uint32(uint16(sum)) + sum>>16
	}

	// One's complement
	var checksum = ^uint16(sum)

	return checksum
}

func main() {
	req := CreateEchoRequest(0, 0, []byte("Hello World! :)"))
	res, err := Marshal(req)

	if err != nil {
		panic(fmt.Errorf("marshal echo request: %v", err))
	}

	var u EchoICMPPacket
	if err := Unmarshal(res, &u); err != nil {
		panic(err)
	}

	fmt.Println(res, "\n", u)
}
