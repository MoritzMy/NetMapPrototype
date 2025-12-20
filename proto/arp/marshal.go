package arp

const (
	ARPPacketLength = 42
)

func Marshal(request ARPRequest) ([]byte, error) {
	b := make([]byte, 0, ARPPacketLength)
	base, err := request.EthernetHeader.Marshal()

	if err != nil {
		return nil, err
	}

	content, err := request.Marshal()

	if err != nil {
		return nil, err
	}

	b = append(b, base...)
	b = append(b, content...)

	return b, nil
}
