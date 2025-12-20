package arp

const (
	ARPPacketLength = 42
)

func Marshal(packet ARPRequest) ([]byte, error) {
	b := make([]byte, 0, ARPPacketLength)
	base, err := packet.EthernetHeader.Marshal()

	if err != nil {
		return nil, err
	}

	content, err := packet.Marshal()

	if err != nil {
		return nil, err
	}

	b = append(b, base...)
	b = append(b, content...)

	return b, nil
}

func Unmarshal(b []byte, packet *ARPRequest) error {
	if err := packet.EthernetHeader.Unmarshal(b[:14]); err != nil {
		return err
	}

	if err := packet.Unmarshal(b[14:]); err != nil {
		return err
	}

	return nil
}
