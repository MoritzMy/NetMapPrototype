package icmp

import "encoding/binary"

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

func Unmarshal[T ICMPPacket](data []byte, zero T) error {
	if err := zero.GetHeaders().Unmarshal(data[0:4]); err != nil {
		return err
	}

	if err := zero.Unmarshal(data[4:]); err != nil {
		return err
	}
	return nil

}
