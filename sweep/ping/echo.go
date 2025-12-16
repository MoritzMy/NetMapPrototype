package ping

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	echoCode            = 0
	echoType            = 8
	checksumPlaceholder = 0
	maxPayload          = 56
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Unmarshaler interface {
	Unmarshall([]byte) (any, error)
}

type ICMPPacket struct {
	Type uint8
	Code uint8
}

type EchoRequest struct {
	ICMPPacket
	Identifier     uint16
	SequenceNumber uint16
	Payload        []byte
}

func CreateEchoRequest(identifier uint16, sequenceNumber uint16, payload []byte) EchoRequest {
	return EchoRequest{
		ICMPPacket: ICMPPacket{
			Type: echoType,
			Code: echoCode,
		},
		Identifier:     identifier,
		SequenceNumber: sequenceNumber,
		Payload:        payload,
	}
}

// Marshal parses the ICMP Type and ICMP Code of the Packet and sets the Checksum Placeholder
func (packet ICMPPacket) Marshal() ([]byte, error) {
	b := make([]byte, 0)
	b = append(b, packet.Type, packet.Code)
	b = binary.BigEndian.AppendUint16(b, checksumPlaceholder)

	return b, nil
}

func (packet EchoRequest) Marshal() ([]byte, error) {
	if len(packet.Payload) > maxPayload {
		return nil, fmt.Errorf("marshal icmp request: payload size %d exceeds limit of %d Bytes", len(packet.Payload), maxPayload)
	}

	b := make([]byte, 0, 8+len(packet.Payload))

	hdr, err := packet.ICMPPacket.Marshal()

	if err != nil {
		return nil, errors.New("failed to parse ICMP Headers")
	}

	b = append(b, hdr...)
	b = binary.BigEndian.AppendUint16(b, packet.Identifier)
	b = binary.BigEndian.AppendUint16(b, packet.SequenceNumber)
	b = append(b, packet.Payload...)

	cs := computeChecksum(b)

	binary.BigEndian.PutUint16(b[2:4], cs)

	return b, nil
}

func (packet *EchoRequest) Unmarshal(data []byte) error {
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
