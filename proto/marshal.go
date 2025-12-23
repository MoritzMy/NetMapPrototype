package proto

func Marshal[T Packet](packet T) ([]byte, error) {
	context, err := packet.Marshal()

	if err != nil {
		return nil, err
	}

	headers, err := packet.GetHeaders().Marshal(context)

	if err != nil {
		return nil, err
	}

	b := append(headers, context...)

	return b, nil
}

func Unmarshal[T Packet](data []byte, zero T) error {
	h := zero.GetHeaders()

	if err := zero.GetHeaders().Unmarshal(data[:h.Len()]); err != nil {
		return err
	}
	if err := zero.Unmarshal(data[h.Len():zero.Len()]); err != nil {
		return err
	}
	return nil
}
