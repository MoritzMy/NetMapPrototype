package eth

import (
	"encoding/binary"
	"fmt"
	"net"
)

type EthernetHeader struct {
	DestinationMAC net.HardwareAddr
	SourceMAC      net.HardwareAddr
	EtherType      uint16
}

func NewEthernetHeader(dest net.HardwareAddr, source net.HardwareAddr, ethertype uint16) EthernetHeader {
	return EthernetHeader{
		DestinationMAC: dest,
		SourceMAC:      source,
		EtherType:      ethertype,
	}
}

func (header EthernetHeader) Marshal() ([]byte, error) {
	if len(header.DestinationMAC) != MACAdressLength {
		return nil, fmt.Errorf("destination MAC adress has length of %d bytes, not the required length of 6 byte", len(header.DestinationMAC))
	}

	if len(header.SourceMAC) != MACAdressLength {
		return nil, fmt.Errorf("source MAC has length of %d bytes, not the required length of 6 bytes", len(header.SourceMAC))
	}

	b := make([]byte, 0, EthernetHeaderLength)
	b = append(b, header.DestinationMAC...)
	b = append(b, header.SourceMAC...)
	b = binary.BigEndian.AppendUint16(b, header.EtherType)
	return b, nil
}

func (header *EthernetHeader) Unmarshal(b []byte) error {
	if len(b) != EthernetHeaderLength {
		return fmt.Errorf("provided byte array does not meet required length of 14")
	}

	header.DestinationMAC = b[0:6]
	header.SourceMAC = b[6:12]
	header.EtherType = binary.BigEndian.Uint16(b[12:14])
	return nil
}
